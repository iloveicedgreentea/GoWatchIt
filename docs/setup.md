## Table of Contents
- [Back To Main](../readme.md)
- [Unraid](./setup_unraid.md)
- [Plex](./setup_plex.md)
- [Jellyfin](./setup_jellyfin.md)
- [HDMI Sync](./setup_hdmi_sync.md)
- [Home Assistant (WIP)](./setup_homeassistant.md)


## Setup
> ⚠️ ⚠️ *Warning: You should really set a compressor on your minidsp for safety as outlined in the [BEQ forum post](https://www.avsforum.com/threads/bass-eq-for-filtered-movies.2995212/). I am not responsible for any damage* ⚠️ ⚠️
### Prerequisites
> ℹ  It is assumed you have the following tools working. Refer to their respective guides for installation help.
* Home Assistant (Optional)
* Plex or Home Assistant supported player
* ezBEQ
* Minidsp (other DSPs may work but I have not tested them. If ezBEQ supports it, it should be work)

You can configure this to only load BEQ profiles, or do everything else besides BEQ. It is up to you.

> ℹ  If you are using MSO, make sure to use the BEQ friendly output export. If you have PEQ on your inputs, BEQ will overwrite them. If you have shared gain, make sure to disable master volume adjustment

### Docker Setup
> ℹ  If you need help deploying with Docker, refer to the [Docker documentation](https://docs.docker.com/get-docker/).
> ℹ  If you are using Jellyfin, read the Jellyfin specific instructions below

1) Deploy the latest version `ghcr.io/iloveicedgreentea/gowatchit:latest` to your preferred Docker environment
    * a docker-compose example is provided in the repo
2) You must mount a volume to `/data`
3) Configure the application via web ui -> `http://(your-server-ip):9999`
4) Set up your player with the instructions below
5) You can change logging timezone by setting the `TZ` env var to your desired timezone 

### Triggering
There are two ways to trigger an action:

1) Home Assistant
2) Plex webhooks

Home Assistant is recommended due to speed and compatibility.

#### Home Assistant
You need to create a rest_commannd and an automation

Rest command
1) TBD # TODO: add rest command example

Automation
1) Create a new blank automation
2) Set the triggers to be BOTH your client playing AND stopping
3) Give each event a unique ID
4) For actions, create a Choose block
5) Have one option get triggered by the play event and the other by the stop event

Example
```yaml
mode: single
triggers:
  - trigger: state
    entity_id:
      - media_player.plex
    to: playing
    id: playing
  - trigger: state
    entity_id:
      - media_player.plex
    from: playing # notice this triggers when it changes FROM playing you can also change this to stop and pause separately
    id: stopped
conditions: []
actions:
  - choose:
      - conditions:
          - condition: trigger
            id:
              - playing
        sequence:
          - action: rest_command.playing
            metadata: {}
            data: {}
      - conditions:
          - condition: trigger
            id:
              - stopped # you can add additional commands for different states like pause
        sequence:
          - action: rest_command.stopped
            metadata: {}
            data: {}
```