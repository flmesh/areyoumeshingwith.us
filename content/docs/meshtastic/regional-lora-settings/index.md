---
title: "Regional LoRa Settings"
# linkTitle:
date: 2026-01-21T03:31:13Z
draft: false
description: This is an active directory for each of the State's Regions and the Meshtastic LoRa RF settings. This document will be updated and maintained as frequency presets change.
noindex: false
# comments: false
nav_weight: 2
nav_icon:
  vendor: bs
  name: wifi
  color: orange
authors:
  - Json_18
  - jbouse
series:
  - Docs
categories:
#  - 
tags:
#  - 
images:
 - /images/regional-lora-settings/Regional-lora-settings.webP
#  - 
---

## Regions

{{< figure src="regional-lora-settings.webP" alt="Florida Mesh Regions" width="50%" class="float-end me-3" >}}

Here at Florida Mesh we've divided up the State into Five distinct regions. They are as follows:

- 🟣 Purple - North West Florida
- 🔴 Red - North Florida
- 🔵 Blue - Central West Florida
- 🟡 Yellow - Central East Florida
- 🟠 Orange - South Florida

Each Region's LoRa settings and Metro based Meshes settings will be called out in each Region Section.

## 🟣 North West Florida

{{< notice info >}}
There is no standardized alternative modes here
{{< /notice >}}

- LoRa: `LongFast`
- Frequency Slot: `0` (default masking)[^longfast]

---

## 🔴 North Florida

{{< notice info >}}
There is no standardized alternative modes here.  
{{< /notice >}}

- LoRa: `LongFast`
- Frequency Slot: `0` (default masking)[^longfast]

---

## 🔵 Central West Florida

{{< notice info >}}
There is no standardized alternative modes here.  
{{< /notice >}}

- LoRa: `LongFast`
- Frequency Slot: `0` (default masking)[^longfast]

---

## 🟡 Central East Florida

{{< notice >}}
As of 2026/21, Sebastian to Layton key is **MEDIUM_FAST** If your in range of the Eastern Seaboard in this region, please see the [South Florida]({{< relref "regional-lora-settings/index.md#sebastian-to-layton-key" >}}) Section below for details
{{< /notice >}}

- LoRa: `LongFast`
- Frequency Slot: `0` (default masking)[^longfast]

---

## 🟠 South Florida

### Sebastian to Layton key

{{< notice >}}
As of 2026/21, Sebastian to Layton key is **MEDIUM_FAST** If your in range of the Eastern Seaboard in this region, Medium_Fast is the frequency you'll want to be on.
{{< /notice >}}

- LoRa: `MediumFast`
- Frequency Slot: `0` (default masking)[^mediumfast]

---

### All other South Florida Areas

- LoRa: `LongFast`
- Frequency Slot: `0` (default masking)[^longfast]

[^longfast]: If your using a non-default channel name, in channel slot `0`, you will need to place in the correct frequency slot manually, which is slot `20`.
[^mediumfast]: If your using a non-default channel name, in channel slot `0`, you will need to place in the correct frequency slot manually, which is slot `45`.
