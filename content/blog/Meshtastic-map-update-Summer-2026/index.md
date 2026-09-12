---
title: "Meshtastic Mesh Map Update Summer 2026"
date: 2026-09-11T00:00:00-04:00
draft: false
description: To radar, or not to radar - that is no longer the question! We have Weather! Also a round of behind-the-scenes improvements to the Florida Mesh Map including infrastructure node display, and better traceroutes.
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
 - guide
 - telemetry
 - meshtastic
images:
  - Closeupstorm.webp
authors:
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
<div class="meshtastic-map-update">
<style>
/* Page-scoped image border for Meshtastic Mesh Map Update post */
.meshtastic-map-update img,
.meshtastic-map-update figure img,
.meshtastic-map-update .hb-featured-img,
.meshtastic-map-update .hb-gallery img,
.meshtastic-map-update .post-thumbnail,
.meshtastic-map-update .hb-figure img {
  border: 0.5px solid #000 !important;
  box-sizing: border-box;
}

.meshtastic-map-update img,
.meshtastic-map-update figure img {
  display: inline-block;
  max-width: 100%;
  height: auto;
}
</style>

<div align="center"><sup>A evening summer storm marches its way across the state, and the Great Lake of Okeechobee.</sup></div>


Since the [Florida Mesh Map][FLMESHMAP] went live over a year ago, we've kept iterating on it behind the scenes. Most of that work is invisible day-to-day, but it adds up to a faster, more informative map, so here's a rundown on some of the most recent changes.

<!--more-->

## Infrastructure Sites

The backbones of much of our regional meshes, Infrastructure nodes like (Routers, Router_Late, Repeaters, and Router_Clients *the last two deprecated, but still existing in the wild*) matter a great deal for how nodes are able to reach out into the net, and can be the diffrence inbetween being online or a party of one. 

We wanted the map to be able to show these connections and the links inbetween the infra sites, while also keeping inmind that for sustainability sake it had to be dynamically taggable, this sadly does limit us from being able to flag nodes of other roles in the infrastructure group, though that may be something to dig into at a later date.

The map now automatically identifies infrastructure-role nodes and surfaces them as their own overlay `Infrastructure Sites` and filtering the map to only them. So you can see at a glance where the backbone of the network actually is. Especially with the new included feature below.

## Infrastructure Connections
To complement the split out infra sites group, we now have a new `infra connections` 
{{< figure src="infra_link.webp" alt="A freqent link across Tampa Bay" width="75%" class="d-block mx-auto" >}}
These new Green neighbors line only displays when theres a captured trace route or SNR reading inbetween two nodes tagged with the `infra` flag, as layed out above. These lines won't be the most frequent to show up, espcially if both infra nodes aren't reporting neighbors or SNR readings. So the lack of this line does not mean the connection isn't there, see tools like [MALLA] or [Meshview][MESHVIEW] for more reliable and detailed connection information and analysis. These lines are meant for cursory review and inspection. 

## Connection Quality Between Nodes

Popups for the connections between nodes now show recent signal history with a rolling look at the last three SNR (Signal-to-Noise-Ratio) between two nodes rather than just the single snapshot that was there previously. This makes it easier to judge whether a link is consistently solid or just had one good packet come through. {{< figure src="infra_link_popup.webp" alt="A freqent link across Tampa Bay" width="75%" class="d-block mx-auto" >}}
The reading are now color coded as well to help better under stand the links at a glance.

| Color                                 | Meaning                                                                                                                                                                                                                                                                       |
|---------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| {{< mapdot color="#00ff33" >}}  GREEN  | Equal or above 0 dB. This is the normal range that you would be able to hear an analog signal (Very strong LoRa signal)                                                                                                                                                      |
| {{< mapdot color="#FFFF00" >}}  YELLOW | Equal or above -7.5dB but not passed 0. This would be a signal thats just below the noise floor, LoRa can work quite well here, making this only a slightly weaker signal  (Weak to Marginal LoRa signal)                                                                                 |
| {{< mapdot color="#ff0000" >}}  RED    | Below -7.5dB. This is a very weak signal and on most LoRa radios a marginally workable signal at best. Though Radio quality is improving, and now at the time of this writing theres radios that are rated to reach and work down to -150dB (Marginal to Nonexistent signal) |

## Weather is Where its At
{{< figure src="florida_storms.webp" alt="An average Florida summer evening" width="100%" class="d-block mx-auto" >}}
With the peak of Hurricane Season now upon us, we wanted to add a frequently used and useful feature of live radar to the map. In this update we now have live data on a 5 min update cadance from the National Weather Service [NEXRAD] service. Which as a resolution ranging from 0.25 to 1km depending on distance to the radar sites.

We also have the option for worldwide radar imagery as well, provided by [RainViewer] which is at an average of 1km resolution, and is also on a 5 min update cadance. 

## A Faster-Loading Map and ETC.

While we're where at it, We went through and conveted the preexisting system of `.png` for device image files over to `.webp` files. This resulted in a good savings of storage space as well as and more importantly faster rendering times for the map. With the ever growing number of devices every little bit helps.
The Devices are now also updated to include as many of the new devices that have either come out on the market or are about to drop. This is a constant ongoing task, so if you'd like to help out, or have found a better quality image of a device or of one we don't yet have. Please feel free to attach it in a PR at the [Github Flmesh Map Repo] or let me know to add to the que.


[FLMESHMAP]: https://map.areyoumeshingwith.us "Florida Mesh Map"
[MALLA]: https://malla.areyoumeshingwith.us "Malla"
[MESHVIEW]: https://meshview.areyoumeshingwith.us "Meshview"
[NEXRAD]: https://mesonet.agron.iastate.edu/docs/nexrad_mosaic/ "NEXRAD"
[RainViewer]:https://www.rainviewer.com/ "RainViewer"
[Github Flmesh Map Repo]: https://github.com/flmesh/meshtastic-map "FL Mesh Meshtastic Map Github Repo"
</div>