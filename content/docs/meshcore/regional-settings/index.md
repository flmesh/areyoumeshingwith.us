---
title: Regional Settings
description: Florida's MeshCore modem preset and region scope scheme.
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

Scoping limits how far a flood propagates, so statewide airtime does not have to carry traffic that only concerns one part of the state.

Two halves of one mechanism:

- A **scope** is a tag attached to a packet by the node that sends it.
- A **region** is a tag a repeater is willing to forward.

A repeater forwards a scoped packet only when it carries that exact tag. Matching is exact and case-sensitive. The name itself is never transmitted — it is hashed into a key, and packets are authenticated against that key — so two operators who spell a tag differently silently end up on separate regions that never exchange traffic.

A name without a leading `#` is treated as if it had one: `us-fl` and `#us-fl` are the same region.

{{< notice note "Not yet ratified" >}}
The scheme below is the form the wider MeshCore community is converging on. It has not been formally agreed, and the Florida subregion list in particular is still open. Expect refinement.
{{< /notice >}}

### The scheme

| Tag | Level | Example |
|---|---|---|
| `us` | Country | `us` |
| `us-southeast` | Meta-region, for traffic crossing state lines | `us-southeast` |
| `us-{state}` | State | `us-fl` |
| `us-{state}-{subregion}` | Part of a state | `us-fl-cfl` |
| `us-{iata}` | Metro area, by airport code | `us-mco` |

Rules:

- Lowercase `a-z`, `0-9`, and hyphen; 29 bytes maximum
- 32 regions maximum per repeater
- Region filtering requires firmware `1.10.0`+
- On firmware below `1.15.0`, follow each `region put` with `region allowf <name>`

`us-fl` is the ISO 3166-2 subdivision code for Florida, lowercased, and `us-ga` / `us-tn` / `us-al` are the same for the states we exchange traffic with. The name costs nothing on air either way — it is hashed to a fixed-size key, so a longer tag does not make a larger packet.

Metro tags use IATA airport codes because observers and mapping tools already standardised on them, and they sit directly under `us` rather than under a state. That is deliberate: a metro's airport is often near a state line and serves more than one — Chattanooga sits at a three-way with north-west Georgia and north-east Alabama, and Charleston's code covers traffic around Savannah — so putting the code behind a state would force duplicate regions like `us-sc-chs` and `us-ga-chs` for one metro.

Subregions are meant to be coarse enough that an operator knows which one they are in without consulting a map — a quadrant of the state, or a metro. Counties are finer than this scheme intends.

### What a repeater carries

Carry the regions you should forward traffic for, and nothing else. Three is typical:

- `us-southeast` — so traffic can cross the state line
- `us-fl` — the state default
- your subregion or metro tag

Do **not** add neighbouring states' tags. `us-southeast` is what carries traffic between states; adding `us-ga` to a Florida repeater pulls all of Georgia's state traffic into Florida airtime.

### What a companion node does

Companions **receive everything**, regardless of region. Scoping never affects what you hear.

Regions only tag what a companion *sends*. Once client support lands, a companion can set a default region for its messages and a region per channel — one region per message, not a list. The intent is a default of `us-fl`, `us-southeast` on channels meant to travel between states, and a narrow local tag on `#testing` or `#wardriving` so that traffic stays put.

### The tree does not cascade

`region` prints its tags with indentation, which reads like a hierarchy. It is not one. Each tag is matched independently.

A repeater carrying only `us-fl` will **not** forward `us-fl-cfl` traffic. Every tag you intend to forward has to be listed explicitly. The hierarchy exists to help people reason about coverage; the repeater only compares transport codes.

This is the single most common cause of traffic that stops at a boundary.

### Unscoped traffic

`*` is unscoped traffic. It floods by default and is unaffected by the region list, so adding regions to a repeater never breaks existing traffic. Nodes that never scope anything keep working exactly as they do now.

### Moving from the earlier tags

Earlier Florida guidance used `us`, `fl`, and a bare county name. `us` is unchanged and stays. `fl` and `us-fl` are *different regions*, not two spellings of one — they hash to different keys and never exchange traffic — and county tags are finer than the scheme now calls for.

Repeaters already deployed should carry both forms until the state has moved over:

```shell {linenos=false}
region put us
region put us-southeast us
region put us-fl us-southeast
region put fl
region save
```

The second argument to `region put` sets the parent shown in the listing. It is presentation only — see above.

### Setting regions on openHop

{{< notice warning "The region CLI is not implemented in openHop Repeater" >}}
Every `region` subcommand — `put`, `save`, `get`, `remove`, `allowf`, `denyf`, `load`, `home` — returns `Error: Region commands not implemented`. The commands above apply to repeaters running MeshCore firmware.
{{< /notice >}}

On [openHop]({{< relref "/docs/meshcore/installation/_index.md" >}}), manage regions in the Console instead; it stores them in the repeater's database rather than in `config.yaml`.

Renaming an existing region there keeps its original key, so it silently carries on matching the old name. Delete the region and add the new one instead of editing it.

### Verifying

```shell {linenos=false}
region
```

Every tag you set should be listed, each showing `F` — flood allowed. On openHop, check the Console's region list instead.

### Further reading

- Background on what regions are and are not — <https://github.com/pinztrek/mesher/blob/main/docs/regions.md>
- Map of regions in use — <https://regions.caboosey.net/>
- Region names and their derived keys — <https://meshmaster.store/regions/>

Application per node role is in [Node Configuration]({{< relref "/docs/meshcore/configuration/index.md" >}}).
