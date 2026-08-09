---
title: Repeater Setup
description: CLI configuration for MeshCore repeaters and room servers on the Florida Mesh.
date: 2026-07-26T00:00:00-04:00
draft: false
noindex: false
nav_weight: 2
nav_icon:
  vendor: bs
  name: toggles
  color: grey
authors:
  - beanfield
  - Json_18
series:
  - Guide
---

Repeaters and room servers are configured over the CLI. Reboot to apply radio changes.

<!--more-->

{{< notice note "Repeater operators" >}}
This page is for infrastructure nodes. If your node is paired to a phone, everything you need is in [Companion Node Setup]({{< relref "/docs/meshcore/companion/index.md" >}}) — none of the settings here apply to you.
{{< /notice >}}

## Settings

```shell {linenos=false}
set radio 910.525,62.5,7,8
set flood.advert.interval 12
set path.hash.mode 1
```

`radio` — frequency, bandwidth, spreading factor, coding rate. The first three must match between two nodes or they cannot demodulate each other at all; coding rate does not have to match, because the receiver reads it from the packet header.

`flood.advert.interval` — hours between flood adverts, range `3`–`168`, default `12`. Shorter intervals put more advert traffic on the air for every repeater that relays them. `advert.interval` separately controls local zero-hop adverts, in minutes.

`path.hash.mode` — `0` = 1-byte, `1` = 2-byte, `2` = 3-byte. Florida uses 2-byte: at 1 byte there are only 256 possible values, so in a mesh this size two repeaters routinely hash to the same byte and the recorded path becomes ambiguous. The cost is that each hop consumes two bytes of a fixed-size path field.

Fixed position, for boards without GPS:

```shell {linenos=false}
set lat <decimal degrees>
set lon <decimal degrees>
```

Boards with a GPS module use `gps on`, `gps sync`, and `gps setloc` instead.

## Region scopes

Which floods this repeater forwards beyond its own area — see [Region Scopes]({{< relref "/docs/meshcore/regions/index.md" >}}).

## Verifying

```shell {linenos=false}
get radio
get tx
get path.hash.mode
get flood.advert.interval
ver
```

- `get radio` returns `910.525,62.5,7,8`
- `ver` reports firmware `1.10.0`+ (required for region filtering), `1.16.0`+ recommended

## Troubleshooting

### Repeater hears very little, or stops hearing traffic entirely

Two receive-side settings, neither of which is set by the modem preset.

```shell {linenos=false}
set radio.rxgain on
set agc.reset.interval 12
```

`radio.rxgain` enables the LoRa transceiver's boosted receive gain — a low-noise front-end mode on SX12xx and LR1110 radios that trades a small increase in current draw for roughly 3 dB of sensitivity. It defaults to `on` in firmware `1.14.1`+, but confirm it with `get radio.rxgain` on any repeater that underperforms.

`agc.reset.interval` periodically restarts the receiver's automatic gain control. A strong nearby transmission can leave the AGC latched at low gain, after which the radio still appears healthy but no longer decodes weak signals — the classic "deaf repeater". The value is in seconds, rounded down to a multiple of 4, and defaults to `0` (disabled). `12` restores sensitivity on an unattended repeater without resetting so often that it clips traffic mid-reception.

### Poor range, or the node is not being heard

```shell {linenos=false}
set tx <dBm>
```

Transmit power, range `1`–`22` dBm. The default varies by board.

This value sets the transceiver's output, not the antenna's. On boards with a power amplifier the radiated power is substantially higher than the number you set. Establish your board's PA gain and your applicable limits before raising it. Driving a PA beyond its rating, or transmitting into a bad or missing antenna, will permanently damage the radio.

Raising TX power does not fix an asymmetric link. If you can hear a repeater but it cannot hear you, the deficit is usually antenna, feedline, or siting — not power.
