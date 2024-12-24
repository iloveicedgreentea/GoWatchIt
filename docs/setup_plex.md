## Table of Contents
- [Back To Main](../readme.md)
- [Unraid](./setup_unraid.md)
- [Plex](./setup_plex.md)
- [Jellyfin](./setup_jellyfin.md)
- [HDMI Sync](./setup_hdmi_sync.md)
- [Home Assistant (WIP)](./setup_homeassistant.md)


## Plex
1) get your player UUID from `https://plex.tv/devices.xml` while logged in (you may need to play something)
  * https://plex.tv/devices.xml?X-Plex-Token=xyz [Getting plex token](https://support.plex.tv/articles/206721658-using-plex-tv-resources-information-to-troubleshoot-app-connections/)
2) Set up Plex to send webhooks to your server IP, port 9999, and the handler endpoint of `/webhook`
    * e.g `(your-server-ip):3000/webhook`
3) Whitelist your server IP in Plex so it can call the API without authentication. [Docs](https://support.plex.tv/articles/200890058-authentication-for-local-network-access/)
4) Add UUID and user filters to the application config
5) Play a movie and check server logs. It should say what it loaded
