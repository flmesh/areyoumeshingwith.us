---
title: Regional Settings
date: 2026-07-26T00:00:00-04:00
draft: false
description: Florida's MeshCore modem preset and region scope scheme.
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

Florida runs one statewide modem preset. Region scopes control which repeaters forward a flood.

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

## Region scopes

Scoping limits how far a flood propagates, so statewide airtime does not have to carry traffic that only concerns one county.

Two halves of one mechanism:

- A **scope** is a tag the sender attaches to a flood.
- A **region** is a tag a repeater is willing to forward.

A repeater forwards a scoped flood only when it carries that exact tag. Matching is exact and case-sensitive. The name itself is never transmitted — only a short code derived from it — so two operators who spell a tag differently silently end up on separate regions that never exchange traffic.

Receiving is unaffected. Scoping governs what a repeater *forwards*, not what a companion *hears*.

### Naming

Florida uses a flat `us` → `fl` → county scheme:

```shell {linenos=false}
region put us
region put fl us
region put manatee fl
region save
```

- Lowercase `a-z`, `0-9`, and hyphen; 29 bytes maximum
- 32 regions maximum per repeater
- On firmware below `1.15.0`, follow each `region put` with `region allowf <name>`
- Region filtering requires firmware `1.10.0`+

### The tree does not cascade

`region` prints its tags with indentation, which reads like a hierarchy. It is not one. Each tag is matched independently.

A repeater carrying only `fl` will **not** forward `manatee` traffic. Every tag you intend to forward has to be listed explicitly, from the root down. Treat the output as a flat list of tags that happens to be printed with indentation.

This is the single most common cause of traffic that stops at a county boundary.

### Unscoped traffic

`*` is unscoped traffic. It floods by default and is unaffected by the region list, so adding regions to a repeater never breaks existing traffic.

### Verifying

```shell {linenos=false}
region
```

Every tag you set should be listed, each showing `F` — flood allowed.

Application per node role is in [Node Configuration]({{< relref "/docs/meshcore/configuration/index.md" >}}).
