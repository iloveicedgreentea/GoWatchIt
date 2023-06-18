# Plex Webhook Automation

*In beta testing. Please monitor ezBEQ profiles*

*Please read the readme start to finish*

## Features

* Load/unload profiles automatically, no user action needed, correct codec detected
* Detect aspect ratio and send command to HA to adjust accordingly
  * Also supports using my MadVR Envy Home Assistant integration 
* Set Master Volume based on media type (movie, TV, etc)
* Trigger lights when playing or stopping automatically
* Mute/Unmute minidsp for night mode/WAF
* Mobile notifications (via HA) to notify for events like loading/unloading BEQ was successful or failed
* Dry run and notification modes to verify BEQ profiles without actually loading them
* All options are highly configurable with hot reload 
* Built in support for Home Assistant and Minidsp webhooks (e.x mute on)

*note: all communication to HA is done via MQTT so you will need to set this up*

I wrote this to be modular and extensible so adding additional listeners is simple. 

Feel free to make PRs or feature requests

## Usage
This tool is web API based. It is extensible by adding "handlers" which are listener endpoints for any function. 

### Handlers
`/plexwebhook`
This endpoint is where you should tell Plex to send webhooks to. It automatically processes them. No further action is needed. This handler does most of the work - Loading BEQ, aspect ratio, lights, volume, etc

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
  # unmute
  - service: rest_command.minidsp
    data:
      command: "off"
```

Using the above you can automate the mute and unmute of your minidsp with anything. I personally use Harmony and trigger this via Emulated Roku. Hold button to mute all subs and lower volume, press same button to unmute and reset volume.

## Setup
Note: this assumes you have ezBEQ, Plex, and HomeAssistant working. Refer to their respective guides for installation help.

You don't strictly need HA and you can use your own systems but I recommend HA.

0) Create `config.json` and set the values appropriately. See below.
1) Either pull `ghcr.io/iloveicedgreentea/plex-webhook-automation:master` or build the binary directly
    * if you deploy a container, mount config.json to a volume called exactly `/config.json`
2) Set up Plex to send webhooks to your server IP, `listenPort`, and the handler endpoint
3) Whitelist your server IP in Plex so it can call the API without authentication. Plex refuses to implement local server auth, so I don't want to implement their locked-in auth method that has historically had outages.
4) Play a movie and check server logs. It should say what it loaded and you should see whatever options you enabled work.
5) Add your UUID to the config.json so it filters by device
6) The app should detect the change and reload itself. If not, restart it.

You should deploy this as a container, systemd unit, etc. 

*side note: you should really set a compressor on your minidsp for safety as outlined in the BEQ forum post*


### MQTT
For flexibility, this uses MQTT to send commands. This is so you can decide what to do with that info. You will need to set MQTT up. Detailed instructions here https://www.home-assistant.io/integrations/mqtt/
  
1) Install mosquitto mqtt add on
2) Install mqtt integration
3) Set up your topics in HA and the tool's config
4) Set up Automations in HA based on the payloads of MQTT

### Payloads
In your Automations, you can action based on these payloads.

#### Lights
```json
{
    "state": "on" || "off"
}
```

#### Master Volume
```json
{
    "type": "movie" || "episode"
}
```

#### Aspect Ratio
```json
{
    "aspect": "2.4" || "2.2" || "1.85" || "1.78"
}
```

### Masking System Support

You can use the default (IMDB) or a MadVR Envy. IMDB works fine but they are very hostile to scraping so there is a chance it may fail, but I tried to add retries for that. 


*Note: if you enable madvr support, you must set up an Automation triggered by MQTT, topic is topicAspectratioMadVrOnly. Run you actions for masking system in that automation. The payload does not matter as its read from the envy. I recommend delaying reading the attribute by 12 seconds or so until the envy scales the display correctly and the attribute changes*

Here is an automation which uses MQTT and Envy attributes ([via my Envy integration](https://github.com/iloveicedgreentea/madvr-envy-homeassistant)). Modify to your needs. My masking system is set up for CIH so I mask off beyond 17:9. 

```yaml
alias: Envy - MQTT - Masking system
description: >-
  Trigger masking if it changed, but not within 5 min so alternating scenes
  don't trigger
