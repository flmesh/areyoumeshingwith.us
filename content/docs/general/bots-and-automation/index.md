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

Bots, scripts, and AI responders are welcome on the Florida mesh. They are also the easiest way to fill the shared channels with traffic nobody asked for, so this page describes where automation belongs.

<!--more-->

Automation here means anything that transmits without a person pressing send: auto-replies, scripts, alerting integrations, and AI agents.

## Where to run it

- Reply by direct message. A bot that answers the person who called it has no reason to broadcast.
- Run on a channel meant for it. Bots talking to bots, and anyone experimenting with AI on the mesh, belong on the testing channel or a private one. Weather scripts belong on the weather channel.
- Leave the emergency channel clear. It has to stay quiet enough to be useful the day it matters.
- Use a tapback instead of a text reply where the app supports one. It confirms receipt without adding a message.

Two bots answering each other on a shared channel is the failure everyone notices, and it is the one these three habits prevent.

## Public channel traffic

Some automation does belong on the public channel. Emergency and severe weather alerting is the clear case: low volume, high value, and nothing for another bot to reply to.

Coordinate before starting one. Several operators relaying the same feed multiplies the traffic without adding information, and an alert meant for one county reaches the whole state if nothing scopes it. Ask on Discord first.

## Channel names

The everyday channels use the same names on both networks so operators do not have to learn two sets. The lists are in [MeshCore Node Configuration]({{< relref "/docs/meshcore/configuration/index.md" >}}#channels) and [Meshtastic Node Configuration]({{< relref "/docs/meshtastic/configuration/index.md" >}}#channels).

None of this is enforced. Neither network has an operator who could enforce it. It works only as a convention people choose to follow.
