---
title: "MQTT Server"
# linkTitle:
date: 2026-04-27T11:05:52-04:00
description: The Florida Mesh MQTT server supports individual account credentials.
noindex: false
# comments: false
nav_weight: 1000
authors:
  - jbouse
nav_icon:
  vendor: bootstrap
  name: hdd-network
  color: '#e24d0e'
series:
  - Docs
categories:
#  - 
tags:
#  - 
images:
#  - 
---

This is meant to enable Florida Mesh to provide better stability for the MQTT server and easier to identify responsible
parties if there an issue pops up that is not possible using shared credentials as originally launched. At a date in the
future, currently still to be determined, the shared credentials will be deprecated and removed completely but for now
only allow channel `uplink` and deny `downlink`.

{{< notice warning >}}
If you are trying to connect to the Florida Mesh MQTT server with anything other than a Meshtastic node, you **must** create
an account.

This includes but is not limited to: a private MQTT Broker/Bridge, Node Red, or any custom application.

Shared credentials do not work for non-nodes at all now.
{{< /notice >}}

To facilitate this, Florida Mesh is utilizing a custom Discord server bot, `@MeshBot`, that can accept commands in
the [#mqtt](https://discord.com/channels/1235735441592942602/1372379064106483843) channel. Once MeshBot has direct
messaged (DM) you, commands can then be sent back via DM rather than in the channel and do not require you to
`@`-mention the bot by name.

You must allow direct messages from server members in Discord before requesting an account. If this privacy setting is disabled, MeshBot will not be able to DM you.

{{< figure src="allow-direct-message.png" alt="Discord User Settings" width="75%" >}}

## Requesting an MQTT Account

Use MeshBot to request a new MQTT account. Mention `@MeshBot` and run the `mqtt.request` command with the `username` argument set to the MQTT username you want to use.

Example:

{{< notice tip "In #mqtt channel" >}}
@MeshBot mqtt.request username:your_username
{{< /notice >}}

{{< figure src="mqtt-request.png" alt="Request MQTT account from MeshBot" >}}

After the account is created, MeshBot will DM the MQTT credentials to you. These credentials are used with `mqtt.areyoumeshingwith.us`. Keep them private and only use them on nodes you operate.

{{< figure src="mqtt-request-dm.png" alt="MQTT account credentials DM from MeshBot" >}}

You can check the MQTT account currently linked to your Discord account with:

{{< notice tip "In #mqtt channel" >}}
@MeshBot mqtt.my-account
{{< /notice >}}

{{< notice tip "In DM with @MeshBot" >}}
mqtt.my-account
{{< /notice >}}

{{< figure src="mqtt-myaccount.png" alt="MQTT my-account command over DM with MeshBot" >}}

If your node needs a different MQTT policy profile, request that change from the Admins. Include the account username, node ID, and a short explanation of what access you need changed.

## Resetting Your MQTT Password

Use MeshBot in the Discord server channel to rotate your MQTT account password:

{{< notice tip "In #mqtt channel" >}}
@MeshBot mqtt.rotate
{{< /notice >}}

{{< notice tip "In DM with @MeshBot" >}}
mqtt.rotate
{{< /notice >}}

{{< figure src="mqtt-rotate.png" alt="MQTT rotate command over DM with MeshBot" >}}

{{< notice tip >}}
Some MeshBot commands like `mqtt.rotate` that make changes require confirmation before the action to be taken is performed.
{{< /notice >}}

MeshBot will DM the replacement credentials to you. After receiving them, update the MQTT password on every node using that account. Old passwords should be considered inactive after a rotation.

If you cannot rotate your own password, ask an Admin to reset it for you.
