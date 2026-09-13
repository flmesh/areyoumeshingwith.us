---
title: Node Configuration
linkTitle: Configuration
description: Configuring MeshCore companion clients and repeater infrastructure for the Florida Mesh.
nav_weight: 3
authors:
  - beanfield
  - Json_18
  - meh098
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

## Firmware

Flash from the [MeshCore flasher](https://flasher.meshcore.io) over USB, using Chrome or Edge. Select the board, then the firmware role: companion, repeater, or room server.

{{< notice note "Use meshcore.io" >}}
The project moved to `meshcore.io`. `meshcore.co.uk` is no longer affiliated with the official project.
{{< /notice >}}

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

Every hop a packet takes is recorded in its path as a short hash of the repeater that forwarded it. At 1 byte there are only 256 possible values, so in a mesh this size two repeaters routinely hash to the same byte, after which routing and telemetry cannot tell which one actually carried the packet. Two bytes gives 65,536 values and makes collisions rare.

The cost is path length: each hop consumes two bytes of a fixed-size path field instead of one, reducing the maximum number of hops a route can record. Florida accepts that trade.

The CLI equivalent is `set path.hash.mode 1`.

## Channels

Florida uses:

| Channel | Key | Region Scope | Purpose |
|---|---|---|---|
| [`Public`](meshcore://channel/add?name=Public&secret=8b3387e9c5cdea6ac9e5edbaa115cd72) | Built in | | General traffic |
| [`FloridaMesh`](meshcore://channel/add?name=FloridaMesh&secret=e9a7128005b364ae010dd6330d693fa6) | `e9a7128005b364ae010dd6330d693fa6` | | Statewide chat |
| [`#emergency`](meshcore://channel/add?name=%23emergency&secret=e1ad578d25108e344808f30dfdaaf926) | Derived from name | | Emergency and severe weather traffic |
| [`#weather`](meshcore://channel/add?name=%23weather&secret=88f502554fee92a1625cfb311546e7cb) | Derived from name | | Weather scripts and alerting |
| [`#testing`](meshcore://channel/add?name=%23testing&secret=cde5e82cf515647dcb547a79a4f065d1) | Derived from name | | Node, script, and bot testing |
| [`#weekly-mesh-net`](meshcore://channel/add?name=%23weekly-mesh-net&secret=284d8129d937833bdd641f21256dced0) | Derived from name | | [Weekly Mesh Net]({{< relref "/docs/general/weekly-mesh-net/index.md" >}}) |

Open this page in your phone's browser and tap a channel name to hand it straight to the MeshCore app, which opens with the name and key already filled in. To add one by hand instead: `Public` is built in and needs nothing entered, the hashtag channels derive their key from their name, so *Add Channel → Join a Hashtag Channel* or a QR code is enough, and `FloridaMesh` is private, so enter its name and key exactly as listed. There is no CLI to add channels.

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

`path.hash.mode` selects the path hash size: `0` = 1-byte, `1` = 2-byte, `2` = 3-byte. Florida uses 2-byte; see [Hop Bytes](#hop-bytes).

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
