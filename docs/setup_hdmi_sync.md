## Table of Contents
- [Back To Main](../readme.md)
- [Unraid](./setup_unraid.md)
- [Plex](./setup_plex.md)
- [Jellyfin](./setup_jellyfin.md)
- [HDMI Sync](./setup_hdmi_sync.md)
- [Home Assistant (WIP)](./setup_homeassistant.md)


This is not enabled yet all docs are WIP

## HDMI Sync Automation
This application supports automatically waiting (pause/play) until HDMI sync is complete. 

Have you ever started something in Plex only to hear audio but see a black screen for 10 seconds? Then everyone in your theater makes fun of you and you cry yourself to sleep? This application will prevent that. 

It supports two ways to get this info currently: a home assistant entity or defining seconds to wait. 

> ℹ If you use Plex webhooks, sync will be imprecise due to how unreliable Plex webhooks are. Switch to Home Assistant triggering

### Home Assistant
* Integration must have a remote or media_player entity that exposes a binary sensor for signal status. For example, my [madVR integration](https://www.home-assistant.io/integrations/madvr/#binary-sensor)

### Time

Using time is the simplest option. Measure how many seconds it takes from you pressing play to the video signal appearing on your screen. This is the input you use for "time".


Check the configuration UI for details on what to input.

## Player Setup
You can use home assistant or Plex native API. I recommend just using home assistant as it is easier to set up.

### Non-Plex
You can use supply the home assistant media_player entity to send pause/play commands to

### Plex Specifics
If using Plex, you MUST get the Player Machine Identifier like so:

1) Play something on your desired player (like a Shield)
2) `curl "http://(player IP):32500/resources"`
    * Note this is *NOT THE SERVER IP!* and *only works while something is actively playing*
3) Copy the `machineIdentifier` value
4) Add this to the Player Machine Identifier field exactly as presented
5) Add the player IP to the player IP field
6) Assign your player a static IP via your router or DHCP server
