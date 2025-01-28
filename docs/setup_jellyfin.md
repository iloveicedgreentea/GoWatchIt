## Table of Contents
- [Back To Main](../readme.md)
- [Unraid](./setup_unraid.md)
- [Plex](./setup_plex.md)
- [Jellyfin](./setup_jellyfin.md)
- [HDMI Sync](./setup_hdmi_sync.md)
- [Home Assistant (WIP)](./setup_homeassistant.md)


## Jellyfin

You must use the [official Jellyfin Webhooks plugin](https://github.com/jellyfin/jellyfin-plugin-webhook/tree/master) to send webhooks to this application.

1) Create a Generic webhook (NOT GenericForm)
2) Add http://(your-server-ip):3000/webhook as the url
3) Types:
  * PlaybackStart
  * PlaybackStopped
4) You can optionally add a user filter
5) Item types: Movies, Episodes

*note: playbackProgress is not supported because it is way too buggy, unreliable, unpredictable*

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
