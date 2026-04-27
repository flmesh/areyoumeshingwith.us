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

## Getting on the Map In Existing Meshes

For nodes that are in established meshes (please check [Florida Mesh Map][MESHMAP] to see where the closest feeders are) all you need to get added to the maps and tools is one config change.
  
LoRa:

- Ok to MQTT: `Checked`
  
This change will allow for MQTT feeder nodes in the area that can hear you via RF to have permission to go and uplink your node's info to the MQTT Map and Telemetry toolset.
If you would like to help feed the map yourself, or are in an area with limited feeders. Please continue on.

## Setting up a MQTT Gateway

If you would like to connect your nodes to the MQTT broker and provide telemetry to both the [Florida Mesh Map][MESHMAP] & [Florida Mesh Telemetry][MALLA], you will need to configure the following settings:

{{< notice note >}}
The Florida Mesh MQTT Server Primary purpose is to provide data to help build and grow the Mesh across the state of Florida; hence this server is set to only allow Uplinking (meaning comunication via this mqtt won't work).
{{< /notice >}}

The data is available to be viewed on both the [Florida Mesh Map][MESHMAP] & [Florida Mesh Telemetry][MALLA].

### Radio Configuration

Channels:

- LongFast (primary)
  - Uplink Enabled: `Checked`
  - Downlink Enabled: `Unchecked`
  - Position enabled: `Checked`
  - Precise Location: *user preference*
  - Precision Slider: *user preference*
  
Device:

- Role: `CLIENT|CLIENT_BASE|CLIENT_MUTE`[^role]

LoRa:

- Presets: `LONG_FAST`[^presets]
- Hop limit: `3`[^hops]
- Ignore MQTT: `Unchecked` *(Can be Checked to stop rouge MQTT data from appearing on your node and hopping though)*
- OK to MQTT: `Checked` *(This gives others permission to uplink your node to MQTT Servers)*

### Module Configuration

MQTT:

- Enabled: `Checked`
- MQTT Server Address: `mqtt.areyoumeshingwith.us`
- MQTT Username: `uplink`[^hubot]
- MQTT Password: `uplink`[^hubot]
- Encryption Enabled: `Checked`
- JSON Enabled: `Unchecked`
- TLS Enabled: *user preference*[^tls]
- Root topic: `msh/US/FL`

 *\* Note: If you followed the above MQTT steps in order, then you need to hit "Send" to update the device at this point before you proceed to the below four steps. Otherwise the send button will be grayed out.*
  
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

[^presets]: Please reference [Regional LoRa Mesh Presets][LORA-PRESETS] for up to date modem presets for each area of the state.
[^tls]: TLS encrypts data transmitted between MQTT clients and the broker for increased security, but may not supported on all platforms.
[^ninfo]: Please only enable Neighbor Info on Basestation and stationary nodes, when enabled on mobile nodes it causes a lot of noise and clutter to the map. Thank you.
[^role]: CLIENT, CLIENT_BASE, or CLIENT_MUTE unless you have a different use case. **Never set this to ROUTER or REPEATER**.
[^hops]: Except in specific use cases *(user preference: but do not set to higher than 5)*.
[^hubot]: Florida Mesh has begun testing on a new feature to allow for individual MQTT accounts which allow both channel uplink & downlink ability.

[MESHMAP]: https://map.areyoumeshingwith.us "Florida Mesh Map"
[MALLA]: https://malla.areyoumeshingwith.us/ "Florida Mesh Telemetry"
[LORA-PRESETS]: {{< relref "docs/meshtastic-docs/regional-lora-settings/index.md" >}} "Regional LoRa Mesh Presets"
