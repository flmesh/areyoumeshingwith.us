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

{{< notice note "Renamed projects" >}}
openHop Repeater was `pyMC_Repeater`; openHop Console was `pyMC Console`. Guides referencing `/etc/pymc_repeater`, the `pymc-repeater` service, or `github.com/dmduran12/…` are out of date. The paths below are current.
{{< /notice >}}

## Requirements

- Raspberry Pi 3, 4, 5, or Zero 2 W
- Raspberry Pi OS Bookworm or later, 32- or 64-bit
- A supported SX1262 or SX1268 LoRa HAT
- Network connection; SPI available, which the installer configures

## 1. Install the repeater

```bash {linenos=false}
git clone https://github.com/openhop-dev/openhop_repeater.git
cd openhop_repeater
sudo bash ./manage.sh install
```

The installer creates a `repeater` service user with hardware access, installs to `/opt/openhop_repeater` in its own venv, creates config, log, and data directories, runs an interactive hardware wizard, and enables the systemd service.

In the wizard, select your exact board. That selection sets your TX power, and the wrong variant gives the wrong power level.

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

Expect exactly one line, under the `radio:` section, at or below the maximum your board's manufacturer specifies. A board with a power amplifier is damaged by a value meant for a bare transceiver.

{{< notice warning >}}
`tx_power` belongs under `radio:`, not under `sx1262:`. A value placed under `sx1262:` is silently ignored and the radio runs at whatever `radio.tx_power` says.
{{< /notice >}}

## 3. Install the Console (optional)

The Console is a static web UI with no backend of its own; the repeater serves it. Install the repeater first.

```bash {linenos=false}
git clone https://github.com/Treehouse-00/pymc_console-dist.git pymc_console
cd pymc_console
sudo bash manage.sh install
```

Browse to `http://<pi-ip>:8000/` and log in with your repeater credentials.

{{< notice note >}}
Follow the repo's `README.md`, not its `INSTALL.md`. The latter is stale and still refers to pre-rename paths and services.
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

The Console does not restart the repeater, because the assets are static. Hard-refresh the browser with `Ctrl+Shift+R`.

The Console's web UI can also upgrade, offering a release channel: `main` for stable snapshots, `dev` for latest commits. `main` is updated infrequently.

### Before you upgrade

- Back up. There is no rollback command.

  ```bash {linenos=false}
  sudo cp -a /etc/openhop_repeater/config.yaml /root/config.yaml.bak
  sudo cp -a /var/lib/openhop_repeater/repeater.db /root/repeater.db.bak
  ```

- Re-check `tx_power` afterwards, since upgrades touch `config.yaml`.
- Upgrading overwrites `/var/lib/openhop_repeater/radio-presets.json` and `radio-settings.json`. Local edits there are lost.
- Database migrations run on first start of the new version and are forward-only.

## References

- openHop Repeater: <https://github.com/openhop-dev/openhop_repeater>
- openHop Console: <https://github.com/Treehouse-00/pymc_console-dist>
- MeshCore firmware: <https://github.com/meshcore-dev/MeshCore>
- NebraHat / ZebraHAT hardware: <https://github.com/wehooper4/Meshtastic-Hardware>
