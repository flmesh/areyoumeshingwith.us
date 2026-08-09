---
title: Weekly Mesh Net
linkTitle: Weekly Net
description: Florida Mesh's weekly on-air net, run concurrently on Meshtastic and MeshCore.
date: 2026-07-26T00:00:00-04:00
draft: false
noindex: false
nav_weight: 2
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

Florida Mesh runs a weekly on-air net to exercise infrastructure, measure range, and put as many operators on the mesh simultaneously. It runs on both networks at the same time.

<!--more-->

## Schedule

Mondays, by region:

- North West Florida: `7:00 PM` local
- North, Central, and South Florida: `8:00 PM` local

Runs two hours.

## Where to check in

### Meshtastic

The net runs on the default public primary channel. A node already on the mesh is already on it.

- Channel: primary, channel `0`, key `AQ==`, frequency slot `0`
- Modem preset: your region's — see [Regional Settings]({{< relref "/docs/meshtastic/regional-lora-settings/index.md" >}})

### MeshCore

The net runs on a hashtag channel. Add it with *Add Channel → Join a Hashtag Channel*; the key derives from the name.

- Channel: `#weekly-mesh-net`
- Key: `284d8129d937833bdd641f21256dced0`

Radio values are the standard Florida preset — see [Regional Settings]({{< relref "/docs/meshcore/regional-settings/index.md" >}}).

## Check-in format

Post to the net channel on your network, incrementing the counter with each message:

```
(NAME) - (CITY) #FLMeshNet message 1
(NAME) - (CITY) #FLMeshNet message 2
```

Report the messages and acknowledgements you receive back to the same net channel.

## Results

Meshtastic check-ins tagged `#FLMeshNet` and the gateways that heard them are at [meshview.areyoumeshingwith.us/net](https://meshview.areyoumeshingwith.us/net). The map icon beside a message lists the receiving gateways.

MeshCore check-ins are captured by observer nodes feeding the [MeshCore telemetry tools](/telemetry/meshcore/); there is no dedicated net view.

Schedule changes are posted on the [Florida Mesh Discord](https://discord.gg/floridamesh).
