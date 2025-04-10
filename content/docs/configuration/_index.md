---
title: Configuration
date: 2025-04-10T00:31:34-04:00
nav_icon:
  vendor: bs
  name: gear
  color: white
---

Below are configuration recommendations for optimizing your Meshtastic nodes.

## MQTT Gateway

If you would like to connect your nodes to the MQTT broker and provide telemetry to the [Florida Mesh Map](https://map.areyoumeshingwith.us), you will need to configure the following settings:

```
  Radio Configuration
    Channels:
    * Primary (usually LongFast)
      * Uplink Enabled: `Checked`
      * Downlink Enabled: `Unchecked` 
      * Position enabled: `Checked`
      * Precise Location: `Unchecked`
      * Precision Slider: `1194ft / 363m or greater`

    LoRa:
    * Ignore MQTT: `Unchecked`
    * OK to MQTT: `Checked`

  Module Configuration
    MQTT:
    * Enabled: `Checked`
    * MQTT Server Address: `mqtt.meshtastic.org`
    * MQTT Username: `meshdev`
    * MQTT Password: `large4cats`
    * Encryption Enabled: `Checked`
    * JSON Enabled: `Unchecked`
    * TLS Enabled: `Unchecked`
    * Root topic: `msh/US/FL`
    * Map reporting: `Checked`
    * Precise location: `Unchecked`
    * Precision Slider: `1194ft / 363m or greater`

    Neighbor Info:
    * Neighbor Info enabled: `Checked`
    * Transmit over LoRa: `Unchecked`
```

*\* Note: If you leave precise location enabled or set precision lower than [1194ft / 363m](https://meshtastic.org/docs/software/integrations/mqtt/#location-precision-filtering), the MQTT server will silently drop the packets.*

## Verifying Your Configuration

After configuring your device, you can verify that your telemetry is being properly reported:

1. Check the [Florida Mesh Map](https://map.areyoumeshingwith.us) - your node should appear within 15-30 minutes
2. Review your device logs for successful MQTT connection messages
3. Confirm your device is sending position updates at the expected intervals

## Troubleshooting

If your node is not appearing on the map:

- Verify internet connectivity on the device
- Check that your MQTT credentials are entered correctly
- Ensure the precision settings meet the minimum requirements (1194ft / 363m)
- Confirm the root topic is set exactly to `msh/US/FL`
- Verify the MQTT module is enabled and properly configured
