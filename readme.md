<!-- README.md -->
<p align="center">
  <img src="./logo.png" alt="Project Logo" width="200" height="200"/>
</p>

<h1 align="center">GoWatchIt</h1>

<p align="center">
  <a href="https://github.com/iloveicedgreentea/GoWatchIt/releases">
    <img src="https://img.shields.io/github/v/release/iloveicedgreentea/gowatchit
" alt="Version" />
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/License-CC_Custom-blue" alt="License MIT" />
  </a>
  <a href="https://github.com/iloveicedgreentea/plex-webhook-automation/actions">
    <img src="https://github.com/iloveicedgreentea/plex-webhook-automation/workflows/Docker/badge.svg" alt="CI Status" />
  </a>
  <a href="https://www.avsforum.com/threads/gowatchit-beq-ezbeq-plex-webhook-automation-tool-official-thread.3264800/">
    <img src="https://img.shields.io/website-up-down-green-red/http/shields.io.svg" alt="Website" />
  </a>
</p>

<p align="center">
  <b>Focus on watching your content, not babysitting it.</b><br>
  Automate Your Theater
</p>

---

## Table of Contents
- [Features](#features)
- [Setup](#setup)
- [Usage](#usage)
- [Home Assistant Quickstart](#Home-Assistant-Quickstart)

## Features

This application is primarily focused on Plex and HomeAssistant but I plan on adding support for other sources in the future. 

Main features:
* Load/unload BEQ profiles automatically, without user action and the correct codec detected
* Set volume based on media type (Movie, Show, etc)
* Trigger lights when playing, pausing, or stopping automatically (e.g turn off lights on play, turn on when paused)
* HDMI Sync detection and automation (pause while HDMI is syncing so you don't sit embarrassed with a audio playing to a black screen)
* Web based UI for configuration

Other cool stuff:
* Mute/Unmute Minidsp automation for things like turning off subs at night
* Various MQTT sensors for volume control, lights, mute status, and current BEQ profile
* Mobile notifications (via HA) to notify for events like loading/unloading BEQ was successful or failed
* Dry run and notification modes to verify BEQ profiles without actually loading them
* Built in support for Home Assistant and Minidsp


## Setup
> ⚠️ ⚠️ *Warning: You should really set a compressor on your minidsp for safety as outlined in the [BEQ forum post](https://www.avsforum.com/threads/bass-eq-for-filtered-movies.2995212/). I am not responsible for any damage* ⚠️ ⚠️
### Prerequisites
> ℹ  It is assumed you have ezBEQ, Plex, and HomeAssistant working. Refer to their respective guides for installation help.
* MQTT Broker
* Home Assistant
* Plex
* ezBEQ
* Minidsp (other DSPs may work but I have not tested them)

### Docker Setup
> ℹ  If you need help deploying with Docker, refer to the [Docker documentation](https://docs.docker.com/get-docker/).


1) Deploy the latest version `ghcr.io/iloveicedgreentea/plex-webhook-automation:latest`. I recommend running this in an orchestrator like Unraid, Docker-Compose, etc
2) Configure the application via web ui -> `http://(you-server-ip):9999`
   * can get your player UUID from `https://plex.tv/devices.xml` while logged in
3) Set up Plex to send webhooks to your server IP, `listenPort`, and the handler endpoint of `/plexwebhook`
    * e.g `(your-server-ip):9999/plexwebhook`
3) Whitelist your server IP in Plex so it can call the API without authentication. [Docs](https://support.plex.tv/articles/200890058-authentication-for-local-network-access/)
4) Play a movie and check server logs. It should say what it loaded and you should see whatever options you enabled work
5) The Application will restart within 5 seconds when config is saved in the UI

### Non-Docker Setup
I don't recommend this as it is more work and you will need to set up systemd or something to keep it running. I don't provide support for this method but if you know what you are doing, it is very easy to build the binary and run it.

TLDR: `make build`


## Usage

### Web UI
The web UI is the primary way to configure this application. It is available at `http://(your-server-ip):9999`

It will automatically restart the application when you save config.

### General Usage
This application will load BEQ profiles automatically when you play something in Plex. It will also set volume, lights, and mute/unmute minidsp if you enable those options. The application itself is not controlling things like lights but relies on Home Assistant to perform the action via MQTT.

## Home Assistant Quickstart

### MQTT
MQTT is used so this application could theoretically be used with any home automation system. Only Home Assistant is officially supported. You will need to set MQTT up first. Detailed instructions here https://www.home-assistant.io/integrations/mqtt/
  
1) Install mosquito mqtt add on
2) Install mqtt integration
3) Set up your topics in HA and the application's config
4) Set up Automations in HA based on the payloads of MQTT

