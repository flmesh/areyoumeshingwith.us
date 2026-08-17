---
title: Bots and Automation
linkTitle: Bots
description: Where bots, scripts, and automated responses belong on the Florida mesh.
date: 2026-08-16T00:00:00-04:00
draft: false
noindex: false
nav_weight: 1000
nav_icon:
  vendor: bs
  name: robot
  color: '#6be'
authors:
  - beanfield
series:
  - Guide
---

Conventions for running automation on the Florida mesh. Automation here means anything that transmits without a person pressing send: auto-replies, scripts, alerting integrations, and AI agents.

<!--more-->

## Airtime is the cost, not screen clutter

Notifications are the visible symptom. The real constraint is airtime: one channel, shared by everyone in range, and every transmission occupies it while repeaters carry it further. Traffic that nobody needed still displaces traffic somebody did.

Automation cuts both ways here. A bot that answers a common question once, by direct message, can lower total traffic by saving five people from asking it on a public channel, and a well-scoped alerting feed can replace a scattered conversation. A bot that acknowledges every message, retries on failure, or answers other bots spends the same airtime and returns nothing. The difference is implementation, not intent.

## Where to run it

- **Reply by direct message.** A bot answering the node that called it does not need to broadcast.
- **Use a channel meant for it.** Bot-to-bot traffic and AI experiments go on the testing channel or a private one. Weather scripts go on the weather channel.
- **Keep the emergency channel clear.** Emergency traffic only.
- **Prefer a tapback to a text reply** where the app supports one. It acknowledges without a second message.

An auto-reply that triggers on another bot's auto-reply loops until one of them stops. Reply by DM, or filter on sender, or both.

## Consider a room server instead

If the content is something people fetch rather than something they need pushed to them, a MeshCore room server or a bulletin board is often the better mechanism. Messages sit on the server until a client connects and reads them, so the cost is one exchange with one node at a time the reader chose, rather than a broadcast every radio in range has to carry whether or not anyone wanted it.

Bulletins, net announcements, reference material, and logs all fit that shape. Time-critical alerts do not, which is why emergency and severe weather alerting stays on a channel.

## Label it as a bot

Put `bot` in the node name, and in the short name where the app has one. A reader scanning a channel can then separate scripted traffic from human traffic without knowing the nodes.

## Public channel traffic

Emergency and severe weather alerting is the case that justifies the public channel: low volume, high value, no reply loop.

Coordinate before running one. Duplicate relays of the same feed multiply traffic without adding information, and an alert scoped to one county floods the state unless something limits it. Ask on Discord first.

## Channel names

The everyday channels use the same names on both networks. The lists, keys, and add links are in [MeshCore Node Configuration]({{< relref "/docs/meshcore/configuration/index.md" >}}#channels) and [Meshtastic Node Configuration]({{< relref "/docs/meshtastic/configuration/index.md" >}}#channels).
