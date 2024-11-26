<!-- README.md -->
<p align="center">
  <img src="./logo-upscale.png" alt="Project Logo" width="200" height="200"/>
</p>

<h1 align="center">GoWatchIt</h1>

<p align="center">
  <a href="https://github.com/iloveicedgreentea/GoWatchIt/releases">
    <img src="https://img.shields.io/github/v/release/iloveicedgreentea/gowatchit" alt="Version" />
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/License-CC_Custom-blue" alt="License MIT" />
  </a>
  <a href="https://github.com/iloveicedgreentea/gowatchit/actions">
    <img src="https://github.com/iloveicedgreentea/gowatchit/workflows/Docker/badge.svg" alt="CI Status" />
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
- [How BEQ Support Works](#How-BEQ-Support-Works)
- [Help](#help)

## Features

Players Supported:
* Plex 
* Jellyfin (no support given, but tested)
* Emby (may work due to jellyfin support, no support given and not tested)

Main features:
* Load/unload BEQ profiles automatically, without user action and the correct codec detected
* HDMI Sync detection and automation (pause while HDMI is syncing so you don't sit embarrassed with a audio playing to a black screen)
* Web based UI for configuration

Other cool stuff:
* Mute/Unmute Minidsp automation for things like turning off subs at night
* Home Assistant notifications to notify for events like loading/unloading BEQ was successful or failed
* Dry run and notification modes to verify BEQ profiles without actually loading them
* Built in support for Home Assistant and Minidsp

## Setup
> ⚠️ ⚠️ *Warning: You should really set a compressor on your minidsp for safety as outlined in the [BEQ forum post](https://www.avsforum.com/threads/bass-eq-for-filtered-movies.2995212/). I am not responsible for any damage* ⚠️ ⚠️
### Prerequisites
> ℹ  It is assumed you have the following tools working. Refer to their respective guides for installation help.
* Home Assistant (Optional)
* Plex or Jellyfin
* ezBEQ
* Minidsp (other DSPs may work but I have not tested them. If ezBEQ supports it, it should be work)

You can configure this to only load BEQ profiles, or do everything else besides BEQ. It is up to you.

### Docker Setup
> ℹ  If you need help deploying with Docker, refer to the [Docker documentation](https://docs.docker.com/get-docker/).
> ℹ  If you are using Jellyfin, read the Jellyfin specific instructions below

1) Deploy the latest version `ghcr.io/iloveicedgreentea/gowatchit:latest` to your preferred Docker environment
    * a docker-compose example is provided in the repo
2) You must mount a volume to `/data`
3) Configure the application via web ui -> `http://(your-server-ip):9999`
4) Set up your player with the instructions below
5) You can change logging timezone by setting the `TZ` env var to your desired timezone 

### Plex Specifics
1) get your player UUID(s) from `https://plex.tv/devices.xml` while logged in
  * https://plex.tv/devices.xml?X-Plex-Token=xyz [Getting plex token](https://support.plex.tv/articles/206721658-using-plex-tv-resources-information-to-troubleshoot-app-connections/)
2) Set up Plex to send webhooks to your server IP, port 9999, and the handler endpoint of `/plexwebhook`
    * e.g `(your-server-ip):9999/plexwebhook`
3) Whitelist your server IP in Plex so it can call the API without authentication. [Docs](https://support.plex.tv/articles/200890058-authentication-for-local-network-access/)
4) Add UUID(s) and user filters to the application config
5) Play a movie and check server logs. It should say what it loaded and you should see whatever options you enabled work

### Jellyfin Specifics

You must use the [official Jellyfin Webhooks plugin](https://github.com/jellyfin/jellyfin-plugin-webhook/tree/master) to send webhooks to this application.

1) Create a Generic webhook (NOT GenericForm)
2) Add http://(your-server-ip):9999/jellyfinwebhook as the url
3) Types:
  * PlaybackStart
  * PlaybackStopped
4) You can optionally add a user filter
5) Item types: Movies, Episodes

*note: playbackProgress is not supported because it is way too buggy and unreliable*

Configure the webhook in whatever way you want but it *must* include the following and in this order:

```json
{
  "DeviceId": "{{DeviceId}}",
  "DeviceName": "{{DeviceName}}",
  "ClientName": "{{ClientName}}",
  "UserId": "{{UserId}}",
  "ItemId": "{{ItemId}}",
  "ItemType": "{{ItemType}}",
  "NotificationType": "{{NotificationType}}",
{{#if_equals NotificationType 'PlaybackStop'}}
    "PlayedToCompletion": "{{PlayedToCompletion}}",
{{/if_equals}}
{{#if_equals NotificationType 'PlaybackProgress'}}
    "IsPaused": "{{IsPaused}}",
{{/if_equals}}
  "Year": "{{Year}}"
}
```

You need to also get the DeviceID from Jellyfin. The easiest way to do this is to enable webhooks as above, play something, go to your jellyfin logs, then search for this exact string - \"DeviceId\


Add the device ID to the application config and add it under the deviceuuidfilter field. If this is not set, HDMI sync detection will not work.

#### Generate API Key

