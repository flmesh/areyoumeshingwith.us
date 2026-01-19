---
title: "Florida Mesh Telemetry Is now Live"
date: 2026-01-19T20:16:08Z
draft: false
description: Florida Mesh Malla Telemetry 
noindex: false
featured: true
pinned: false
# comments: false
series:
#  - 
categories:
#  - 
tags:
  - telemetry
images:
#  - 
authors:
  - beanfield
  - Json_18
# menu:
#   main:
#     weight: 100
#     params:
#       icon:
#         vendor: bs
#         name: book
#         color: '#e24d0e'
---

## Florida Mesh Telemetry Is Now Officially Live!

We're excited to announce that the Florida Mesh Telemetry dashboard is now live! This new tool gives you real-time insights into our mesh network's health, activity, and coverage across Florida.

### What is Malla?

Our telemetry system is powered by [Malla](https://github.com/zenitraM/malla), an open-source web analyzer built specifically for Meshtastic networks. Malla captures packets from our MQTT broker and presents them through an intuitive dashboard with powerful analytics.

### What Can You See?

Visit [Mesh Telemetry](https://malla.areyoumeshingwith.us/) to explore:

- **Live network statistics** - Track total nodes, active nodes, gateway diversity, and message volume
- **Interactive maps** - Visualize node locations and RF links across Florida
- **Network graphs** - See how packets route through the mesh with multi-hop visualizations
- **Signal analysis** - Review RSSI/SNR distributions and identify longest links
- **Activity trends** - Monitor 7-day network activity and protocol usage
- **Packet browser** - Search and filter individual messages with export capabilities

Whether you're troubleshooting your node, planning new installations, or just curious about mesh activity in your area, the telemetry dashboard provides the data you need.

Check it out at **[malla.areyoumeshingwith.us](https://malla.areyoumeshingwith.us/)** and let us know what you think!

### Want to Contribute?

Help us expand telemetry coverage across Florida! To see your node appear in Malla and contribute to the mesh network data:

**For existing mesh participants:**
- Simply enable **"OK to MQTT"** in your LoRa settings
- This allows nearby MQTT gateway nodes to upload your telemetry to the dashboard
- Your node will appear on both the [Florida Mesh Map](https://map.areyoumeshingwith.us) and Malla within 15-30 minutes

**Expand coverage by becoming a gateway:**
- Set up your node as an MQTT gateway to feed telemetry directly to our broker
- Gateway nodes help relay data from other nodes in your area, significantly expanding network visibility
- This is especially valuable in areas with limited existing coverage

For detailed configuration instructions, check out our [recommended node settings](https://areyoumeshingwith.us/docs/configuration/) guide.