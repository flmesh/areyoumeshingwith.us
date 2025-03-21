---
title: "Florida Mesh Map Go-Live"
date: 2025-03-21T13:09:14-04:00
draft: false
description: The Florida Mesh Map goes live!
noindex: false
featured: true
pinned: false
# comments: false
series:
#  - 
categories:
#  - 
tags:
 - map
images:
#  - 
authors:
  - jbouse
# menu:
#   main:
#     weight: 100
#     params:
#       icon:
#         vendor: bs
#         name: book
#         color: '#e24d0e'
---

# Florida Mesh Map

Within the Florida Mesh community we keep trying to take an RF first approach and keep things as off-grid as best we can. It does help to know where everyone else is so we have begun to work on being able to map the existing nodes out there. This is still also following an opt-in approach as well for those that don't wish their nodes to be shown.

We did not want to re-invent the wheel so we went with the excellent example of [Liam Cottle][MESHMAP]'s Meshtastic Map that utilizes a MQTT feed to populate. To start we forked Liam's [meshtastic-map][GITHUB] GitHub repository, but at this point have not made any modifications to the code. We then setup a private MQTT broker for Florida Mesh for it to point to.

## RF First

With our RF First initiative the concept was to limit the complexity to implement the mapping. To this purpose we intended for not everyone needing to enable MQTT to get nodes on the map, and instead encourage RF meshing with the data being fed in via strategic MQTT GW nodes that had good local mesh coverage along with solid network connectivity.

Taking this approach means that the minimum modification necessary on a node is to simply enable the `OK to MQTT` option under their nodes LoRa configuration. This adds the flag to packets sent out that they are agreeing to be sent over MQTT. That's it! If you enable that option and your node meshes with an MQTT GW node it will be ingested. If you want to be able to participate in communications over channels gatewayed over MQTT or DM nodes over MQTT then you also need to disable the `Ignore MQTT` otherwise your node will simply ignore packets flaged as `via MQTT` as if they didn't exist.

With this deployment model there is no need to enable MQTT on your individual nodes. The MQTT GW nodes currently are enabled for LongFast and some other testing channels. Like any node we are limited by the firmware channel limit for how many can be configured so the current plan
is to only gateway channels that might benefit communications between all local meshes connected. There are some areas that have setup local channels and shared the PSK for them.

## MQTT GW Nodes

The Florida Mesh team has been identifying nodes to be used as MQTT GW nodes and setting up the individual credentials for each. This means there won't be any general public credentials shared, at least at this point, and if a node is observed as causing issues we can easily disable it without affecting the entire network.

We are starting out small and growing as we test and vet the setup, making changes as needed as we go. With that said though, the [Florida Mesh Map][FLMESHMAP] is live and linked on the site menu. We are also looking for ways to work with neighboring mesh communities to bridge data to increase the overall view while still focusing on our local communities.

[MESHMAP]: https://meshtastic.liamcottle.net/ "Liam Cottle's Meshtastic Map"
[GITHUB]: https://github.com/liamcottle/meshtastic-map "Liam Cottle - Meshtastic Map - GitHub"
[FLMESHMAP]: https://map.areyoumeshingwith.us "Florida Mesh Map"