Features that will write to Topics
Current BEQ Profile
Lights
Minidsp mute status
Volume
Playing status

These Topics allow you to trigger automations in HA based on sensor values such as triggering HVAC when playing status is true.

Here are some sensor examples

```yaml
mqtt:
  binary_sensor:
    - name: "subs_muted"
      state_topic: "theater/subs/status"
      payload_on: "false"
      payload_off: "true"
    - name: "plex_playing"
      state_topic: "theater/plex/playing"
      payload_on: "true"
      payload_off: "false"
  sensor:
    - name: "lights"
      state_topic: "theater/lights/front"
      value_template: "{{ value_json.state }}"
    - name: "volume"
      state_topic: "theater/denon/volume"
      value_template: "{{ value_json.type }}"
    - name: "beq_current_profile"
      state_topic: "theater/beq/currentprofile"


```

### Automation Example
Here is an example of an automation to change lights based on MQTT.

Assuming you have the following sensor:
```yaml
mqtt:
  sensor:
    - name: "lights"
      state_topic: "theater/lights/front"
      value_template: "{{ value_json.state }}"
```

This will turn the light(s) on/off depending on the state of the sensor, state is changed by a message sent to the topic

```yaml
alias: MQTT - Theater Lights
description: Trigger lights when mqtt received, depending on state
trigger:
  - platform: mqtt
    topic: theater/lights/front
condition: []
action:
  - if:
      - condition: state
        entity_id: sensor.lights
        state: "on"
    then:
      - service: light.turn_on
        data: {}
        target:
          entity_id: light.caseta_r_wireless_in_wall_dimmer
  - if:
      - condition: state
        entity_id: sensor.lights
        state: "off"
    then:
      - service: light.turn_off
        data: {}
        target:
          entity_id: light.caseta_r_wireless_in_wall_dimmer
mode: queued
max: 10
```


### Handlers
`/plexwebhook`
This endpoint is where you should tell Plex to send webhooks to. It automatically processes them. No further action is needed. This handler does most of the work - Loading BEQ,  lights, volume, etc

`/minidspwebhook`
This endpoint accepts commands used by minidsp-rs which are performed by EZbeq. Here is how to trigger it with Home Assistant

```yaml
rest_command:
  minidsp:
    url: "http://192.168.88.56:9999/minidspwebhook"
    method: POST
    payload: '{"command": "{{ command }}" }'
    content_type:  'application/json'
    verify_ssl: false
```

And then inside an automation, you make an action
```yaml
  # unmute subs
  - service: rest_command.minidsp
    data:
      command: "off" (or "on")
```

Using the above you can automate the mute and unmute of your minidsp with any automation source.
 
You can then do cool stuff like create a binary sensor to show the state of the subs based on the MQTT topic

### Config
The only supported way to configure this is via the web UI. You can also edit the config.json file directly but this is not supported and will probably break things.

Refer to the web UI for the latest config options.

### Authentication
You must whitelist your server IP in "List of IP addresses and networks that are allowed without auth"

Why? Plex refuses to implement client to server authentication and you must go through their auth servers. I may eventually implement their auth flow but it is not a priority.

### Debugging
These are environment variables you can set to get more info

`LOG_LEVEL=debug` to have it print out debug logs while running

`SUPER_DEBUG=true` for each line to also have a trace to its call site and line number

