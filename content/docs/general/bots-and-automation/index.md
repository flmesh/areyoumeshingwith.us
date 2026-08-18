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

Automation earns its airtime in plenty of places: severe weather and emergency alerting, range validation, equipment testing, new user onboarding, and proximity alerts. Not all of them belong on every channel. Remember, airtime is shared and finite: one frequency for every node in range, used again by each repeater that forwards a packet. A bot that answers by direct message saves airtime. One that replies to everything, retries failures, or answers other bots wastes it and creates noise.

## Best Practices

- **Keep the emergency channels clear.** Emergency traffic only.
- Add the word `bot` to your node name to make it easy for users and other automated responders to identify.
- **Reply by direct message, not the channel.** A bot answering the node that called it does not need to broadcast. Direct messages are also cheaper on the air: both networks learn a path to the destination and route later messages along it rather than flooding the mesh, [Meshtastic](https://meshtastic.org/docs/overview/mesh-algo/) with next-hop routing and [MeshCore](https://docs.meshcore.io/faq/#54-q-how-does-a-node-discover-a-path-to-its-destination-and-then-use-it-to-send-messages-in-the-future-instead-of-flooding-every-message-it-sends-like-meshtastic) with the path recorded in the delivery report.
- **Use tapbacks** where the app supports one. It acknowledges without a second message.
- **Avoid replying on the default public channel.** Use one of the above two methods or choose an appropriate [MeshCore]({{< relref "/docs/meshcore/configuration/index.md" >}}#channels) or [Meshtastic]({{< relref "/docs/meshtastic/configuration/index.md" >}}#channels) channel.
- Keep replies concise and pertinent.
- Consider your reach. Local weather forecasts and the functional status of your automated responder do not belong on an MQTT uplinked channel that reaches nodes across the state.
- When in doubt, start small and on a private channel.
- If the content is something people fetch rather than something they need pushed to them, consider a BBS or Room Server.

## Projects worth looking at

- [Mesh Monitor](https://meshmonitor.org). Self-hosted web dashboard for Meshtastic, MeshCore, and MQTT sources in one deployment: maps, messaging, telemetry, and remote device administration. Its automation covers auto-responders, scheduled messages, auto-traceroute, geofence triggers, and push notifications, extensible with Python or Bash scripts.
- [meshing-around](https://github.com/SpudGunMan/meshing-around). Python bot suite for Meshtastic, run over serial, TCP, or BLE across multiple nodes. Keyword responder, BBS with mail and store-and-forward, scheduled messages, NOAA and USGS weather, river, tide and earthquake lookups, EAS alerts, wiki and satellite pass queries, proximity alerts, network and hardware test commands, optional LLM integration, and games.
- [Meshtastic SAME EAS Alerter](https://github.com/RCGV1/Meshtastic-SAME-EAS-Alerter). Decodes SAME emergency alerts off the air with an RTL-SDR and forwards them to a node over serial or TCP. No internet connection required.
