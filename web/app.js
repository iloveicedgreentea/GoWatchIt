function populateFields(config) {
    // .checked or .value for each kind for each item in HTML

    // EZBeq
    document.getElementById('ezbeq-adjustmastervolumewithprofile').checked = config.ezbeq.adjustmastervolumewithprofile;
    document.getElementById('ezbeq-enabled').checked = config.ezbeq.enabled;
    document.getElementById('ezbeq-avrip').value = config.ezbeq.avrip;
    document.getElementById('ezbeq-dryrun').checked = config.ezbeq.dryrun;
    document.getElementById('ezbeq-enabletvbeq').checked = config.ezbeq.enabletvbeq;
    document.getElementById('ezbeq-notifyendpointname').value = config.ezbeq.notifyendpointname;
    document.getElementById('ezbeq-notifyonload').checked = config.ezbeq.notifyonload;
    document.getElementById('ezbeq-port').value = config.ezbeq.port;
    document.getElementById('ezbeq-preferredauthor').value = config.ezbeq.preferredauthor;
    const slotsArray = config.ezbeq.slots;
    slotsArray.forEach(slot => {
        document.getElementById(`slot${slot}`).checked = true;
    });
    document.getElementById('ezbeq-stopplexifmismatch').checked = config.ezbeq.stopplexifmismatch;
    document.getElementById('ezbeq-url').value = config.ezbeq.url;
    document.getElementById('ezbeq-useavrcodecsearch').checked = config.ezbeq.useavrcodecsearch;
    document.getElementById('ezbeq-avrbrand').value = config.ezbeq.avrbrand;


    // HomeAssistant
    document.getElementById('homeassistant-enabled').checked = config.homeassistant.enabled;
    document.getElementById('homeassistant-url').value = config.homeassistant.url;
    document.getElementById('homeassistant-port').value = config.homeassistant.port;
    document.getElementById('homeassistant-token').value = config.homeassistant.token;
    document.getElementById('homeassistant-triggerlightsonevent').checked = config.homeassistant.triggerlightsonevent;
    document.getElementById('homeassistant-triggeravrmastervolumechangeonevent').checked = config.homeassistant.triggeravrmastervolumechangeonevent;
    document.getElementById('homeassistant-remoteentityname').value = config.homeassistant.remoteentityname;
    document.getElementById('homeassistant-playscriptname').value = config.homeassistant.playscriptname;
    document.getElementById('homeassistant-pausescriptname').value = config.homeassistant.pausescriptname;
    document.getElementById('homeassistant-stopscriptname').value = config.homeassistant.stopscriptname;

    // MQTT
    document.getElementById('mqtt-enabled').checked = config.mqtt.enabled;
    document.getElementById('mqtt-url').value = config.mqtt.url;
    document.getElementById('mqtt-username').value = config.mqtt.username;
    document.getElementById('mqtt-password').value = config.mqtt.password;
    document.getElementById('mqtt-topiclights').value = config.mqtt.topiclights;
    document.getElementById('mqtt-topicvolume').value = config.mqtt.topicvolume;
    document.getElementById('mqtt-topicbeqcurrentprofile').value = config.mqtt.topicbeqcurrentprofile;
    document.getElementById('mqtt-topicminidspmutestatus').value = config.mqtt.topicminidspmutestatus;
    document.getElementById('mqtt-topicplayingstatus').value = config.mqtt.topicplayingstatus;

    // Plex
    document.getElementById('plex-enabled').checked = config.plex.enabled;
    document.getElementById('plex-url').value = config.plex.url;
    document.getElementById('plex-port').value = config.plex.port;
    document.getElementById('plex-ownernamefilter').value = config.plex.ownernamefilter;
    document.getElementById('plex-deviceuuidfilter').value = config.plex.deviceuuidfilter;
    document.getElementById('plex-playermachineidentifier').value = config.plex.playermachineidentifier;
    document.getElementById('plex-playerip').value = config.plex.playerip;
    document.getElementById('plex-enabletrailersupport').checked = config.plex.enabletrailersupport;
    // jellyfin
    document.getElementById('jellyfin-enabled').checked = config.jellyfin.enabled;
    document.getElementById('jellyfin-url').value = config.jellyfin.url;
    document.getElementById('jellyfin-port').value = config.jellyfin.port;
    document.getElementById('jellyfin-ownernamefilter').value = config.jellyfin.ownernamefilter;
    document.getElementById('jellyfin-deviceuuidfilter').value = config.jellyfin.deviceuuidfilter;
    document.getElementById('jellyfin-playermachineidentifier').value = config.jellyfin.playermachineidentifier;
    document.getElementById('jellyfin-userID').value = config.jellyfin.userID;
    document.getElementById('jellyfin-apiToken').value = config.jellyfin.apiToken;

    // Signal
    document.getElementById('signal-enabled').checked = config.signal.enabled;
    document.getElementById('signal-source').value = config.signal.source;
}

