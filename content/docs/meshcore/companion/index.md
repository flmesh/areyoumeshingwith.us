---
title: Companion Node Setup
linkTitle: Companion Setup
description: The app settings for a phone-paired MeshCore node on the Florida Mesh. The only page most people need.
date: 2026-07-26T00:00:00-04:00
draft: false
noindex: false
nav_weight: 1
nav_icon:
  vendor: bs
  name: phone
  color: blue
authors:
  - beanfield
  - Json_18
series:
  - Guide
---

If your node is paired to your phone, this page is the whole setup. Three settings and a list of channels.

<!--more-->

{{< notice tip "This is everything you need" >}}
The other pages in this section are for people running a repeater or room server. Region scopes, transmit power, and the command line do not apply to a companion node — you can skip all of it.
{{< /notice >}}

## 1. Radio

Select the `USA/Canada` preset, then change the coding rate to `8`.

| Setting | Value |
|---|---|
| Frequency | `910.525 MHz` |
| Bandwidth | `62.5 kHz` |
| Spreading factor | `7` |
| Coding rate | `8` |

Frequency, bandwidth, and spreading factor have to match the rest of the mesh exactly — if any one of them is wrong your node cannot hear anyone and nobody can hear it. There is no partial or degraded mode.

Coding rate is the exception: it is read from each packet, so a node left on the preset default of `5` still exchanges traffic. Set it to `8` anyway, because it makes your own weak links more reliable.

## 2. Hop Bytes

Settings → Experimental → **Hop Bytes: `2`**

## 3. Channels

`Public` is built in. Add these:

`#wardriving` · `#emergency` · `#hamradio` · `#testing` · `#florida` · `#weather` · `#weekly-mesh-net`

Add them with *Add Channel → Join a Hashtag Channel*, or by scanning a QR code. A hashtag channel derives its key from its name, so everyone who joins by the same name lands on the same channel. These are shared public channels, not private ones.

`#weekly-mesh-net` carries the [Weekly Mesh Net]({{< relref "/docs/weekly-mesh-net/index.md" >}}).

## Check your work

- The radio screen reads `910.525` / `62.5` / `7` / `8`
- Hop Bytes is `2`
- You can see other nodes appearing in your contact list

That's the whole configuration. Nothing else on this site is required to use the mesh.

## Background

**Why coding rate 8.** Coding rate is LoRa's forward-error-correction ratio. `5` sends 5 symbols for every 4 of payload; `8` sends 8. The extra parity lets a receiver rebuild packets that arrive corrupted, which is what decides whether a marginal link works at all. The cost is airtime — the payload portion takes up to 60% longer. On a mesh this size that trade favours reliability.

**Why Hop Bytes 2.** Each hop a packet takes is recorded as a short hash of the node that forwarded it. At 1 byte there are only 256 possible values, so in a mesh this size two repeaters routinely collide and the path becomes ambiguous. Two bytes gives 65,536 values. The cost is that each hop uses two bytes of a fixed-size path field instead of one.
