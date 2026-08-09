---
title: openHop Repeater Installation Guide
description: Installing a MeshCore repeater on a Raspberry Pi with openHop Repeater and the optional openHop Console.
nav_weight: 4
authors:
  - beanfield
  - Json_18
series:
  - Guide
date: 2026-07-31T00:00:00-04:00
nav_icon:
  vendor: bs
  name: wrench
  color: green
---

Setting up a MeshCore repeater on a Raspberry Pi using **openHop Repeater**, a Python repeater daemon, and the optional **openHop Console** web UI.

<!--more-->

{{< notice note "Repeater operators" >}}
Building infrastructure. If your node is paired to a phone, you do not need any of this — see [Companion Node Setup]({{< relref "/docs/meshcore/companion/index.md" >}}).
{{< /notice >}}

{{< notice note "Renamed projects" >}}
openHop Repeater was `pyMC_Repeater`; openHop Console was `pyMC Console`. Guides referencing `/etc/pymc_repeater`, the `pymc-repeater` service, or `github.com/dmduran12/…` are out of date. The paths below are current.
{{< /notice >}}

## TX power on 2 W boards

{{< notice warning >}}
Some LoRa HATs are permanently damaged by the default TX power setting. Establish which board you have before you set a power level.
{{< /notice >}}

Boards built around a 2 W module — the Ebyte **E22-900M33S**, roughly 33 dBm — contain an SX1262 driving an internal power amplifier. The PA's absolute-maximum RF input is about `+10 dBm`, and Ebyte's manual states that the output power of the RF chip cannot exceed 9 dBm.

`tx_power` is the **SX1262's** output, applied 1:1 with no compensation for PA gain. The software default of `22` is roughly 16× over the PA's absolute maximum. The wehooper4 NebraHat/ZebraHAT documentation states that a power level above `8` will damage the PA.

| Board type | Safe `tx_power` |
|---|---|
| 2 W / 33 dBm — NebraHat 2W, ZebraHAT 2W, Femtofox 2W | `8` or lower |
| 1 W / 30 dBm — NebraDuo, PiMesh, Zebra | `18` |
| Bare SX1262, no PA — Waveshare, MeshAdv, uConsole | up to `22` |

The RF output gives no warning. The chain is already saturated at the top of the safe range, so the excess dissipates inside the module as heat. Nothing in the software validates this — the API accepts anything up to 30 dBm. A 1 W and a 2 W variant of the same HAT are often visually identical and named almost identically.

## Requirements

- Raspberry Pi 3, 4, 5, or Zero 2 W
- Raspberry Pi OS Bookworm or later, 32- or 64-bit
- A supported SX1262 or SX1268 LoRa HAT
- Network connection; SPI available — the installer configures it

## 1. Install the repeater

```bash {linenos=false}
git clone https://github.com/openhop-dev/openhop_repeater.git
cd openhop_repeater
sudo bash ./manage.sh install
```

The installer creates a `repeater` service user with hardware access, installs to `/opt/openhop_repeater` in its own venv, creates config, log, and data directories, runs an interactive hardware wizard, and enables the systemd service.

In the wizard, select your exact board. That selection sets your TX power — the wrong variant gives the wrong power level.

| Path | Contents |
|---|---|
| `/opt/openhop_repeater/` | venv and application files |
| `/etc/openhop_repeater/config.yaml` | configuration |
| `/var/lib/openhop_repeater/` | `repeater.db`, metrics, radio presets |
| `/var/log/openhop_repeater/` | logs, also in the journal |
| `openhop-repeater.service` | the systemd unit |

## 2. Verify TX power before transmitting

Immediately after install, and after any config change:

```bash {linenos=false}
sudo grep -n tx_power /etc/openhop_repeater/config.yaml
```

Expect exactly one line, under the `radio:` section, at or below your board's limit.

{{< notice warning >}}
`tx_power` belongs under `radio:`, not under `sx1262:`. A value placed under `sx1262:` is silently ignored and the radio runs at whatever `radio.tx_power` says.
{{< /notice >}}

## 3. Install the Console (optional)

The Console is a static web UI with no backend of its own — the repeater serves it. Install the repeater first.

```bash {linenos=false}
git clone https://github.com/Treehouse-00/pymc_console-dist.git pymc_console
cd pymc_console
sudo bash manage.sh install
```

Browse to `http://<pi-ip>:8000/` and log in with your repeater credentials.

{{< notice note >}}
Follow the repo's `README.md`, not its `INSTALL.md` — the latter is stale and still refers to pre-rename paths and services.
{{< /notice >}}

## 4. Managing the service

```bash {linenos=false}
sudo systemctl status openhop-repeater
sudo systemctl restart openhop-repeater
sudo journalctl -u openhop-repeater -f
```

## Upgrading

`manage.sh upgrade` installs whatever is in your working tree; it fetches nothing. `git pull` first.

```bash {linenos=false}
cd ~/openhop_repeater
git pull
sudo bash ./manage.sh upgrade
```

