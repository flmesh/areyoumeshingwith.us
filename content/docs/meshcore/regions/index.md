---
title: Region Scopes
description: Florida's MeshCore region scope scheme and how scoped flooding works.
date: 2026-07-26T00:00:00-04:00
draft: false
noindex: false
nav_weight: 3
nav_icon:
  vendor: bs
  name: diagram-3
  color: orange
authors:
  - beanfield
  - Json_18
series:
  - Docs
---

Scoping limits how far a flood propagates, so statewide airtime does not have to carry traffic that only concerns one county.

<!--more-->

{{< notice note "Repeater operators" >}}
Region scopes are a repeater setting. They govern what a repeater *forwards*, never what a node *hears*, so a companion node is unaffected by any of this — see [Companion Node Setup]({{< relref "/docs/meshcore/companion/index.md" >}}).
{{< /notice >}}

## How it works

Two halves of one mechanism:

- A **scope** is a tag the sender attaches to a flood.
- A **region** is a tag a repeater is willing to forward.

A repeater forwards a scoped flood only when it carries that exact tag. Matching is exact and case-sensitive. The name itself is never transmitted — only a short code derived from it — so two operators who spell a tag differently silently end up on separate regions that never exchange traffic.

## Naming

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

## The tree does not cascade

`region` prints its tags with indentation, which reads like a hierarchy. It is not one. Each tag is matched independently.

A repeater carrying only `fl` will **not** forward `manatee` traffic. Every tag you intend to forward has to be listed explicitly, from the root down. Treat the output as a flat list of tags that happens to be printed with indentation.

This is the single most common cause of traffic that stops at a county boundary.

## Unscoped traffic

`*` is unscoped traffic. It floods by default and is unaffected by the region list, so adding regions to a repeater never breaks existing traffic.

## Verifying

```shell {linenos=false}
region
```

Every tag you set should be listed, each showing `F` — flood allowed.
