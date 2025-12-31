---
title: Meshtasticd Installation Guide
linkTitle: Meshtasticd Installation
description: Guide to install Meshtasticd on a Raspberry Pi
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
