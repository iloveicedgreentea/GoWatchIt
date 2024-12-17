


#### Plex
1) get your player UUID(s) from `https://plex.tv/devices.xml` while logged in
  * https://plex.tv/devices.xml?X-Plex-Token=xyz [Getting plex token](https://support.plex.tv/articles/206721658-using-plex-tv-resources-information-to-troubleshoot-app-connections/)
2) Set up Plex to send webhooks to your server IP, port 9999, and the handler endpoint of `/plexwebhook`
    * e.g `(your-server-ip):9999/plexwebhook`
3) Whitelist your server IP in Plex so it can call the API without authentication. [Docs](https://support.plex.tv/articles/200890058-authentication-for-local-network-access/)
4) Add UUID(s) and user filters to the application config
5) Play a movie and check server logs. It should say what it loaded and you should see whatever options you enabled work
