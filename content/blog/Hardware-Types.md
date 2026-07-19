---
title: "Hardware Types"
date: 2026-07-19T15:47:00-04:00
draft: false
description: Hardware archtypes for use in mesh systems.
noindex: false
featured: true
pinned: true
authors:
    - Json_18
# comments: false
series:
#  -
categories:
#  - 
tags:  
 - ESP32
 - NRF52
 - Linux
#  - 
images: 
# - 
#  - 
# menu:
#   main:
#     weight: 100
#     params:
#       icon:
#         vendor: bs
#         name: book
#         color: '#e24d0e'
---

# Types of Hardware for Meshtastic and Meshcore Systems.

<!--more-->
So, one of the most frequent questions I get asked from new users is "*What hardware should I get for my first node?*"
This used to be a fairly simple question to answer. But as the hobby and market has grown, so has the hardware options and systems that can be used. In this blog post we'll step through the main hardware architypes that are currently available to use.


### ESP32
ESP32 (developed by Espressif Systems) is a versatile microcontroller platform that’s easy to program, flash, and work with thanks to its dual-core processor, built-in Wi-Fi and Bluetooth, and broad community support.

- Has ample amount of GPIO pins (ranging around 30 to 40 pins)

- Typically has 4-16 MB of flash storage and 512 KB of RAM. which allows for larger tables and directories to be kept in memory.

- Power usage is fairly high when used on Meshtastic (cores are always active). On Meshcore as a client the power usage is more manageable. (Meshcore client mode disables the TX of the radio, which helps a lot on power savings)

- Due to the onboard WiFi these boards make for an easy and cheap MQTT feeder node

![M5 Stack C6L](/images/Hardware/M5STACK_C6L.webp)


#### What is the hardware good for?

| ESP32 | Meshtastic | Meshcore |
| :---: | :---: | :---: |
| Client Device | Useable but short battery life | Good battery life in client mode |
| Solar powered node | Useable but needs a larger battery and panel then other hardware types | Useable but needs a larger battery and panel then other hardware types |
| Base station on shore power (powered by mains) | Well suited | Well suited |
| infrastructure Device | Needs WIFI or serial access for firmware [updating]



### NRF52
NRF52 devices (developed by Nordic Semiconductor) are another versatile microcontroller platform that while a little more difficult to program, does have the advantage of far better efficiency then ESP32 devices for battery powered and solar operations. 
- Runs a single core coretex CPU
- Doesn't have WiFi but does carry Bluetooth low energy 5.0
- Memory is far smaller than ESP32, usually averaging around 256 KB to 1 MB of flash and 64 KB of RAM
- Uses a Bootloader and needs to be mounted like a portable storage device to be able to have new firmware added on. This allows firmware to be remotely uploaded to the device via BT. But if the bootloader becomes corrupted due to flakey power supplies or other issues it will require reformatting of the MBR (master boot record). 
- Can run for multiple days at a high duty cycle off of a single 18650 battery cell.

#### What is the hardware good for?

| NRF52 | Meshtastic | Meshcore |
| :---: | :---: | :---: |
| Client Device | Good battery life | Good battery life |
| Solar powered node | Preferred system for most solar nodes | Preferred system for most solar nodes |
| Base station on shore power (powered by mains) | Useable but only easily accessible via BT | Useable but only easily accessible via BT|
| infrastructure Device | Plan for serial access or BT access for firmware updates (Note: can get soft locked if the link with updating device is poor, will require physical access to clear) | 

### Linux / Python
Linux based systems are far more advanced then regular development boards, requiring a basic knowledge of linux and commandline But also doesn't limit you to any preset boards or systems. It also gives the very strong advantage of updating and building in place while not replacing the whole or part of the operating system like ESP32 or NRF52 boards do. This is a crucial advantage on infrastructure sites that are hard to get too or would require a climber to go up and correct the system. 
- Can be built on any operating system that can use python
- Is power efficient when running on low power boards
- Doesn't need firmware updates but updates in place using `sudo apt update`
- Has higher storage and ram limits, removing the need for database entry limits.

#### What is the hardware good for?

| Linux | Meshtastic | Meshcore |
| :---: | :---: | :---: |
| Client Device | Depends on device, good for unorthodox devices | Depends on device, good for unorthodox devices |
| Solar powered node | If built with on a power efficient SoC it will have the power draw of a ESP32 but also the advantage of a linux device |  If built with on a power efficient SoC it will have the power draw of a ESP32 but also the advantages of a linux device |
| Base station on shore power (powered by mains) | Great option for power users | Great option for power users |


[updating]: https://areyoumeshingwith.us/docs/meshtastic/installation/#7-esp32-firmware-upgrades-for-nodes-connected-to-a-pi- "updating"
