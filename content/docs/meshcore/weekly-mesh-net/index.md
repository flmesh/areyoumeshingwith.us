---
title: Weekly Mesh Net
linkTitle: Weekly Net
description: Florida Mesh's weekly on-air net on MeshCore.
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

The net runs on a hashtag channel. Add it with *Add Channel → Join a Hashtag Channel*; the key derives from the name.

- Channel: `#weekly-mesh-net`
- Key: `284d8129d937833bdd641f21256dced0`

Radio values are the standard Florida preset. See [Regional Settings]({{< relref "/docs/meshcore/regional-settings/index.md" >}}).

## Check-in format

Post to the net channel, incrementing the counter with each message:

```
(NAME) - (CITY) #FLMeshNet message 1
(NAME) - (CITY) #FLMeshNet message 2
```

Report the messages and acknowledgements you receive back to the same net channel.

## Results

Check-ins are captured by observer nodes feeding the [MeshCore telemetry tools](/telemetry/meshcore/); there is no dedicated net view.

Schedule changes are posted on the [Florida Mesh Discord](https://discord.gg/floridamesh).
