# Home Theater Automation Server

*Please read the readme start to finish*

## Features

This has become more of a general purpose home theater automation server. It is centered around Plex, but I am adding support for generic sources of info as well.

Meat & Potatos:
* Load/unload BEQ profiles automatically, no user action needed, correct codec detected
* Set Volume based on media type (movie, TV, etc)
* Trigger lights when playing, pausing, or stopping automatically
* HDMI Sync detection and automation (pause while HDMI is syncing so you don't sit embarrassed with a audio playing to a black screen)
* Web based UI for configuration

Other cool stuff:
* Mute/Unmute Minidsp automation for things like turning off subs at night
* Detect aspect ratio and send command to HA to adjust accordingly
  * Also supports using my MadVR Envy Home Assistant integration 
* Various MQTT sensors for volume control, lights, mute status, and current BEQ profile
* Mobile notifications (via HA) to notify for events like loading/unloading BEQ was successful or failed
* Dry run and notification modes to verify BEQ profiles without actually loading them
* Built in support for Home Assistant and Minidsp

## Setup
Note: this assumes you have ezBEQ, Plex, and HomeAssistant working. Refer to their respective guides for installation help.

You don't strictly need HA and you can use your own systems but I recommend HA.

Simple Way:

1) Deploy via Docker -> `ghcr.io/iloveicedgreentea/plex-webhook-automation:$version` 
2) Set up config via web ui -> `http://(server-ip):9999`
   * can get UUID from `https://plex.tv/devices.xml`
3) Set up Plex to send webhooks to your server IP, `listenPort`, and the handler endpoint of `/plexwebhook`
    * `(server-ip):9999/plexwebhook`
3) Whitelist your server IP in Plex so it can call the API without authentication. 
4) Play a movie and check server logs. It should say what it loaded and you should see whatever options you enabled work.
5) App will restart within 5 seconds when config is changed

Manual Way:

0) Create `config.json` and set the values appropriately. See below.
1) Either pull `ghcr.io/iloveicedgreentea/plex-webhook-automation:$version` or build the binary directly
    * if you deploy a container, mount config.json to a volume called exactly `/config.json`
2) Set up Plex to send webhooks to your server IP, `listenPort`, and the handler endpoint of `/plexwebhook`
    * `(server-ip):9999/plexwebhook`
3) Whitelist your server IP in Plex so it can call the API without authentication. Plex refuses to implement local server auth with an API, so I don't want to implement their locked-in auth method that has historically had outages (8/30/23 is the latest one by the way).
4) Add your UUID to the config.json so it filters by device
    * can get UUID from `https://plex.tv/devices.xml` or run the tool and play something, check the logs
5) Play a movie and check server logs. It should say what it loaded and you should see whatever options you enabled work.
6) The app should detect the change and reload itself. This doesn't usually work with Docker, so will need a manual restart.

You should deploy this as a container, systemd unit, etc. 

*side note: you should really set a compressor on your minidsp for safety as outlined in the BEQ forum post, outside the scope here but you have been warned, I am not responsible for any damages*


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
  # unmute subs
  - service: rest_command.minidsp
    data:
      command: "off" (or "on")
```

Using the above you can automate the mute and unmute of your minidsp with any automation source.
 
You can then do cool stuff like create a binary sensor to show the state of the subs based on the MQTT topic

### MQTT
For flexibility, this uses MQTT to send commands. This is so you can decide what to do with that info. You will need to set MQTT up. Detailed instructions here https://www.home-assistant.io/integrations/mqtt/
  
1) Install mosquito mqtt add on
2) Install mqtt integration
3) Set up your topics in HA and the tool's config
4) Set up Automations in HA based on the payloads of MQTT

Used Topics:
Aspect ratio
Current BEQ Profile
Lights
Minidsp mute status
Volume
Playing status


Here are some sensor examples

```yaml
mqtt:
  binary_sensor:
    - name: "subs_muted"
      state_topic: "theater/subs/status"
      payload_on: "false"
      payload_off: "true"
    - name: "plex_playing"
      state_topic: "theater/plex/playing"
      payload_on: "true"
      payload_off: "false"
  sensor:
    - name: "lights"
      state_topic: "theater/lights/front"
      value_template: "{{ value_json.state }}"
    - name: "volume"
      state_topic: "theater/denon/volume"
      value_template: "{{ value_json.type }}"
    - name: "aspectratio"
      state_topic: "theater/jvc/aspectratio"
      value_template: "{{ value_json.aspect }}"
    - name: "beq_current_profile"
      state_topic: "theater/beq/currentprofile"


