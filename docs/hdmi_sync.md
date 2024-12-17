
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