trigger:
  - platform: mqtt
    topic: theater/envy/aspectratio
condition: []
action:
  - delay:
      hours: 0
      minutes: 0
      seconds: 12
      milliseconds: 0
  - if:
      - condition: numeric_state
        entity_id: remote.envy
        attribute: aspect_ratio
        above: 0
        below: 1.89
    then:
      - service: switch.turn_on
        data: {}
        target:
          entity_id: switch.masking_down
      - delay:
          hours: 0
          minutes: 0
          seconds: 35
          milliseconds: 0
      - service: switch.turn_off
        data: {}
        target:
          entity_id: switch.masking_down
  - if:
      - condition: numeric_state
        entity_id: remote.envy
        attribute: aspect_ratio
        below: 10
        above: 1.88
    then:
      - service: switch.turn_on
        data: {}
        target:
          entity_id: switch.masking_up
      - delay:
          hours: 0
          minutes: 0
          seconds: 35
          milliseconds: 0
      - service: switch.turn_off
        data: {}
        target:
          entity_id: switch.masking_up
mode: single

```

### HA Quickstart
Here is an example of an automation to change lights based on MQTT.

Assuming you have the following sensor:
```yaml
mqtt:
  sensor:
    - name: "lights"
      state_topic: "theater/lights/front"
      value_template: "{{ value_json.state }}"
```

This will turn the light(s) on/off depending on the state of the sensor, state is changed by a message sent to the topic

```yaml
alias: MQTT - Theater Lights
description: Trigger lights when mqtt received, depending on state
trigger:
  - platform: mqtt
    topic: theater/lights/front
condition: []
action:
  - if:
      - condition: state
        entity_id: sensor.lights
        state: "on"
    then:
      - service: light.turn_on
        data: {}
        target:
          entity_id: light.caseta_r_wireless_in_wall_dimmer
  - if:
      - condition: state
        entity_id: sensor.lights
        state: "off"
    then:
      - service: light.turn_off
        data: {}
        target:
          entity_id: light.caseta_r_wireless_in_wall_dimmer
