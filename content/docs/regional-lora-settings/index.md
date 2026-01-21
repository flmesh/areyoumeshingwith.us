---
title: "Regional LoRa Settings"
# linkTitle:
date: 2026-01-21T03:31:13Z
draft: false
description: This is an active directory for each of the State's Regions and the Meshtastic LoRa RF settings. This document will be updated and maintained as frequency presets change.
noindex: false
# comments: false
nav_weight: 2
# nav_icon:
#   vendor: bootstrap
#   name: toggles
#   color: '#e24d0e'
authors:
  - Json_18
series:
  - Docs
categories:
#  - 
tags:
#  - 
images:
#  - 
# menu:
#   main:
#     weight: 100
#     params:
#       icon:
#         vendor: bs
#         name: book
#         color: '#e24d0e'
nav_icon:
  vendor: bs
  name: wifi
  color: orange
---

## Regions


{{< figure src="regional-lora-settings.webP" alt="Florida Mesh Regions" width="50%" class="float-end me-3" >}}
#### Here at Florida Mesh we've divided up the State into Five distinct regions. They are as follows:
* <span style="display: inline-block; width: 16px; height: 16px; background-color: #aa66cd; margin-right: 8px;"></span>Purple – North West Florida
* <span style="display: inline-block; width: 16px; height: 16px; background-color: #a80000; margin-right: 8px;"></span>Red – North Florida
* <span style="display: inline-block; width: 16px; height: 16px; background-color: #0084a8; margin-right: 8px;"></span>Blue – Central West Florida
* <span style="display: inline-block; width: 16px; height: 16px; background-color: #d1c800; margin-right: 8px; border: 1px solid #999;"></span>Yellow – Central East Florida
* <span style="display: inline-block; width: 16px; height: 16px; background-color: #e69800; margin-right: 8px;"></span>Orange – South Florida

### Each Region's LoRa settings and Metro based Meshes settings will be called out in each Region Section.

## <span style="display: inline-block; width: 23px; height: 23px; background-color: #aa66cd; margin-right: 8px;"></span>North West Florida

>##### Special Notes: There is no standardized alternative modes here.  

>LoRa - `Long Fast`

>Frequency Slot - `0` (default masking) * if your using a nonstandard Channel, in Channel slot - `0`, you will need to place in the correct frequency slot manually, which is slot `20`.

---
---

## <span style="display: inline-block; width: 23px; height: 23px; background-color: #a80000; margin-right: 8px;"></span>North Florida

>##### Special Notes: There is no standardized alternative modes here.  

>LoRa - `LongFast`

>Frequency Slot - `0` (default masking) * if your using a nonstandard Channel, in Channel slot - `0`, you will need to place in the correct frequency slot manually, which is slot `20`.

---
---

## <span style="display: inline-block; width: 23px; height: 23px; background-color: #0084a8; margin-right: 8px;"></span>Central West Florida

>##### Special Notes: There is no standardized alternative modes here.  

>LoRa - `LongFast`

>Frequency Slot - `0` (default masking) * if your using a nonstandard Channel, in Channel slot - `0`, you will need to place in the correct frequency slot manually, which is slot `20`.

---
---

## <span style="display: inline-block; width: 23px; height: 23px; background-color: #d1c800; margin-right: 8px; border: 1px solid #999;"></span>Central East Florida

>##### Special Notes: There is no standardized alternative modes here.  

>LoRa - `LongFast`

>Frequency Slot - `0` (default masking) * if your using a nonstandard Channel, in Channel slot - `0`, you will need to place in the correct frequency slot manually, which is slot `20`.

---
---

## <span style="display: inline-block; width: 23px; height: 23px; background-color: #e69800; margin-right: 8px;"></span>South Florida

>##### Special Notes: As of 2025/09, Port St. Lucie to Islamorada is **MEDIUM_FAST** If your in range of the Eastern Seaboard in this region, Medium_Fast is the frequency you'll want to be on. 

  >LoRa - `MediumFast`

  >Frequency Slot - `0` (default masking) * if your using a nonstandard Channel, in Channel slot - `0`, you will need to place in the correct frequency slot manually, which is slot `45`.

---

>##### All other South Florida Areas:

>LoRa - `LongFast`

>Frequency Slot - `0` (default masking) * if your using a nonstandard Channel, in Channel slot - `0`, you will need to place in the correct frequency slot manually, which is slot `20`.

---
---