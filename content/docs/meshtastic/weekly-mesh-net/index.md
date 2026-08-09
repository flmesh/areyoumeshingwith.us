---
title: Weekly Mesh Net
linkTitle: Weekly Net
description: Florida Mesh's weekly on-air net on Meshtastic.
date: 2026-07-26T00:00:00-04:00
draft: false
noindex: false
nav_weight: 5
nav_icon:
  vendor: bs
  name: broadcast
  color: green
authors:
  - beanfield
  - Json_18
series:
  - Guide
---

Florida Mesh runs a weekly on-air net to exercise infrastructure, measure range, and put as many operators on the mesh simultaneously.

<!--more-->

## Schedule

Mondays, by region:

- North West Florida: `7:00 PM` local
- North, Central, and South Florida: `8:00 PM` local

Runs two hours.

## Where to check in

The net runs on the default public primary channel. A node already on the mesh is already on it.

- Channel: primary, channel `0`, key `AQ==`, frequency slot `0`
- Modem preset: your region's — see [Regional Settings]({{< relref "/docs/meshtastic/regional-lora-settings/index.md" >}})

## Check-in format

Post to the primary channel, incrementing the counter with each message:

```
(NAME) - (CITY) #FLMeshNet message 1
(NAME) - (CITY) #FLMeshNet message 2
```

Report the messages and acknowledgements you receive back to the same channel.

## Results

Check-ins tagged `#FLMeshNet` and the gateways that heard them are at [meshview.areyoumeshingwith.us/net](https://meshview.areyoumeshingwith.us/net). The map icon beside a message lists the receiving gateways.

Schedule changes are posted on the [Florida Mesh Discord](https://discord.gg/floridamesh).