The same applies to the Console, so the script itself is current before it runs.

```bash {linenos=false}
cd ~/pymc_console
git pull
sudo bash manage.sh upgrade
```

The Console does not restart the repeater — the assets are static. Hard-refresh the browser with `Ctrl+Shift+R`.

The Console's web UI can also upgrade, offering a release channel: `main` for stable snapshots, `dev` for latest commits. `main` is updated infrequently.

### Before you upgrade

- Back up. There is no rollback command.

  ```bash {linenos=false}
  sudo cp -a /etc/openhop_repeater/config.yaml /root/config.yaml.bak
  sudo cp -a /var/lib/openhop_repeater/repeater.db /root/repeater.db.bak
  ```

- Re-check `tx_power` afterwards — upgrades touch `config.yaml`.
- Upgrading overwrites `/var/lib/openhop_repeater/radio-presets.json` and `radio-settings.json`. Local edits there are lost.
- Database migrations run on first start of the new version and are forward-only.

## Optional: enforce the TX power limit

On a 2 W board, a startup guard stops a bad value ever reaching the radio. A wrong `tx_power` is not applied live — it sits dormant in `config.yaml` until the next restart, then goes out on the first packet. A systemd drop-in catches it at exactly that moment.

Create `/usr/local/sbin/check-txpower`, owned `root:root`, mode `755`:

```sh {linenos=false}
#!/bin/sh
set -u
LIMIT=8
CONFIG="${1:-/etc/openhop_repeater/config.yaml}"
PY=/opt/openhop_repeater/venv/bin/python

val=$("$PY" -c 'import sys,yaml;print(yaml.safe_load(open(sys.argv[1]))["radio"]["tx_power"])' \
      "$CONFIG" 2>/dev/null) || { echo "guard: cannot read radio.tx_power" >&2; exit 1; }

case "$val" in ''|*[!0-9-]*) echo "guard: '$val' not an integer" >&2; exit 1 ;; esac
[ "$val" -gt "$LIMIT" ] && { echo "guard: tx_power=$val exceeds $LIMIT dBm" >&2; exit 1; }
echo "guard: tx_power=$val OK"
```

Then `/etc/systemd/system/openhop-repeater.service.d/10-txpower-guard.conf`:

```ini {linenos=false}
[Unit]
StartLimitIntervalSec=60
StartLimitBurst=3

[Service]
ExecStartPre=/usr/local/sbin/check-txpower
```

Finish with `sudo systemctl daemon-reload`. The daemon now refuses to start if `tx_power` exceeds the limit. A drop-in is used deliberately: upgrades overwrite the main unit file and the Python package, but never `*.service.d/`. Test both paths before relying on it.

## Gotchas

**Raspberry Pi OS keeps the journal in RAM.** It ships `/usr/lib/systemd/journald.conf.d/40-rpi-volatile-storage.conf` with `Storage=volatile` to reduce SD-card wear, so logs vanish on reboot. To make them persistent, add a drop-in whose filename sorts *after* `40-` — drop-ins merge by filename across directories, so a `10-` file is silently ignored:

```bash {linenos=false}
# /etc/systemd/journald.conf.d/99-persistent.conf
[Journal]
Storage=persistent
SystemMaxUse=200M
MaxRetentionSec=1month
```

Changing `Storage` migrates nothing — after restarting journald, run `sudo journalctl --flush`. Weigh this against SD-card wear.

**Duty cycle is not enforced by default.** `duty_cycle.enforcement_enabled` ships as `false`. A repeater retransmits constantly, so at high power this is a thermal consideration as well as a regulatory one.

**Config changes in the web UI are not always applied live.** Frequency, bandwidth, spreading factor, and coding rate take effect immediately. TX power does not — it applies at the next restart. Re-check before restarting.

**Backup/Restore in the Console performs no validation.** Importing a config from another node overwrites the whole `radio:` block, TX power included. Re-check afterwards.

## Troubleshooting

| Symptom | Check |
|---|---|
| Service won't start | `sudo journalctl -u openhop-repeater -n 50` |
| Console shows nothing | Hard-refresh (`Ctrl+Shift+R`); confirm `web.web_path` in `config.yaml` |
| Radio fails to initialise | Board selection and SPI/GPIO pins in the `sx1262:` section |
| No packets | Frequency, bandwidth, spreading factor, and coding rate must match the local mesh — see [Repeater Setup]({{< relref "/docs/meshcore/repeater-setup/index.md" >}}) |
| Wrong TX power after upgrade | `sudo grep -n tx_power /etc/openhop_repeater/config.yaml` |

## References

- openHop Repeater — <https://github.com/openhop-dev/openhop_repeater>
- openHop Console — <https://github.com/Treehouse-00/pymc_console-dist>
- MeshCore firmware — <https://github.com/meshcore-dev/MeshCore>
- NebraHat / ZebraHAT hardware — <https://github.com/wehooper4/Meshtastic-Hardware>