function buildFinalConfig() {
    const slotsArray = [];
    for (let i = 1; i <= 4; i++) {
        if (document.getElementById(`slot${i}`).checked) {
            slotsArray.push(i);
        }
    }
    const ezbeqConfig = {
        "adjustmastervolumewithprofile": document.getElementById('ezbeq-adjustmastervolumewithprofile').checked,
        "enabled": document.getElementById('ezbeq-enabled').checked,
        "avrip": document.getElementById('ezbeq-avrip').value,
        "dryrun": document.getElementById('ezbeq-dryrun').checked,
        "enabletvbeq": document.getElementById('ezbeq-enabletvbeq').checked,
        "notifyendpointname": document.getElementById('ezbeq-notifyendpointname').value,
        "notifyonload": document.getElementById('ezbeq-notifyonload').checked,
        "port": document.getElementById('ezbeq-port').value,
        "preferredauthor": document.getElementById('ezbeq-preferredauthor').value,
        "slots": slotsArray,
        "stopplexifmismatch": document.getElementById('ezbeq-stopplexifmismatch').checked,
        "url": document.getElementById('ezbeq-url').value,
        "avrbrand": document.getElementById('ezbeq-avrbrand').value,
        "useavrcodecsearch": document.getElementById('ezbeq-useavrcodecsearch').checked
    };
    const homeAssistantConfig = {
        "enabled": document.getElementById('homeassistant-enabled').checked,
        "url": document.getElementById('homeassistant-url').value,
        "port": document.getElementById('homeassistant-port').value,
        "token": document.getElementById('homeassistant-token').value,
        "triggerlightsonevent": document.getElementById('homeassistant-triggerlightsonevent').checked,
        "triggeravrmastervolumechangeonevent": document.getElementById('homeassistant-triggeravrmastervolumechangeonevent').checked,
        "remoteentityname": document.getElementById('homeassistant-remoteentityname').value,
        "playscriptname": document.getElementById('homeassistant-playscriptname').value,
        "pausescriptname": document.getElementById('homeassistant-pausescriptname').value,
        "stopscriptname": document.getElementById('homeassistant-stopscriptname').value
    };

    const mqttConfig = {
        "enabled": document.getElementById('mqtt-enabled').checked,
        "url": document.getElementById('mqtt-url').value,
        "username": document.getElementById('mqtt-username').value,
        "password": document.getElementById('mqtt-password').value,
        "topiclights": document.getElementById('mqtt-topiclights').value,
        "topicvolume": document.getElementById('mqtt-topicvolume').value,
        "topicbeqcurrentprofile": document.getElementById('mqtt-topicbeqcurrentprofile').value,
        "topicminidspmutestatus": document.getElementById('mqtt-topicminidspmutestatus').value,
        "topicplayingstatus": document.getElementById('mqtt-topicplayingstatus').value
    };

    const plexConfig = {
        "enabled": document.getElementById('plex-enabled').checked,
        "url": document.getElementById('plex-url').value,
        "port": document.getElementById('plex-port').value,
        "ownernamefilter": document.getElementById('plex-ownernamefilter').value,
        "deviceuuidfilter": document.getElementById('plex-deviceuuidfilter').value,
        "playermachineidentifier": document.getElementById('plex-playermachineidentifier').value,
        "playerip": document.getElementById('plex-playerip').value,
        "enabletrailersupport": document.getElementById('plex-enabletrailersupport').checked
    };
    const jellyfinConfig = {
        "enabled": document.getElementById('jellyfin-enabled').checked,
        "url": document.getElementById('jellyfin-url').value,
        "port": document.getElementById('jellyfin-port').value,
        "ownernamefilter": document.getElementById('jellyfin-ownernamefilter').value,
        "deviceuuidfilter": document.getElementById('jellyfin-deviceuuidfilter').value,
        "playermachineidentifier": document.getElementById('jellyfin-playermachineidentifier').value,
        "userID": document.getElementById('jellyfin-userID').value,
        "apiToken": document.getElementById('jellyfin-apiToken').value
    };
    const signalConfig = {
        "enabled": document.getElementById('signal-enabled').checked,
        "source": document.getElementById('signal-source').value
    };
    // Build the final config JSON
    const finalConfig = {
        "ezbeq": ezbeqConfig,
        "homeassistant": homeAssistantConfig,
        "mqtt": mqttConfig,
        "plex": plexConfig,
        "jellyfin": jellyfinConfig,
        "signal": signalConfig
    };

    return finalConfig;

}


document.addEventListener('DOMContentLoaded', function () {
    async function loadConfig() {
        const response = await fetch('/get-config');
        if (response.ok) {
            return await response.json();
        } else {
            throw new Error('Failed to fetch config');
        }
    }

    loadConfig()
        .then(config => {
            // Populate fields from config
            populateFields(config);
        })
        .catch(error => {
            console.error('Error in loading config:', error);
        });


    document.getElementById('ezbeqForm').addEventListener('submit', async function (e) {
        e.preventDefault();

        // We moved the logic to build the finalConfig object here
        const finalConfig = buildFinalConfig();
        console.log(JSON.stringify(finalConfig))
        try {
            await submitConfig(finalConfig);
            showNotification("Configuration saved successfully.");
        } catch (error) {
            showNotification("Failed to save configuration.", false);
        }
    });

});
async function submitConfig(data) {
    try {
        const response = await fetch('/save-config', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            return await response.json();
        } else {
            const text = await response.text();  // Get more information from the server
            throw new Error(`Failed to save configuration, server says: ${text}`);
        }
    } catch (error) {
        console.error("An error occurred while trying to send the request:", error);
        throw error; // Re-throw to be caught by the caller
    }
}



function showNotification(message, isSuccess = true) {
    const notification = document.getElementById("notification");
    notification.textContent = message;

    if (isSuccess) {
        notification.className = "notification-success";
    } else {
        notification.className = "notification-error";
    }

    // Remove the notification after 4 seconds
    setTimeout(() => {
        notification.textContent = "";
        notification.className = "";
    }, 4000);
}