mode: single
```

### Config

create file named config.json, paste this in, remove the comments after

```json
{
    "homeAssistant": {
        "url": "http://123.123.123.123",
        "port": "8123",
        "enabled": true,
        // get a token from your user profile
        "token": "ey.xyzjwt",
        // Trigger functions to change the following
        "triggerAspectRatioChangeOnEvent": true,
        "triggerLightsOnEvent": true,
        "triggerAvrMasterVolumeChangeOnEvent": true
    },
    // all communication to HA is done via MQTT. Set up automations to run scripts
    "mqtt": {
        // url to broker and user/pass to use. Set up mosquitto via HA add on then add an HA user
        "url": "tcp://123.123.123.123:1883",
        "username": "sdf",
        "password": "123",
        // these are arbitrary strings
        "topicLights": "theater/lights/front",
        "topicVolume": "theater/denon/volume",
        "topicAspectratio": "theater/jvc/aspectratio"
    },
    "plex": {
        // your main owner account, will filter webhooks so others dont trigger
        // leave blank if you dont want to filter on accounts
        "ownerNameFilter": "PLEX_OWNER_NAME to filter events on",
        // filter based on device UUID so only the client you want triggers things, or leave blank
        // Must be UUID. Easy way to get it is running this in debug mode and then play a movie
        "deviceUUIDFilter": "",
        "url": "http://xyz",
        "port": "32400",
        // if you enable trailers before movies, it can process it like turn off lights. no BEQ 
        "enableTrailerSupport": true || false
    },
    "ezbeq": {
        // note this will use slot1/config1. I don't see a good reason to support multiple slots since this is event driven
        "url": "http://xyz",
        "port": "8080",
        "enabled": true,
        // support BEQ for TV shows also, some exist
        "enableTvBeq": true,
        // will log what it will do, but will not load BEQ profiles
        "dryRun": false,
        // some BEQ catalogs have negative MV adjustment. Recommend to true unless you really like bass, can cause damage
        "adjustMasterVolumeWithProfile": true,
        // Trigger HA to notify you when it loads so you can double check stuff. Will also trigger with dryrun enabled
        "notifyOnLoad": true,
        // name of the endpoint in HA to send notification to. Look at the notify service in HA to see endpoints
        "notifyEndpointName": "mobile_app_iphone",
        // which author you want. none or blank will find the best match according to ezbeq application
        "preferredAuthor": "aron7awol" || "mobe1969" || "none" || "",
        // slots you want to apply beq configs. minidsp 2x4hd has four PRESET slots. Not tested on anything but 2x4hd
        "slots": [1],
        // use an IP enabled Denon AVR to get the codec instead of querying plex. This is faster and more reliable
        "useAVRCodecSearch": true,
        "DenonIP": "",
        "DenonPort": "23",
    },
    "main": {
        "listenPort": "9999"
    }
}
```

### Authentication
You must whitelist your server IP in "List of IP addresses and networks that are allowed without auth"

Why? Plex refuses to implement client to server authentication and you must go through their auth servers. I don't want to do that so this is my form of protest.

A local attacker hijacking my server and sending commands to Plex is not remotely in my threat model. 

### Debug mode
`export LOG_LEVEL=debug` to have it print out debug logs

`export SUPER_DEBUG=true` for each line to have a trace to its call site and line number

If using a container you can set the above as environment variables. 

## How BEQ Support Works
On play and resume, it will load the profile. On pause and stop, it will unload it (so you don't forget). It has some logic to cache the profile so if you pause and unpause, the profile will get loaded much faster as it skips searching the DB and stuff. 

If enabled, it will also send a notification to Home Assistant via Notify. 

For safety, the tool tries to unload the profile when it loads up each time in case it crashed or was killed previously, and will unload before playing anything so it doesn't start playing something with the wrong profile. 

### Matching
The tool will search the catalog and match based on codec (Atmos, DTS-X, etc), title, year, and edition. I have tested with multiple titles and everything matched as expected.

*If you get an incorrect match, please open a github issue asap*

### Editions

This tool will do its best to match editions. It will look for one of the following:
1) Plex edition metadata. Set this from your server in the UI
2) Looking at the file name if it contains `Unrated, Ultimate, Theatrical, Extended, Director, Criterion`

There is no other reliable way to get the edition. If an edition is not matched, BEQ will fail to load for safety. If a BEQCatalog entry has a blank edition, then edition will not matter and it will match based on the usual criteria.

If you find repeated match failures because of editions, open a github issue with debug logs of you triggering `media.play`

## Building Binary
GOOS=xxxx make build

## Tech Nerd Stuff / Development

This uses a modular architecture via handlers. The main action points are `func ProcessWebhook` which processes and sends the payload to a 
channel processed by `func PlexWorker` which runs in the background. 

`func eventRouter` uses flags and switches to determine what to do. Additional actions can easily be added here. The actionable functions run as coroutines for maximum speed. Going from play to lights off is instantanous and aspect ratio detection takes about 1.5 seconds.

`ezbeq`, `plex`, amd `homeassistant` packages have reusable clients so their functions can easily be used by other handlers.

`logger` is a reusable logging package which sets some nice defaults and stuff like debug mode

`models` contains all the structs needed for serialization/deserialization

### Adding handlers
Check `main.go` for how to implement a new handler. Call `mux.Handle()` to add the new handler

Variable aspect movies will use the widest aspect listed in IMDB

### Audio stuff
Here are some examples of what kind of codec tags Plex will spit out based on file metadata

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