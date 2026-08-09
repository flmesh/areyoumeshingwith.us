---
title: Regional Settings
description: Florida's MeshCore modem preset, and the region scoping work in progress.
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

Florida runs one statewide modem preset. Region scoping is still being worked out.

<!--more-->

## Modem preset

MeshCore's `USA/Canada` preset, with the coding rate raised to `8`.

- Frequency: `910.525 MHz`
- Bandwidth: `62.5 kHz`
- Spreading factor: `7`
- Coding rate: `8`

CLI: `set radio 910.525,62.5,7,8`. In the app, select the preset and then change the coding rate.

### Which values must match

Frequency, bandwidth, and spreading factor must be identical on both ends. A mismatch in any of the three means the radios cannot demodulate each other at all. There is no degraded mode.

Coding rate is the exception. The receiver reads it from the packet header, so a node on `8` and a node on the preset default of `5` still exchange traffic.

### Why coding rate 8

Coding rate is LoRa's forward-error-correction ratio. `5` sends 5 symbols for every 4 of payload (4/5); `8` sends 8 (4/8). The additional parity lets a receiver reconstruct packets that arrive corrupted, which raises the decode margin on weak, noisy, and multipath-affected links, the marginal paths that determine how far a mesh actually reaches.

The cost is airtime. The payload portion of a packet takes 8/5 as long, up to 60% more; total packet time rises by less, since the preamble and header are fixed. On a network of this density that trade favours reliability: a packet that decodes once beats a shorter packet that has to be repeated.

## Region scopes

{{< notice warning "Work in progress" >}}
The naming convention is still being worked out across the wider MeshCore community, and Florida has not chosen its own tags. Expect the formats below to change.
{{< /notice >}}

### How they work

A region is a tag carried on a packet. Every repeater keeps a list of the regions it is willing to relay, and it forwards a scoped packet only when the tag matches something in that list. The name is never sent over the air; it is hashed, and the repeater compares hashes.

Matching is exact and the list is flat. A repeater holding `us-fl` does not forward `us-fl-cfl` traffic just because one name begins with the other. The hierarchy in the names is there to help people reason about coverage, not to route anything.

Unscoped traffic ignores all of this and floods the way it always has, so adding regions to a repeater does not cut anyone off.

### Who configures them

Repeater operators, mostly. A repeater can hold many regions at once, and the guidance so far is to carry only what you should be relaying: your own area, plus a wider tag for the traffic that genuinely needs to travel further.

Companion nodes receive everything regardless of scope, so nothing here affects what you can hear. A companion's region setting only tags what it sends, one region per message, with a default and optionally a per-channel choice once clients support it.

### The formats under discussion

| Format | Example |
|---|---|
| `us` | `us` |
| `us-southeast` | `us-southeast` |
| `us-{state}` | `us-fl` |
| `us-{state}-{subregion}` | `us-fl-cfl` |
| `us-{iata}` | `us-mco` |

The `us-` prefix is there because bare two-letter state codes are ambiguous: `sc` is both South Carolina and Seychelles. `us-fl` is also the ISO 3166-2 code for the state, so the form is already a standard rather than a local invention.

Subregions are meant to be coarse enough that an operator knows which one they are in without consulting a map. IATA codes sit directly under `us` rather than beneath a state, because an airport often sits near a state line and serves across it, and they are expected to be used mainly for wardriving so they line up with Meshmapper.

Florida will need `us-fl` and a subregion list. Neither is decided.

### Further reading

- [Region Filtering](https://blog.meshcore.io/2026/01/20/region-filtering)
- [Default Scope Region](https://blog.meshcore.io/2026/04/17/default-scope)
- [Proposal: predefined standardized region scopes](https://github.com/meshcore-dev/MeshCore/issues/2495)
- [Regions, scopes, and the future of MeshCore](https://github.com/meshcore-dev/MeshCore/discussions/2375)
- [Docs unclear on region behaviour](https://github.com/meshcore-dev/MeshCore/issues/1747)

Application per node role is in [Node Configuration]({{< relref "/docs/meshcore/configuration/index.md" >}}).
