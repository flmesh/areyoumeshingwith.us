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

If your node is paired to your phone, this section plus [Channels](#channels) is the whole setup. Nothing in the repeater section applies to you.

### Radio

Select the `USA/Canada` preset, then change the coding rate to `8`:

| Setting | Value |
|---|---|
| Frequency | `910.525 MHz` |
| Bandwidth | `62.5 kHz` |
| Spreading factor | `7` |
| Coding rate | `8` |

Frequency, bandwidth, and spreading factor must match the rest of the mesh exactly. Coding rate does not have to, since the receiver reads it from each packet, but set it to `8` anyway because it makes your own weak links more reliable.

### Hop Bytes

Settings → Experimental → **Hop Bytes: `2`**

### Check your work

- The radio screen reads `910.525` / `62.5` / `7` / `8`
- Hop Bytes is `2`
- The channels below are joined
- Other nodes appear in your contact list

That is the whole radio configuration for a companion node. Channels are next.

### Why Hop Bytes 2

Every hop a packet takes is recorded in its path as a short hash of the repeater that forwarded it. At 1 byte there are only 256 possible values, so in a mesh this size two repeaters routinely hash to the same byte, after which routing and telemetry cannot tell which one actually carried the packet. Two bytes gives 65,536 values and makes collisions rare.

The cost is path length: each hop consumes two bytes of a fixed-size path field instead of one, reducing the maximum number of hops a route can record. Florida accepts that trade.

The CLI equivalent is `set path.hash.mode 1`.

## Channels

`Public` is built in. Florida also uses:

`#emergency` · `#hamradio` · `#testing` · `#florida` · `#weather` · `#weekly-mesh-net`

Add them with *Add Channel → Join a Hashtag Channel* or a QR code; there is no CLI to add channels. A hashtag channel derives its key from its name, so every node joining by the same name lands on the same channel. These are shared public channels, not private ones. `#weekly-mesh-net` carries the [Weekly Mesh Net]({{< relref "/docs/meshcore/weekly-mesh-net/index.md" >}}).

Join the ones you intend to use. There is no cost to joining a channel you rarely read, and traffic on it reaches you either way.

If you run Meshmapper, list every channel you know of under **Public Channels** in its admin settings. That is what the wardriving app listens on for passive RX logs, and a longer list means more coverage recorded for the same drive.

## Region scopes

Region scopes decide which repeaters relay a flood beyond its own area. They are set on repeaters, and companion nodes receive everything regardless, so most operators have nothing to configure here yet.

Florida has not chosen its tags and the naming convention is still moving. See [Region Scopes]({{< relref "/docs/meshcore/regional-settings/index.md" >}}#region-scopes) for how they work, the formats under discussion, and the CLI commands.

## Repeater & room server nodes

CLI configuration. Reboot to apply radio changes.

```shell {linenos=false}
set radio 910.525,62.5,7,8
set flood.advert.interval 12
set path.hash.mode 1
```

`radio` sets frequency, bandwidth, spreading factor, and coding rate. The first three must match between two nodes or they cannot demodulate each other; coding rate does not. See [Regional Settings]({{< relref "/docs/meshcore/regional-settings/index.md" >}}).

`flood.advert.interval` is the number of hours between flood adverts, range `3`–`168`, default `12`. Shorter intervals put more advert traffic on the air for every repeater that relays them. `advert.interval` separately controls local zero-hop adverts, in minutes.

`path.hash.mode` selects the path hash size: `0` = 1-byte, `1` = 2-byte, `2` = 3-byte. Florida uses 2-byte; see [Why Hop Bytes 2](#why-hop-bytes-2).

Fixed position, for boards without GPS:

```shell {linenos=false}
set lat <decimal degrees>
set lon <decimal degrees>
```

Boards with a GPS module use `gps on`, `gps sync`, and `gps setloc` instead.

### Verifying Your Configuration

```shell {linenos=false}
get radio
get tx
get path.hash.mode
get flood.advert.interval
ver
```

- `get radio` returns `910.525,62.5,7,8`
- `ver` reports the firmware version, so you can confirm the node is running the one you intended

### Troubleshooting

#### Repeater hears very little, or stops hearing traffic entirely

Two receive-side settings, neither of which is set by the modem preset.

```shell {linenos=false}
set radio.rxgain on
set agc.reset.interval 12
```

`radio.rxgain` enables the LoRa transceiver's boosted receive gain, a low-noise front-end mode on SX12xx and LR1110 radios that trades a small increase in current draw for roughly 3 dB of sensitivity. It defaults to `on` in firmware `1.14.1`+, but confirm it with `get radio.rxgain` on any repeater that underperforms.

`agc.reset.interval` periodically restarts the receiver's automatic gain control. A strong nearby transmission can leave the AGC latched at low gain, after which the radio still appears healthy but no longer decodes weak signals, the classic "deaf repeater". The value is in seconds, rounded down to a multiple of 4, and defaults to `0` (disabled). `12` restores sensitivity on an unattended repeater without resetting so often that it clips traffic mid-reception.

#### Poor range, or the node is not being heard

```shell {linenos=false}
set tx <dBm>
```

Transmit power, range `1`–`22` dBm. The default varies by board.

This value sets the transceiver's output, not the antenna's. On boards with a power amplifier the radiated power is substantially higher than the number you set. Establish your board's PA gain and your applicable limits before raising it. Driving a PA beyond its rating, or transmitting into a bad or missing antenna, will permanently damage the radio.

Raising TX power does not fix an asymmetric link. If you can hear a repeater but it cannot hear you, the deficit is usually antenna, feedline, or siting, not power.
