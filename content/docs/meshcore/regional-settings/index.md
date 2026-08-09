---
title: Regional Settings
description: Florida's MeshCore modem preset.
date: 2026-07-26T00:00:00-04:00
draft: false
noindex: false
nav_weight: 2
nav_icon:
  vendor: bs
  name: wifi
  color: orange
authors:
  - beanfield
  - Json_18
series:
  - Docs
---

Florida runs one statewide modem preset.

<!--more-->

## Modem preset

MeshCore's `USA/Canada` preset, with the coding rate raised to `8`.

- Frequency: `910.525 MHz`
- Bandwidth: `62.5 kHz`
- Spreading factor: `7`
- Coding rate: `8`

CLI: `set radio 910.525,62.5,7,8`. In the app, select the preset and then change the coding rate.

### Which values must match

Frequency, bandwidth, and spreading factor must be identical on both ends. A mismatch in any of the three means the radios cannot demodulate each other at all — there is no degraded mode.

Coding rate is the exception. The receiver reads it from the packet header, so a node on `8` and a node on the preset default of `5` still exchange traffic.

### Why coding rate 8

Coding rate is LoRa's forward-error-correction ratio. `5` sends 5 symbols for every 4 of payload (4/5); `8` sends 8 (4/8). The additional parity lets a receiver reconstruct packets that arrive corrupted, which raises the decode margin on weak, noisy, and multipath-affected links — the marginal paths that determine how far a mesh actually reaches.

The cost is airtime. The payload portion of a packet takes 8/5 as long, up to 60% more; total packet time rises by less, since the preamble and header are fixed. On a network of this density that trade favours reliability: a packet that decodes once beats a shorter packet that has to be repeated.

Application per node role is in [Node Configuration]({{< relref "/docs/meshcore/configuration/index.md" >}}).
