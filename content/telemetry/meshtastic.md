---
title: Meshtastic Telemetry
type: docs
menu:
  main:
    identifier: telemetry
    weight: 5

---
 
<!--more-->

## Direct Links to the Meshtastic Telemetry Tools

- #### [Malla](https://malla.areyoumeshingwith.us) 
  [Click](#malla) Here to learn more about Malla 

- #### [Meshview](https://meshview.areyoumeshingwith.us)
  [Click](#meshview) Here to learn more about Meshview

## About the Tools

### Malla

 One of our telemetry systems is powered by [Malla](https://github.com/zenitraM/malla), an open-source web analyzer built specifically for Meshtastic networks. Malla captures packets from our MQTT broker and presents them through an intuitive dashboard with powerful analytics.
Some of the features that you can use with Malla include
- Live Graphs and Statistics up to the minute on the health of the mesh
- Captured Packets and what nodes heard them
- Traceroute graphs (what nodes hear a trace route attempt, and where did it end)
- Node-and-Link Network Graph of the mesh and the signal quality inbetween links
- Line of Sight Analysis
- Longest Observed Links
- And More



### Meshview

 Our second telemetry systems is powered by [Meshview](https://github.com/pablorevilla-meshtastic/meshview). It is also a open-source web analyzer built specifically for Meshtastic networks. Like Malla it also captures packets from our MQTT broker and presents them in an easy to use dashboard, though it does provide some unique features over Malla. 
Here are some of the features Meshview provides.
- Captured Packets and what nodes heard them
- Easy to read dashboard for all Public messages (Messages that are sent via the unencrypted channels.)
- Live Map that splits the data out by the channel name, as well as notes when data is received on the MQTT and from where.
- Daily Graphs that are easily exportable
- Static Node-and-Link Graphs
- Weekly Mesh Net Packet Filter (filters out and displays any messages that has #FLMeshNet)
- And More