1) Navigate to the dashboard
2) Click on“API Keys” under “Advanced” 
3) Click “Create”
4) Add API Key to the application config

### Emby

Emby support is added by coincidence due to Jellyfin support (its the same API codebase). I have not tested it and do not plan to. Configure it the same way as Jellyfin and it might work but I cannot guarantee it.

### Non-Docker Setup
I don't recommend this as it is more work and you will need to set up systemd or something to keep it running. I don't provide support for this method but if you know what you are doing, it is very easy to build the binary and run it.

TLDR: `make build`

## Usage

### Web UI
The web UI is the primary way to configure this application. It is available at `http://(your-server-ip):9999`

It will automatically restart the application when you save.

Each section has an enable/disable toggle. If you disable a section, it will not be used. For example, if you disable BEQ, it will not load BEQ profiles. If you disable MQTT, it will not send MQTT messages.

### General Usage
This application will load BEQ profiles automatically when you play something in Plex. It will also set volume, lights, and mute/unmute minidsp if you enable those options. The application itself is not controlling things like lights but relies on Home Assistant to perform the action via MQTT. In theory, you could use any home automation system but Home Assistant is the only one officially supported but anything that can receive MQTT messages should work.

## Home Assistant Quickstart

### Handlers
`/plexwebhook`

`/jellyfin` 

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

One use case is to mute the subs at night. You can use the time integration to trigger this at a certain time or with a button press.

### Config
The only supported way to configure this is via the web UI. You can dump the current config via the `/config` endpoint.

### Plex Authentication
You must whitelist your server IP in "List of IP addresses and networks that are allowed without auth"

Why? Plex refuses to implement client to server authentication and you must go through their auth servers. I may eventually implement their auth flow but it is not a priority.

### Logs
`/logs`
It will return the current logs as of the last request. It will not stream logs. You can use this to get logs for debugging. Refresh the page to get the latest logs.

### Debugging
These are environment variables you can set to get more info

`LOG_LEVEL=debug` to have it print out debug logs while running

`SUPER_DEBUG=true` for each line to also have a trace to its call site and line number

## How BEQ Support Works
On play and resume, it will load the profile. On pause and stop, it will unload it (so you don't forget to). It has some logic to cache the profile so if you pause and unpause, the profile will get loaded much faster as it skips searching the DB and stuff. 

If enabled, it will also send a notification to Home Assistant via Notify so you can send an alert to your phone for example. 

For safety, the application tries to unload the profile when it loads up each time in case it crashed or was killed previously, and will unload before playing anything so it doesn't start playing something with the wrong profile. 

### Matching
The application will search the catalog and match based on codec (Atmos, DTS-X, etc), title, year, TMDB, and edition. I have tested with multiple titles and everything matched as expected.

> ⚠️ *If you get an incorrect match, please open a github issue with the full log output and expected codec and title*

Jellyfin may have some issues matching as I have found it will sometimes just not return a TMDB. This has nothing to do with me. Jellyfin is generally just quite buggy. There is a configuration option that you should probably enable in the Jellyfin section which lets you skip TMDB matching. It will instead use the title name which could be prone to false negatives. 

### Editions

This application will do its best to match editions. It will look for one of the following:
1) Plex edition metadata. Set this from your server in the Plex UI
2) Looking at the file name if it contains `Unrated, Ultimate, Theatrical, Extended, Director, Criterion` (case insensitive)

There is no other reliable way to get the edition. If an edition is not matched, BEQ will fail to load for safety reasons (different editions have different masterings, etc). If a BEQCatalog entry has a blank edition, then edition will not matter and it will match based on the usual criteria.

If you want to match an edition when none is found locally, you can enable Loose Edition Matching. It will match if:

1) We sent a blank edition in the request
2) BEQ catalogue returns an edition

This is only useful if you have issues matching editions.

### HDMI Sync Automation
This application supports automatically waiting until HDMI sync is complete. 

Have you ever started something in Plex only to hear audio but see a black screen for 10 seconds? Then everyone in your theater makes fun of you and you cry yourself to sleep? This application will prevent that. 

It supports two ways to get this info currently: my Envy integration, or you can pass in seconds to wait. Using time is the simplest option. Measure how many seconds it takes from you pressing play to the video signal appearing on your screen. This is the input you use for "time".

Check the configuration UI for details on what to input.

If using Plex, you MUST get the Player Machine Identifier like so:

1) Play something on your desired player (like a shield)
2) `curl "http://(player IP):32500/resources"`
    * Note this is *NOT THE SERVER IP!* and *only works while something is actively playing*
3) Copy the `machineIdentifier` value
4) Add this to the Player Machine Identifier field exactly as presented
5) Add the player IP to the player IP field
6) Assign your player a static IP via your router or DHCP server

*Jellyfin HDMI sync is not implemented yet*

### Audio stuff
Here are some examples of what kind of codec tags Plex will have based on file metadata

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

## Help
First check your logs to see whats happening via the `/logs` endpoint.

If you need help or support due to an error or bug, you must file an issue. If you have a general question, you can ask in the Discussions tab or the AVS Forum post (linked as website above)