## How BEQ Support Works
On play and resume, it will load the profile. On pause and stop, it will unload it (so you don't forget to). It has some logic to cache the profile so if you pause and unpause, the profile will get loaded much faster as it skips searching the DB and stuff. 

If enabled, it will also send a notification to Home Assistant via Notify so you can send an alert to your phone for example. 

For safety, the application tries to unload the profile when it loads up each time in case it crashed or was killed previously, and will unload before playing anything so it doesn't start playing something with the wrong profile. 

### Matching
The application will search the catalog and match based on codec (Atmos, DTS-X, etc), title, year, and edition. I have tested with multiple titles and everything matched as expected.

> ⚠️ *If you get an incorrect match, please open a github issue with the full log output and expected codec and title*

### Editions

This application will do its best to match editions. It will look for one of the following:
1) Plex edition metadata. Set this from your server in the Plex UI
2) Looking at the file name if it contains `Unrated, Ultimate, Theatrical, Extended, Director, Criterion`

There is no other reliable way to get the edition. If an edition is not matched, BEQ will fail to load for safety reasons (different editions have different mastering, etc). If a BEQCatalog entry has a blank edition, then edition will not matter and it will match based on the usual criteria.

If you find repeated match failures because of editions, open a github issue with debug logs of you triggering `media.play`

### HDMI Sync Automation
This application supports automatically waiting until HDMI sync is complete. 

Have you ever started something in Plex only to hear audio but see a black screen for 10 seconds? Then everyone in your theater makes fun of you and you cry yourself to sleep? This application will prevent that. 

It supports four ways to get this info: my JVC integration, my Envy integration, a generic binary_sensor, or you can pass in seconds to wait.

If using the first two methods, you just need to install the integration, then set `remoteentityname` to the name of the remote entity. Set `signal.source` to either `jvc` or `envy`.

If using a binary_sensor, you need to create an automation which will set the state to `off` when there is NO SIGNAL and `on` when there is. Getting that data is up to you. Set the `signal` config to the name of the binary sensor (e.g signal, if the entity is binary_sensor.signal).

If using seconds, provide the number of seconds to wait as a string such as "13" for 13 seconds. You can time how long your sync time is and add it here. It will pause for that amount of time then continue playing.

You also must set `plex.playerMachineIdentifier` and `plex.playerIP`. To get this:
1) Play something on your desired player (like a shield)
2) `curl "http://(player IP):32500/resources"`
    * Note this is *NOT THE SERVER IP!* and *only works while something is actively playing*
3) Copy the `machineIdentifier` value
4) Add this to that config field exactly as presented

### Audio stuff
Here are some examples of what kind of codec tags Plex will spit out based on file metadata

 TrueHD 7.1
Unknown (TRUEHD 7.1) --- Surround 7.1 (TRUEHD) 
English (TRUEHD 7.1) --- TrueHD 7.1 (English) 
English (TRUEHD 7.1) --- Surround 7.1 (English TRUEHD) 
English (TRUEHD 7.1) --- English (TRUEHD 7.1) 

 Atmos
