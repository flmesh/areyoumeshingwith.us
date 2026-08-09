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

Scoping limits how far a flood propagates, so statewide airtime does not have to carry traffic that only concerns one county.

Region scopes are a repeater setting. They govern what a repeater *forwards*, never what a node *hears*, so a companion node is unaffected by any of this.

Two halves of one mechanism:

- A **scope** is a tag the sender attaches to a flood.
- A **region** is a tag a repeater is willing to forward.

A repeater forwards a scoped flood only when it carries that exact tag. Matching is exact and case-sensitive. The name itself is never transmitted — it is hashed into a key, and packets are authenticated against that key — so two operators who spell a tag differently silently end up on separate regions that never exchange traffic.

A name without a leading `#` is treated as if it had one: `us-fl` and `#us-fl` are the same region.

### Florida's scheme

Three levels, each a complete hyphenated tag:

| Tag | Covers |
|---|---|
| `us-southeast` | Interstate traffic between Florida and its neighbours |
| `us-fl` | Statewide |
| `us-fl-<county>` | One county, e.g. `us-fl-manatee` |

Rules:

- Lowercase `a-z`, `0-9`, and hyphen; 29 bytes maximum
- 32 regions maximum per repeater
- Region filtering requires firmware `1.10.0`+
- On firmware below `1.15.0`, follow each `region put` with `region allowf <name>`

### Why the `us-` prefix

`us-fl` is the ISO 3166-2 subdivision code for Florida, lowercased, and `us-ga` / `us-nc` / `us-sc` are the same for our neighbours. Georgia and western North Carolina are already implementing that form, so it is both the standard and what the adjacent states will actually match against.

Bare two-letter state codes are ambiguous internationally — `sc` is Seychelles, `ga` is Gabon, `nc` is New Caledonia, `tn` is Tunisia. Florida is one of the states that happens *not* to collide with a country code, so `fl` was never dangerous on its own; the prefix is about matching the region, not avoiding a clash.

The name costs nothing on air either way. It is hashed to a fixed-size key, so a longer tag does not make a larger packet.

### Why county names, not airport codes

Some networks use IATA codes for local scopes. Florida uses county names because IATA only covers places with an airport — most rural counties have none — and because the area a code implies is ambiguous: `srq` is "Sarasota–Bradenton", which is two counties. Metro repeaters may carry an IATA tag *in addition* if it is locally useful, but the county tag is the one that must be there.

### Interstate traffic

`us-southeast` is the tag Georgia and western North Carolina carry for traffic that crosses a state line. A Florida repeater without it will not forward their scoped floods, and they will not forward Florida's.

Every Florida repeater can carry it; boundary and high-site repeaters in particular should. The recommendation from the operators in the adjacent states is to carry `us-southeast` rather than adding each neighbour's state tag to your repeaters.

### During the transition

Earlier Florida guidance used bare tags — `us`, `fl`, and a bare county name. Those are *different regions*, not shorthand for the new ones: `fl` and `us-fl` hash to different keys and never exchange traffic.

Repeaters already deployed should carry both forms until the state has moved over:

```shell {linenos=false}
region put us-southeast
region put us-fl us-southeast
region put us-fl-manatee us-fl
region put fl
region put manatee fl
region save
```

Drop the bare `us` tag. No network outside Florida scopes a flood to the whole country, so it matches nothing.

The second argument to `region put` sets the parent shown in the listing. It is presentation only — see below.

### The tree does not cascade

`region` prints its tags with indentation, which reads like a hierarchy. It is not one. Each tag is matched independently.

A repeater carrying only `us-fl` will **not** forward `us-fl-manatee` traffic. Every tag you intend to forward has to be listed explicitly, from the root down. Treat the output as a flat list of tags that happens to be printed with indentation.

This is the single most common cause of traffic that stops at a county boundary.

### Unscoped traffic

`*` is unscoped traffic. It floods by default and is unaffected by the region list, so adding regions to a repeater never breaks existing traffic.

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

Region names in use across the wider network, and their derived keys, are collected at <https://meshmaster.store/regions/>.

Application per node role is in [Node Configuration]({{< relref "/docs/meshcore/configuration/index.md" >}}).
