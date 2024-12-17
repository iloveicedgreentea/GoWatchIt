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
<!-- TODO: add other pages here -->
- [Setup](./docs/setup.md)
- [Usage](#usage)
- [Home Assistant Quickstart](#Home-Assistant-Quickstart)
- [How BEQ Support Works](./docs/beq.md)
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
The web UI mainly serves to configure this application. It is available at `http://(your-server-ip):9999`

Each section has an enable/disable toggle. If you disable a section, it will not be used. For example, if you disable BEQ, it will not load BEQ profiles. Options will not be shown if the section is disabled.

You can also check application logs. It will fetch logs on the page automatically.

## Home Assistant Quickstart

### Endpoints

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
Look at the Logs view in the UI, query the `/logs` endpoint, check stdout, or read the log in `/data`

### Debugging
These are environment variables you can set to get more info

`LOG_LEVEL=debug` to have it print out debug logs while running
`LOG_FILE=true` to have it write logs to a file
`LOG_ENV=local` leave blank when using docker

## Help
First check your logs to see whats happening in the Logs view.

If you need help or support due to an error or bug, you must file an issue. If you have a general question, you can ask in the Discussions tab or the AVS Forum post (linked as website above)