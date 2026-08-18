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

{{< notice note "Work in progress" >}}
The naming convention is still being worked out across the wider MeshCore community, and Florida has not chosen its own tags. Expect the formats below to change.
{{< /notice >}}

### How they work

A region is a tag carried on a packet. Every repeater keeps a list of the regions it is willing to relay, and it forwards a scoped packet only when the tag matches something in that list. The name is never sent over the air; it is hashed, and the repeater compares hashes.

Matching is exact and the list is flat. A repeater holding `us-fl` does not forward `us-fl-cfl` traffic just because one name begins with the other. The hierarchy in the names is there to help people reason about coverage, not to route anything.

Unscoped traffic ignores all of this and floods the way it always has, so adding regions to a repeater does not cut anyone off.

### Who configures them

Repeater operators, mostly. A repeater can hold many regions at once, and the guidance so far is to carry only what you should be relaying: your own area, plus a wider tag for the traffic that genuinely needs to travel further.

Companion nodes receive everything regardless of scope, so nothing here affects what you can hear. A companion's region setting only tags what it sends. The app's experimental settings hold a default scope that applies to every flood packet the node originates, and a channel configured with its own scope overrides that default.

### The formats under discussion

| Format | Example |
|---|---|
| `us` | `us` |
| `us-southeast` | `us-southeast` |
| `us-{state}` | `us-fl` |
| `us-{state}-{subregion}` | `us-fl-cfl` |
| `us-{iata}` | `us-mco` |

`us-fl` is the ISO 3166-2 subdivision code for Florida, so the prefixed form is an existing standard rather than a local invention. A bare `fl` is only an abbreviation, and two-letter codes on their own stop being unique once a mesh reaches beyond one country.

Subregions are meant to be coarse enough that an operator knows which one they are in without consulting a map. IATA codes sit directly under `us` rather than beneath a state, because an airport near a state line serves operators on both sides of it and a state-prefixed code would split one metro into two regions. They are expected to be used mainly for wardriving, so they line up with Meshmapper.

Florida will need `us-fl` and a subregion list. Neither is decided.

### Setting them over the CLI

Regions are configured on the repeater, not in the app. Add each one you intend to relay, then save:

```shell {linenos=false}
region put us-fl
region put us-fl-cfl us-fl
region save
```

The second argument sets the parent shown in the listing. It affects the display only, since each name is hashed whole and matched on its own.

List what the node is carrying, and remove one you no longer want:

```shell {linenos=false}
region
region remove us-fl-cfl
```

Each entry should show `F`, meaning flood allowed. On firmware older than `1.15.0`, follow each `region put` with `region allowf <name>`.

`region default <name>` sets the scope the node puts on packets it originates itself, and `region default <null>` clears it.

{{< notice note "openHop repeaters" >}}
openHop Repeater does not implement any of these. Every `region` subcommand returns `Error: Region commands not implemented`, so manage regions in the openHop Console instead. Renaming an entry there keeps its original key, so delete it and add the new name rather than editing it.
{{< /notice >}}

Application per node role is in [Node Configuration]({{< relref "/docs/meshcore/configuration/index.md" >}}).
