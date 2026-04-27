---
title: "MeshBot Commands"
# linkTitle:
date: 2026-04-27T00:00:00-04:00
description: Public MeshBot command usage for the Florida Mesh Discord server.
noindex: false
# comments: false
nav_weight: 1001
authors:
  - jbouse
nav_icon:
  vendor: bootstrap
  name: robot
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

MeshBot is the Florida Mesh Discord bot used for MQTT account self-service and mesh troubleshooting commands.

<!--more-->

Commands can be sent in a Discord server channel by mentioning `@MeshBot` first:

{{< notice tip "In a server channel" >}}
@MeshBot help
{{< /notice >}}

After MeshBot has direct messaged you, commands can also be sent directly to the bot without the `@MeshBot` mention:

{{< notice tip "In DM with @MeshBot" >}}
help
{{< /notice >}}

You must allow direct messages from server members in Discord if you expect MeshBot to send credentials or command results by DM.

## Help

Use `help` to list available commands. Use `--help` with a command to show command-specific usage.

{{< notice tip >}}
@MeshBot help

@MeshBot mqtt.request --help

@MeshBot node.trace --help
{{< /notice >}}

Some commands also accept `--argument value` syntax instead of `argument:value` syntax.

## Basic Commands

### ping

Checks whether MeshBot is responding.

{{< notice tip >}}
@MeshBot ping
{{< /notice >}}

### info

Shows non-sensitive runtime information for MeshBot, including bot name, version, adapter, Node.js version, and uptime.

{{< notice tip >}}
@MeshBot info
{{< /notice >}}

### where am i

Shows the Discord server and channel MeshBot sees for the current message. In a DM, it reports that the conversation is a DM.

{{< notice tip >}}
@MeshBot where am i
{{< /notice >}}

## MQTT Account Commands

### mqtt.request

Requests a new MQTT account linked to your Discord account. The `username` argument is required.

{{< notice tip >}}
@MeshBot mqtt.request username:your_username

@MeshBot mqtt.request --username your_username
{{< /notice >}}

MeshBot validates the requested username, creates the account with the default MQTT profile, and DMs the generated credentials.

Common aliases:

- `mqtt request`
- `request mqtt`

### mqtt.my-account

Shows the MQTT account linked to your Discord account, including username, status, profile, and creation time.

{{< notice tip >}}
@MeshBot mqtt.my-account
{{< /notice >}}

Common aliases:

- `mqtt my-account`
- `mqtt account`
- `my mqtt account`

### mqtt.rotate

Rotates the password for your MQTT account and DMs the replacement credentials.

{{< notice tip >}}
@MeshBot mqtt.rotate
{{< /notice >}}

This command requires confirmation before MeshBot changes the password. After rotating the password, update every node or service using that MQTT account.

Common aliases:

- `mqtt rotate`
- `rotate mqtt password`

## Node Troubleshooting Commands

### node.logs

Queries MQTT broker logs for a Meshtastic node ID.

{{< notice tip >}}
@MeshBot node.logs clientid:a1b2c3d4

@MeshBot node.logs clientid:!a1b2c3d4

@MeshBot node.logs clientid:a1b2c3d4 minutes:30

@MeshBot node.logs clientid:a1b2c3d4 minutes:30 limit:10
{{< /notice >}}

Arguments:

- `clientid` is required. Use an 8-character hexadecimal node ID, optionally prefixed with `!`, or the unsigned decimal equivalent.
- `minutes` is optional and defaults to `15`.
- `limit` is optional and defaults to `20`. Results are capped at `100`.

Common aliases:

- `node logs`
- `show node logs`
- `logs for node`

### node.trace

Queries Floodgate message events by sender, receiver, or Floodgate message ID.

{{< notice tip >}}
@MeshBot node.trace from:a1b2c3d4

@MeshBot node.trace to:!a1b2c3d4

@MeshBot node.trace from:a1b2c3d4 to:9e9fba3c minutes:30

@MeshBot node.trace id:1234567890

@MeshBot node.trace id:1234567890 minutes:30 limit:10
{{< /notice >}}

Arguments:

- Provide either `id`, `from`, or `to`.
- `from` and `to` can be combined to narrow the trace.
- `id` cannot be combined with `from` or `to`.
- Node IDs can be 8-character hexadecimal values, optionally prefixed with `!`, the unsigned decimal equivalent, or `all` for broadcast.
- `minutes` is optional and defaults to `15`.
- `limit` is optional and defaults to `20`. Results are capped at `100`.

Common aliases:

- `node trace`
- `trace node`

## Admin-Only Commands

MeshBot also has admin-only commands for MQTT account resets, profile management, bans, and AUTHZ reports. Those commands require a configured Discord admin role and are not public self-service commands.
