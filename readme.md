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
</p>

---

## Table of Contents
- [Features](#features)
- [Setup](#docs/setup.md)
- [Usage](#usage)
- [Home Assistant Quickstart](#Home-Assistant-Quickstart)
- [How BEQ Support Works](#docs/beq.md)
- [Help](#help)

## Features

Main features:
* Load/unload BEQ profiles automatically, without user action and the correct codec/edition detected
* HDMI Sync detection and automation (pause while HDMI is syncing so you don't sit embarrassed with a audio playing to a black screen in front of your friends)
* Web based UI

Players Supported:
* Kodi (through Home Assistant)
* Plex (Webhooks and through Home Assistant)
* Jellyfin (through Home Assistant)
* Emby (may work due to jellyfin support, no support given and not tested)
* Apple TV and technically any player that exposes the correct metadata (title, year, codec, edition, tmdb)

Other cool stuff:
* Mute/Unmute Minidsp
* Home Assistant notifications to notify for events like loading/unloading BEQ was successful or failed
* Dry run and notification modes to verify BEQ profiles without actually loading them
* Built in support for Home Assistant and Minidsp
* API to get BEQ status

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

## Help
First check your logs to see whats happening via the `/logs` endpoint.

If you need help or support due to an error or bug, you must file an issue. If you have a general question, you can ask in the Discussions tab or the AVS Forum post (linked as website above)