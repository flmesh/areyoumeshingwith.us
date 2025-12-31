---
title: Meshtasticd Installation Guide
linkTitle: Meshtasticd Installation, ESP32 remote Firmware Updating
description: Guide to install Meshtasticd on a Raspberry Pi and ESP32 Firmware Updating
nav_weight: 2
authors:
  - Json_18
series:
  - Guide
  date: 2025-12-30T00:35-04:00
nav_icon:
  vendor: bs
  name: cloud-download
  color: green
---


###### Thank You to Cmgamble71 on the FLorida Mesh Discord Server for putting together this guide
Here is a step-by-step guide to get the Meshtastic CLI running on your Raspberry Pi.
Version used, P2 Zero 2W, Raspian OS Lite-64 bit

## 1. Update Your System ⚙️

First, it's always a good idea to make sure your Raspberry Pi's package list and installed software are up to date.
Open a terminal and run the following command:

`Bash`


`sudo apt update && sudo apt full-upgrade -y`

If desired, adjust hostname and expand file system using sudo raspi-config



## 2. Install Python and Dependencies

Raspberry Pi OS usually comes with Python 3, but we'll run this command to ensure you have it, along with pip (the Python package installer) and venv (the virtual environment tool).

`Bash`


`sudo apt install python3 python3-pip python3-venv -y`



## 3. Create a Project Directory and Virtual Environment

A virtual environment is an isolated space for Python projects, so the packages you install for Meshtastic won't interfere with your system's global Python setup.

  1. **Create a folder** for your Meshtastic project and move into it.
     `Bash`
     `mkdir ~/meshtastic-cli`
     `cd ~/meshtastic-cli`


  2. **Create the virtual environment** inside this folder. We'll name the environment venv.
     `Bash`
     `python3 -m venv venv`

You will now see a new folder named venv inside your meshtastic-cli directory.

## 4. Activate the Virtual Environment

You must "activate" the environment before you can use it. This is a crucial step that is often missed.
Run the following command. Note the . or source at the beginning.

`Bash`


`source venv/bin/activate`


Your command prompt should now change to show (venv) at the beginning. This indicates that the virtual environment is **active**.
(venv) pi@raspberrypi:~/meshtastic-cli $

## 5. Install the Meshtastic CLI

Now that your virtual environment is **active**, you can safely install the Meshtastic CLI using pip.

`Bash`


`pip install -U meshtastic`


This will download and install the meshtastic package and all its required dependencies **only** inside your venv, keeping your system clean.  The -U option will also upgrade the packages if they’re already installed and newer versions are available.

## 6. Verify the Installation ✅

To make sure everything worked correctly, check the version of the installed tool.

`Bash`


`meshtastic --version`


If the installation was successful, it will print the version number, something like meshtastic version 2.x.x.
You're all set! You can now connect your Meshtastic device via USB and run commands like meshtastic --nodes or meshtastic --info.

## 7. ESP32 Firmware Upgrades For Nodes Connected to a PI ✅
You can install or upgrade meshtastic on esp32 devices using the provided scripts in each firmware zip (device-install.sh or device-update.sh).  At the time of this writing, these scripts do not correctly implement the trigger to put them in Device Firmware Update (DF) mode.  You can either put them in DFU manually (different for each device type) or manually invoke esptool with a baud rate of 1200 to trigger DFU mode from your raspberry pi.

Firmware update example, tested on StationG2:
`Bash`

`python -m esptool --baud 1200 write-flash 0x10000 firmware-station-g2-2.6.11.60ec05e-update.bin`

Once the upgrade is complete, your station g2 should automatically reboot with the updated firmware. 

``` shell
(meshtastic-venv) ubuntu@rpi-g2admin:/tmp/firmware-esp32s3-2.7.7.5ae4ff9 $ python -m esptool --baud 1200 write-flash 0x10000  firmware-station-g2-2.7.7.5ae4ff9-update.bin 
esptool v5.1.0
Connected to ESP32-S3 on /dev/ttyACM0:
Chip type:          ESP32-S3 (QFN56) (revision v0.2)
Features:           Wi-Fi, BT 5 (LE), Dual Core + LP Core, 240MHz, Embedded PSRAM 8MB (AP_3v3)
Crystal frequency:  40MHz
USB mode:           USB-Serial/JTAG
MAC:                30:ed:a0:38:8c:24
Stub flasher running.
Configuring flash size...
Flash will be erased from 0x00010000 to 0x001f9fff...
Wrote 2006976 bytes (1275511 compressed) at 0x00010000 in 18.7 seconds (859.7 kbit/s).
Hash of data verified.
Hard resetting via RTS pin...
```




Note: see the Downloading section in [Notes](#important-final-notes) for details on downloading the correct firmware file. 

## Important Final Notes

Deactivating: When you are finished using the CLI, you can leave the virtual environment by simply typing:
`Bash`
`deactivate`


Reactivating: The next time you want to use the Meshtastic CLI, you must navigate back to your project folder and reactivate the environment:
`Bash`
`cd ~/meshtastic-cli`
`source venv/bin/activate`

Downloading: The latest stable (beta) or alpha firmware is in the `Firmware` section at the bottom of the [downloads](https://meshtastic.org/downloads/) page.  Once you select the version, you will be redirected to the github page for that release.  Scroll down to the bottom `Assets` section, right click on the `firmware-esp32s3-{version}.zip` link, and select `copy link address` from your browser’s drop down to get the full link in your clipboard. Then go to your ssh session on the raspberry pi and download it to an appropriate location of your choosing, unzip it, and continue with the install. For example, to download 2.6.11 into /tmp (non-persistent storage):

``` shell
cd /tmp
wget https://github.com/meshtastic/firmware/releases/download/v2.6.11.60ec05e/firmware-esp32s3-2.6.11.60ec05e.zip  
unzip -d firmware-esp32s3-2.6.11.60ec05e firmware-esp32s3-2.6.11.60ec05e.zip
cd firmware-esp32s3-2.6.11.60ec05e
	```


Note: the example above downloads firmware to /tmp which is not persisted if the node is rebooted. If you would like to keep firmware long term for an emergency offline recovery / rollback, choose another location.  
