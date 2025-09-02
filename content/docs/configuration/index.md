---
title: Node Configuration
linkTitle: Configuration
description: Recommendations for configuring your node and contributing to the mesh.
nav_weight: 3
authors:
  - beanfield
series:
  - Guide
date: 2025-04-10T00:31:34-04:00
nav_icon:
  vendor: bs
  name: toggles
  color: white
---

Below are configuration recommendations for optimizing your Meshtastic nodes.

## MQTT Gateway

If you would like to connect your nodes to the MQTT broker and provide telemetry to the [Florida Mesh Map][MESHMAP], you will need to configure the following settings:

### Radio Configuration
  Channels:
  * LongFast (primary)
    * Uplink Enabled: `Checked`
    * Downlink Enabled: `Unchecked` 
    * Position enabled: `Checked`
    * Precise Location: *user preference*
    * Precision Slider: *user preference*

  Device:
  * Role: `Client, unless you have a different use case. *Never set this to ROUTER*`

  LoRa:
  * Ignore MQTT: `Unchecked`
  * OK to MQTT: `Checked`

### Module Configuration
  MQTT:
  * Enabled: `Checked`
  * MQTT Server Address: `mqtt.areyoumeshingwith.us`
  * MQTT Username: `uplink`
  * MQTT Password: `uplink`
  * Encryption Enabled: `Checked`
  * JSON Enabled: `Unchecked`
  * TLS Enabled: *user preference* \*
  * Root topic: `msh/US/FL`
  * Map reporting: `Checked`
  * Map reporting interval (seconds): `10800`
  * Precise location: *user preference*
  * Precision Slider: *user preference*

*\* Note: TLS encrypts data transmitted between MQTT clients and the broker for increased security, but it is not supported on all platforms.*

  Neighbor Info: 
  * Neighbor Info enabled: `Checked`
  * Update interval (seconds): `14400`
  * Transmit over LoRa: `Unchecked`

## Verifying Your Configuration

After configuring your device, you can verify that your telemetry is being properly reported:
1. Check the [Florida Mesh Map][MESHMAP] - your node should appear within 15-30 minutes
2. Review your device logs for successful MQTT connection messages
3. Confirm your device is sending position updates at the expected intervals

## Troubleshooting

If your node is not appearing on the map:
- Verify internet connectivity on the device
- Check that your MQTT credentials are entered correctly
- Ensure the precision settings meet the minimum requirements (1194ft / 363m)
- Confirm the root topic is set exactly to `msh/US/FL`
- Verify the MQTT module is enabled and properly configured

[MESHMAP]: https://map.areyoumeshingwith.us "Florida Mesh Map"
[PRECISION]: https://meshtastic.org/docs/software/integrations/mqtt/#location-precision-filtering "Location Precision Filtering"
