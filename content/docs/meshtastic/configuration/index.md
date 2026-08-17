---
title: Node Configuration
linkTitle: Configuration
description: Recommendations for configuring your node and contributing to the mesh.
nav_weight: 3
authors:
  - beanfield
  - Json_18
  - PockyBum522
  - jbouse
series:
  - Guide
date: 2026-01-19T00:31:34-04:00
nav_icon:
  vendor: bs
  name: toggles
  color: grey
---

Below are configuration recommendations for optimizing your Meshtastic nodes for getting on the [Florida Mesh Map][MESHMAP] & [Florida Mesh Telemetry][MALLA].

## Channels

Alongside the primary channel, Florida uses:

| Channel | Key | Purpose |
|---|---|---|
| [`FloridaMesh`](https://meshtastic.org/e/?add=true#CjcSIGXQsmeHp7fzWcnl7zY7Qb666dJQpaHVMoWKVS-MbIZwGgtGbG9yaWRhTWVzaCgBMAE6AggN) | `ZdCyZ4ent/NZyeXvNjtBvrrp0lClodUyhYpVL4xshnA=` | Statewide chat |
| [`emergency`](https://meshtastic.org/e/?add=true#ChYSAQEaCWVtZXJnZW5jeSgBMAE6AggN) | `AQ==` | Emergency and severe weather traffic |
| [`weather`](https://meshtastic.org/e/?add=true#Cg4SAQEaB3dlYXRoZXIoAQ) | `AQ==` | Weather scripts and alerting |
| [`testing`](https://meshtastic.org/e/?add=true#Cg4SAQEaB3Rlc3RpbmcoAQ) | `AQ==` | Node, script, and bot testing |
| [`telemetry`](https://meshtastic.org/e/?add=true#ChISAQEaCXRlbGVtZXRyeSgBOgA) | `AQ==` | Remote infrastructure nodes only, see below |

Open a channel name on a device with the app installed to add it, or enter the name and key by hand as a secondary channel. Name and key together define the channel: same name, same key, same channel. Names are case sensitive and limited to 11 characters. A node has eight channel slots including the primary.

The [Weekly Mesh Net]({{< relref "docs/general/weekly-mesh-net/index.md" >}}) runs on your region's default public channel, not on any of these. Bot-to-bot traffic goes on `testing` or a private channel; keep `emergency` clear. See [Bots and Automation]({{< relref "docs/general/bots-and-automation/index.md" >}}).

{{< notice note "Check for ?add=true" >}}
The links above are `?add=true`: each adds one channel and leaves your other channels and LoRa config untouched.

A `meshtastic.org/e/` link *without* `?add=true` replaces the entire channel set, makes its first channel your primary, and overwrites LoRa config including region, hop limit, and TX power. Verify before scanning a link or QR code from anyone.
{{< /notice >}}

### telemetry

Infrastructure only. Neighbor Info does not transmit over LoRa while the primary channel is the default one, so a remote infrastructure node sharing neighbor data sets `telemetry` as its **primary** rather than a secondary. The add link installs it as a secondary, where it has no effect.

## Getting on the Map In Existing Meshes

For nodes that are in established meshes (please check [Florida Mesh Map][MESHMAP] to see where the closest feeders are) all you need to get added to the maps and tools is one config change.
  
LoRa:

- Ok to MQTT: `Checked`
  
This change will allow for MQTT feeder nodes in the area that can hear you via RF to have permission to go and uplink your node's info to the MQTT Map and Telemetry toolset.
If you would like to help feed the map yourself, or are in an area with limited feeders. Please continue on.

## Setting up a MQTT Gateway

If you would like to connect your nodes to the MQTT broker and provide telemetry to both the [Florida Mesh Map][MESHMAP] & [Florida Mesh Telemetry][MALLA], you will need to configure the following settings:

{{< notice note >}}
The Florida Mesh MQTT Server Primary purpose has been to provide data to help build and grow the Mesh across the state of Florida; hence this server from the start only allowed Uplinking
(meaning comunication via this mqtt would not work). The data has been available to be viewed on both the [Florida Mesh Map](https://map.areyoumeshingwith.us) &
[Florida Mesh Telemetry](https://malla.areyoumeshingwith.us).

Due to demand of users and need to be able to quickly identify and address issues, the Florida Mesh MQTT Server has moved to individual MQTT accounts rather than a singular shared credential.
This now means you will need to use `@MeshBot` on the Discord server to [request your account]({{< relref "docs/general/mqtt-server/index.md" >}}) in order to connect to the Florida Mesh MQTT Server.
{{< /notice >}}

### Radio Configuration

Channels:

- LongFast (primary)
  - Uplink Enabled: `Checked`
  - Downlink Enabled: `Unchecked`[^downlink]
  - Position enabled: `Checked`
  - Precise Location: *user preference*
  - Precision Slider: *user preference*
  
Uplink and downlink are per channel:

| Channel | Uplink | Downlink | Zero-hopped |
|---|---|---|---|
| primary (`LongFast`) | `Checked` | `Unchecked` | Yes |
| `FloridaMesh` | `Checked` | `Checked` | No |
| `emergency` | `Checked` | `Checked` | No |
| `weather`, `testing` | `Checked` | `Unchecked` | No |
| `telemetry` | `Checked` | `Unchecked` | Yes |

`FloridaMesh` and `emergency` take downlink so traffic on them reaches the whole state instead of only nodes in radio range of the sender.

The [Florida broker](https://github.com/flmesh/emqx) zeroes `hop_limit` and `hop_start` in flight on the default modem-preset channels and on `telemetry`. Packets it delivers on those reach only nodes hearing a gateway directly and are not rebroadcast, so MQTT traffic cannot flood the radio mesh. `FloridaMesh` and `emergency` are not zero-hopped; downlinked traffic on them propagates at its hop limit.

Device:

- Role: `CLIENT|CLIENT_BASE|CLIENT_MUTE`[^role]

LoRa:

- Presets: `LONG_FAST`[^presets]
- Hop limit: `3`[^hops]
- Ignore MQTT: `Unchecked`[^ignore-mqtt]
- OK to MQTT: `Checked`[^ok-mqtt]

### Module Configuration

MQTT:

- Enabled: `Checked`
- MQTT Server Address: `mqtt.areyoumeshingwith.us`
- MQTT Username: *your MQTT account credentials*[^hubot]
- MQTT Password: *your MQTT account credentials*[^hubot]
- Encryption Enabled: `Checked`
- JSON Enabled: `Unchecked`
- TLS Enabled: *user preference*[^tls]
- Root topic: `msh/US/FL`

{{< notice tip >}}
If you followed the above MQTT steps in order, then you need to hit "Send" to update the device at this point before you proceed to the below four steps. Otherwise the send button will be grayed out.
{{< /notice >}}
  
- Map reporting: `Checked`
- Map reporting interval (seconds): `10800`
- Precise location: *user preference*
- Precision Slider: *user preference*

Neighbor Info:

- Neighbor Info enabled: `Checked`[^ninfo]
- Update interval (seconds): `14400`
- Transmit over LoRa: `Unchecked`

### Special Settings for Mobile Nodes and NRF52 Nodes

If your wanting to run a MQTT feeder from a mobile NRF52 or ESP32 based node. You'll need to enable an additional setting.

MQTT:

- Proxy to Client: `Enabled`

This will let the node pass MQTT data via the phone to the MQTT server.

If your using a NRF52 based node for a Base Station for example its best practice to use a secondary node that can either run off WIFI or ethernet, like a ESP32 based node or Portduino (Linux) based node, to operate as the feeder node.
If you do run secondary feeder node, remember to set it to `Client_Mute` to prevent it from causing unnecessary noise.
You can feed directly from a NRF52 via a phone but longterm its not the most reliable, since the connection is done over Bluetooth. Its best for adhoc quick deployments and mobile use.

## Verifying Your Configuration

After configuring your device, you can verify that your telemetry is being properly reported:

1. Check the [Florida Mesh Map][MESHMAP] & [Florida Mesh Telemetry][Malla] - your node should appear within 15-30 minutes, check for the hardware MAC ID first, since this will be the first part to populate
2. Review your device debug logs for successful MQTT connection messages
3. Confirm your device is sending position updates at the expected intervals

## Troubleshooting

If your node is not appearing on the map:

- Verify internet connectivity on the device
- If its a NRF52 based node confirm that `Proxy to Client` is enabled.
- Check that your MQTT credentials are entered correctly
- Ensure the precision settings meet the minimum requirements (1194ft / 363m)
- Confirm the root topic is set exactly to `msh/US/FL`
- Verify the MQTT module is enabled and properly configured

## FAQ

### I receive too many notifications on a channel

Notification behaviour is local to your phone. Messages keep arriving either way, and your node keeps relaying traffic for the mesh.

- **Mute the conversation.** Select the channel in the message list, then **Mute notifications**: 8 hours, 1 week, or always.
- **Disable notifications for the app** in your phone's settings.
- **Remove the channel** in the app's channel settings, freeing the slot. Rejoin later with the name and key above.

[^presets]: Please reference [Regional LoRa Settings]({{< relref "regional-lora-settings/index.md" >}}) for up to date modem presets for each area of the state.
[^tls]: TLS encrypts data transmitted between MQTT clients and the broker for increased security, but may not supported on all platforms. It is known that the Android App version above 2.7.13 may have issues with TLS enabled.
[^ninfo]: Please only enable Neighbor Info on Basestation and stationary nodes, when enabled on mobile nodes it causes a lot of noise and clutter to the map. Thank you.
[^role]: CLIENT, CLIENT_BASE, or CLIENT_MUTE unless you have a different use case. **Never set this to ROUTER or REPEATER**.
[^hops]: Except in specific use cases *(user preference: but do not set to higher than 5)*.
[^downlink]: If you wish to receive traffic from MQTT and not just send traffic to MQTT then you may enable with your individual MQTT account credentials.
[^hubot]: Florida Mesh has a new requirement to [request an individual MQTT account]({{< relref "docs/general/mqtt-server/index.md" >}}) which allow both channel `uplink` & `downlink` ability.
[^ignore-mqtt]: Can be Checked to stop rouge MQTT data from appearing on your node and hopping though. Consider unchecking if you uncheck `downlink`.
[^ok-mqtt]: This gives permission to uplink your node to MQTT Servers. You will not appear on the Meshtastic Map if this is not enabled.

[MESHMAP]: https://map.areyoumeshingwith.us "Florida Mesh Map"
[MALLA]: https://malla.areyoumeshingwith.us/ "Florida Mesh Telemetry"