```

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
    "aspect": "2.4" || "2.2" || "1.85" || "1.78" || etc
}
```

All other payloads are "true" or "false" so create an MQTT binary sensor.

### Masking System Support

There are two ways to do masking. One is using this tool. The other is using my MadVR Envy Home Assistant integration and using the Envy's aspect ratio (Aspect dec or Aspect int) attribute in an automation. The second way is most useful if you have a curtain like masking system, not a drop down one.

If using this tool, you can use the default source (IMDB) or MadVR Envy. IMDB works fine but they are very hostile to scraping so there is a chance it may fail because they intentionally change their HTML to break scrapers, but I tried to add retries for that. It is provided as a best effort but IMDB can be unreliable, if not wrong. Variable Aspect Ratio movies like Interstellar will use the widest AR reported by IMDB as it is the most likely option. 

If you use an Envy with my method, you will get real time masking adjustments. I built my own infinite ratio CIH masking system.

*Note: if you enable madvr support, you must set up an Automation triggered by MQTT, topic needs to be named topicAspectratioMadVrOnly. Run your actions for masking system in that automation. The payload does not matter as its read from the envy.*

Here is an automation which Envy attributes ([via my Envy integration](https://github.com/iloveicedgreentea/madvr-envy-homeassistant)). Modify to your needs. My masking system is set up for CIH so I mask off beyond 16:9. 

Example Automation:
```yaml
alias: "Envy: Masking System"
description: ""
trigger:
  - platform: state
    entity_id:
      - remote.envy
    attribute: aspect_dec
condition:
  - condition: state
    entity_id: input_boolean.wide1
    state: "off"
action:
  - if:
      - condition: numeric_state
        entity_id: remote.envy
        attribute: aspect_dec
        above: 0
        below: 1.79
      - condition: state
        entity_id: input_boolean.wide1
        state: "off"
        enabled: false
    then:
      - service: switch.turn_on
        data: {}
        target:
          entity_id: switch.masking_close
      - stop: ""
    alias: "1.78"
  - if:
      - condition: numeric_state
        entity_id: remote.envy
        attribute: aspect_dec
        above: 2.33
        below: 6
      - condition: state
        entity_id: input_boolean.wide1
        state: "off"
        enabled: false
    then:
      - service: switch.turn_on
        data: {}
        target:
          entity_id: switch.masking_open
      - stop: ""
        error: false
    alias: scope
  - if:
      - condition: numeric_state
        entity_id: remote.envy
        attribute: aspect_dec
        above: 1.78
        below: 1.87
      - condition: state
        entity_id: input_boolean.wide1
        state: "off"
        enabled: false
    then:
      - service: switch.turn_on
        data: {}
        target:
          entity_id: switch.masking_1_85_1
      - stop: ""
    alias: "1.85"
  - if:
      - condition: numeric_state
        entity_id: remote.envy
        attribute: aspect_dec
        above: 1.86
        below: 1.91
      - condition: state
        entity_id: input_boolean.wide1
        state: "off"
        enabled: false
    then:
      - service: switch.turn_on
        data: {}
        target:
          entity_id: switch.masking_1_9_1
      - stop: ""
    alias: "1.9"
  - if:
      - condition: numeric_state
        entity_id: remote.envy
        attribute: aspect_dec
        above: 1.9
        below: 2.1
      - condition: state
        entity_id: input_boolean.wide1
        state: "off"
        enabled: false
    then:
      - service: switch.turn_on
        data: {}
        target:
          entity_id: switch.masking_2_0_1
      - stop: ""
    alias: "2.0"
  - if:
      - condition: numeric_state
        entity_id: remote.envy
        attribute: aspect_dec
        above: 2.1
        below: 2.33
      - condition: state
        entity_id: input_boolean.wide1
        state: "off"
        enabled: false
    then:
      - service: switch.turn_on
        data: {}
        target:
          entity_id: switch.masking_open
      - stop: ""
    alias: "2.2"
mode: queued
max: 10
```

I have the condition set because the Envy internally sets things to 16:9 while it is syncing which can cause false positives for masking.

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
mode: queued
max: 10
```

### Config
You can configure this with the web UI or manually. The app will be restarted after changing any values.

All fields are required unless otherwise stated

If editing manually, create a file named config.json, paste this in, remove the comments after

```json

    // note the case is lowercase
    "homeassistant": {
        "enabled": true,
        "url": "http://123.123.123.123",
        "port": "8123",
        // get a token from your user profile
        "token": "ey.xyzjwt",
        // trigger functions to change the following
        "triggeraspectratiochangeonevent": true,
        "triggerlightsonevent": true,
        "triggeravrmastervolumechangeonevent": true,
        // optional if using hdmi sync via an entity. the name of the remote entities, or a binary sensor. If binary sensor, set source.source to "binary_sensor". See readme
        "remoteentityname": "jvc",
        // optional: set only if plex.enabled is false. names of scripts which call Home assistant scripts to play/pause/stop
        "playscriptname": "",
        "pausescriptname": "",
        "stopscriptname": ""
    },
    // all communication to ha is done via mqtt. set up automations to run scripts
    "mqtt": {
        // url to broker and user/pass to use. set up mosquito via ha add on then add an ha user
        "url": "tcp://123.123.123.123:1883",
        "username": "sdf",
        "password": "123",
        // these are arbitrary strings
        "topiclights": "theater/lights/front",
        "topicvolume": "theater/denon/volume",
        "topicaspectratio": "theater/jvc/aspectratio",
        // will publish the current profile here
        "topicbeqcurrentprofile": "theater/beq/currentprofile",
        // write mute status
        "topicminidspmutestatus": "theater/minidsp/mutestatus",
        // if tool is playing or not
        "topicplayingstatus": "theater/plex/playing",
    },
    "plex": {
        // if you don't use plex
        "enabled": true,
        // your main owner account, will filter webhooks so others don't trigger
        // leave blank if you don't want to filter on accounts
        "ownernamefilter": "plex_owner_name to filter events on",
        // filter based on device uuid so only the client you want triggers things, or leave blank
        // must be uuid. easy way to get it is playing anything and searching logs for 'got a request from uuid:'
        // or check the devices. You can split it with a comma for multiple devices
        "deviceuuidfilter": "",
        "url": "http://xyz",
        "port": "32400",
        // if you enable trailers before movies, it can process it like turn off lights. no beq 
        "enabletrailersupport": true || false,
        // optional if not using hdmi sync get this from "http://(player ip):32500/resources". Required for pausing/playing plex . check readme
        "playermachineidentifier": "uuid",
        "playerip": "xxx.xxx.xxx.xxx"
    },
    "ezbeq": 
        "url": "http://xyz",
        "port": "8080",
        "enabled": true,
        // support beq for tv shows also, some existyhgt
        "": true,
        // will log what it will do, but will not load beq profiles
        "dryrun": false,
        // some beq catalogs have negative mv adjustment. recommend to true unless you really like bass, can cause damage
        "adjustmastervolumewithprofile": true,
        // trigger ha to notify you when it loads so you can double check stuff. will also trigger with dryrun enabled
        "notifyonload": true,
        // name of the endpoint in ha to send notification to. look at the notify service in ha to see endpoints
        "notifyendpointname": "mobile_app_iphone",
        // which author to filter on. blank will find the best match according to ezbeq 
        "preferredauthor": "aron7awol" or "mobe1969" or "other supported author" or "",
        // slots you want to apply beq configs. minidsp 2x4hd has four preset slots. not tested on anything but 2x4hd
        "slots": [1],
        // use an ip enabled denon avr to get the codec instead of querying plex
        // requires a madvr envy for now
        // much slower but more accurate as it will get the actual codec playing
        // will also compare denon and plex to ensure correct codec is playing (sometimes plex will incorrectly transcode. might be a shield bug) (not ready)
        "useavrcodecsearch": false,
        // optional if not using above field
        "denonip": "",
        "denonport": "23",
        // tell plex to stop if the playing codec does not match expected like when it transcodes atmos for no reason
        // requires denon AVR above
        "stopplexifmismatch": true
    ,
    // what to use for signal source
    "signal": {
      // true if you want to pause plex until hdmi sync is done
      "enabled": true,
        // jvc, envy, or name of the binary sensor (see readme), or you can specify seconds to pause for like "13"
        "source": "jvc"
    },
    "main": {
        "listenport": "9999"
    }

```

### Authentication
You must whitelist your server IP in "List of IP addresses and networks that are allowed without auth"

Why? Plex refuses to implement client to server authentication and you must go through their auth servers. I don't want to do that so this is my form of protest.

A local attacker hijacking my server and sending commands to Plex is not a concern. 

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

*If you get an incorrect match, please open a github issue with the full log output and expected codec and title*

### Editions

This tool will do its best to match editions. It will look for one of the following:
1) Plex edition metadata. Set this from your server in the UI
2) Looking at the file name if it contains `Unrated, Ultimate, Theatrical, Extended, Director, Criterion`

There is no other reliable way to get the edition. If an edition is not matched, BEQ will fail to load for safety reasons (different mastering, etc). If a BEQCatalog entry has a blank edition, then edition will not matter and it will match based on the usual criteria.

If you find repeated match failures because of editions, open a github issue with debug logs of you triggering `media.play`

### HDMI Sync Automation
This tool supports automatically waiting until HDMI sync is done. Have you ever started something in Plex only to hear audio but see a black screen for 10 seconds? Then everyone you are watching a movie with makes fun of you and you cry yourself to sleep? This tool will prevent that. 

It supports four ways to get this info: my JVC integration, my Envy integration, a generic binary_sensor, or you can pass in seconds to wait.

If using the first two methods, you just need to install the integration, then set `remoteentityname` to the name of the remote entity. Set `signal.source` to either `jvc` or `envy`.

If using a binary_sensor, you need to create an automation which will set the state to `off` when there is NO SIGNAL and `on` when there is. Getting that data is up to you. Set the `signal` config to the name of the binary sensor (e.g signal, if the entity is binary_sensor.signal).

If using seconds, provide the number of seconds to wait as a string such as "13" for 13 seconds. You can time how long your sync time is and add it here. It will pause for that amount of time then continue playing.

You also must set `plex.playerMachineIdentifier` and `plex.playerIP`. To get this:
1) Play something on your desired player (like a shield)
2) `curl "http://(player IP):32500/resources"`
    * Note this is *NOT THE SERVER IP!* and *only works while something is actively playing*
3) Copy the `machineIdentifier` value
4) Add this to that config field exactly as presented

## Building Binary
GOOS=xxxx make build

## Tech Nerd Stuff / Development

This uses a modular architecture via handlers. The main action points are `func ProcessWebhook` which processes and sends the payload to a 
channel processed by `func PlexWorker` which runs in the background. 

`func eventRouter` uses flags and switches to determine what to do. Additional actions can easily be added here. The actionable functions run as coroutines for maximum speed. Going from play to lights off is instantaneous and aspect ratio detection takes about 1.5 seconds.

`ezbeq`, `plex`, amd `homeassistant` packages have reusable clients so their functions can easily be used by other handlers.

`logger` is a reusable logging package which sets some nice defaults and stuff like debug mode

`models` contains all the structs needed for serialization/deserialization

### Adding handlers
Check `main.go` for how to implement a new handler. 

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
