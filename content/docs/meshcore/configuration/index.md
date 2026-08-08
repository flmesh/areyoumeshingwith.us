---
title: Node Configuration
description: Configuring MeshCore companion clients and repeater infrastructure for the Florida Mesh.
nav_weight: 3
authors:
  - beanfield
  - Json_18
series:
  - Guide
date: 2026-07-26T00:00:00-04:00
nav_icon:
  vendor: bs
  name: toggles
  color: grey
---

Companion nodes are configured in the app. Repeaters and room servers are configured over the CLI. Radio values come from [Regional Settings]({{< relref "/docs/meshcore/regional-settings/index.md" >}}).

<!--more-->

## Companion nodes

Radio — select the `USA/Canada` preset, then set the coding rate to `8`:

- Frequency: `910.525 MHz`
- Bandwidth: `62.5 kHz`
- Spreading factor: `7`
- Coding rate: `8`

Channels:

- `Public` — built in
- `#wardriving`, `#emergency`, `#hamradio`, `#testing`, `#florida`, `#weather`, `#weekly-mesh-net`

A hashtag channel derives its key from its name, so every node joining by the same name lands on the same channel. These are shared public channels, not private ones. Add them with *Add Channel → Join a Hashtag Channel* or a QR code; there is no CLI to add channels. `#weekly-mesh-net` carries the [Weekly Mesh Net]({{< relref "/docs/weekly-mesh-net/index.md" >}}).

Settings → Experimental:

- Hop Bytes: `2`

### Why Hop Bytes 2

Every hop a packet takes is recorded in its path as a short hash of the repeater that forwarded it. At 1 byte there are only 256 possible values, so in a mesh this size two repeaters routinely hash to the same byte — routing and telemetry then cannot tell which one actually carried the packet. Two bytes gives 65,536 values and makes collisions rare.

The cost is path length: each hop consumes two bytes of a fixed-size path field instead of one, reducing the maximum number of hops a route can record. Florida accepts that trade.

The CLI equivalent is `set path.hash.mode 1`.

## Repeater & room server nodes

CLI configuration. Reboot to apply radio changes.

```shell {linenos=false}
set radio 910.525,62.5,7,8
set flood.advert.interval 12
set path.hash.mode 1
```

`radio` — frequency, bandwidth, spreading factor, coding rate. The first three must match between two nodes or they cannot demodulate each other; coding rate does not. See [Regional Settings]({{< relref "/docs/meshcore/regional-settings/index.md" >}}).

`flood.advert.interval` — hours between flood adverts, range `3`–`168`, default `12`. Shorter intervals put more advert traffic on the air for every repeater that relays them. `advert.interval` separately controls local zero-hop adverts, in minutes.

`path.hash.mode` — `0` = 1-byte, `1` = 2-byte, `2` = 3-byte. Florida uses 2-byte; see [Why Hop Bytes 2](#why-hop-bytes-2).

Fixed position, for boards without GPS:

```shell {linenos=false}
set lat <decimal degrees>
set lon <decimal degrees>
```

Boards with a GPS module use `gps on`, `gps sync`, and `gps setloc` instead.

Region scopes are set with `region put` and `region save` — see [Regional Settings]({{< relref "/docs/meshcore/regional-settings/index.md" >}}).

## Verifying Your Configuration

```shell {linenos=false}
get radio
region
get tx
get path.hash.mode
get flood.advert.interval
ver
```

- `get radio` returns `910.525,62.5,7,8`
- `region` lists every tag you set, each showing `F` (flood allowed)
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
