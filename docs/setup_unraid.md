## Table of Contents
- [Back To Main](../readme.md)
- [Unraid](./setup_unraid.md)
- [Plex](./setup_plex.md)
- [Jellyfin](./setup_jellyfin.md)
- [HDMI Sync](./setup_hdmi_sync.md)
- [Home Assistant (WIP)](./setup_homeassistant.md)

## Unraid setup

1) Add Container
2) Repository - `ghcr.io/iloveicedgreentea/gowatchit:latest`
3) Network Type - `bridge`
4) Volumes:
   * Data volume - Container Path:` /data`
5) Ports:
   * Container Port: `3000`