English (TRUEHD 7.1) --- TrueHD Atmos 7.1 Remixed (English) 
English (TRUEHD 7.1) --- TrueHD Atmos 7.1 (English) 
English (TRUEHD 7.1) --- TrueHD 7.1. Atmos (English 7.1) 
English (TRUEHD 7.1) --- TrueHD 7.1 Atmos (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1+13 objects / 48 kHz / 4691 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1+11 objects / 48 kHz / 4309 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1 / 48 kHz / 5026 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1 / 48 kHz / 4396 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1 / 48 kHz / 4353 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1 / 48 kHz / 3258 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1 / 48 kHz / 3041 kbps / 16-bit (English) 
English (TRUEHD 7.1) --- Dolby Atmos/TrueHD Audio / 7.1-Atmos / 48 kHz / 5826 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby Atmos/TrueHD Audio / 7.1-Atmos / 48 kHz / 4535 kbps / 24-bit (English) 

 TrueHD 5.1
English (TRUEHD 5.1) --- English Dolby TrueHD (5.1) 
English (TRUEHD 5.1) --- English (TRUEHD 5.1) 
English (TRUEHD 5.1) --- Dolby TrueHD Audio / 5.1 / 48 kHz / 4130 kbps / 24-bit (English) 
English (TRUEHD 5.1) --- Dolby TrueHD Audio / 5.1 / 48 kHz / 1522 kbps / 16-bit (English) 
English (TRUEHD 5.1) --- TrueHD 5.1 (English) 
English (TRUEHD 5.1) --- Dolby TrueHD 5.1 (English) 

 DTS:X 
English (DTS-HD MA 7.1) --- DTS:X/DTS-HD Master Audio / 7.1-X / 48 kHz / 4458 kbps / 24-bit (English) 
English (DTS-HD MA 7.1) --- DTS:X/DTS-HD Master Audio / 7.1-X / 48 kHz / 4255 kbps / 24-bit (English) 
English (DTS-HD MA 7.1) --- DTS:X 7.1 (English DTS-HD MA) 

 DTS-HD MA 7.1
English (DTS-HD MA 7.1) --- Surround 7.1 (English DTS-HD MA) 
English (DTS-HD MA 7.1) --- DTS-HD MA 7.1 (English) 

 DTS-MA 6.1
中文 (DTS-HD MA 6.1) --- DTS-HD Master Audio / 6.1 / 48 kHz / 4667 kbps / 24-bit (中文) 

 DTS-HD MA 5.1
中文 (DTS-HD MA 5.1) --- Mandarin DTS-HD MA 5.1 (中文) 
中文 (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 2360 kbps / 16-bit (中文) 
Indonesia (DTS-HD MA 5.1) --- DTS-HD Master Audio / Indonesian / 5.1 / 48 kHz / 3531 kbps / 24-bit (Indonesia) 
Indonesia (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 3576 kbps / 24-bit (Indonesia) 
English (DTS-HD MA 5.1) --- Surround 5.1 (English DTS-HD MA) 
English (DTS-HD MA 5.1) --- English / DTS-HD Master Audio / 5.1 / 48 kHz / 4104 kbps / 24-bit (DTS Core: 5.1 / 48 kHz / 1509 kbps / 24-bit) 
English (DTS-HD MA 5.1) --- English / DTS-HD Master Audio / 5.1 / 48 kHz / 2688 kbps / 24-bit 
English (DTS-HD MA 5.1) --- English (DTS-HD MA 5.1) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio English 2877 kbps 5.1 / 48 kHz / 2877 kbps / 24-bit 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / English / 3336 kbps / 5.1 Channels / 48 kHz / 24-bit (DTS Core: 5.1 Channels / 48 kHz / 1509 kbps / 24-bit) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48kHz / 5128 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 4189 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 4107 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 3900 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 3746 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 3600 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 3596 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 3233 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD MA 5.1 (English) 
English (DTS-HD MA 5.1) --- DTS-HD MA @ 1509 kbps (English 5.1) 
English (DTS-HD MA 5.1) --- DTS HD MA 5.1 (English) 
English (DTS-HD MA 5.1) --- DTS (MA) / 2181 kbps / 48 KHz / 24-Bit / 5.1 (English) 

 DTS HD HRA 7.1
English (DTS-HD HRA 7.1) --- Surround 7.1 (English DTS-HD HRA) 

 DTS 5.1
English (DTS 5.1) --- English (DTS 5.1) 
English (DTS 5.1) --- DTS 5.1 (English) 

Unknown (AAC Stereo) --- Unknown (AAC Stereo) 

 AC3 stereo
English (AC3 Stereo) --- English (AC3 Stereo) 

 AC3
English (AC3 5.1) --- Surround (English AC3 5.1) 
English (AC3 5.1) --- English (AC3 5.1) 
English (AC3 5.1) --- AC3 5.1 @ 640 Kbps (English) 

 Misc
English (PCM Mono) --- Mono (English PCM) 
English (FLAC Stereo) --- Original Dolby Stereo (Laserdisc USA LD68993) (English FLAC) 
English (FLAC Stereo) --- FLAC Audio / 1266 kbps / 2.0 / 48 kHz / 24-bit (English) 

English (FLAC 5.1) --- Main Audio (English FLAC 5.1) 
English (FLAC 5.1) --- FLAC 5.1 @ 2954 kbps / 24-bit (English) 
English (FLAC 5.1) --- English (FLAC 5.1) 

English (EAC3 Stereo) --- English (EAC3 Stereo) 
English (EAC3 5.1) --- Main Audio (English EAC3 5.1) 
English (EAC3 5.1) --- English (EAC3 5.1) 
