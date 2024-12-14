import { ConfigSection } from './config';

// This object stores all config options
export const CONFIG_SCHEMA: ConfigSection[] = [
    {
        name: 'ezbeq',
        enabled: false,
        options: [
            {
                key: 'enabled',
                label: 'Enabled',
                description: 'Use EzBEQ',
                type: 'checkbox',
                defaultValue: false,
                section: 'ezbeq'
            },
            {
                key: 'adjustmastervolumewithprofile',
                label: 'Adjust Master Volume With Profile',
                description: 'Most BEQ profiles have a Master Volume adjustment. Most lower the MV but some raise it. I recommend this on because the BEQ is created with MV in mind.',
                type: 'checkbox',
                defaultValue: true,
                section: 'ezbeq'
            },
            {
                key: 'dryrun',
                label: 'Dryrun',
                description: 'Keep EzBEQ in dryrun mode',
                type: 'checkbox',
                defaultValue: false,
                section: 'ezbeq'
            },
            {
                key: 'enabletvbeq',
                label: 'TV BEQ',
                description: 'Check TV media types for BEQ',
                type: 'checkbox',
                defaultValue: false,
                section: 'ezbeq'
            },
            {
                key: 'notifyonload',
                label: 'Notify On Load',
                description: 'Notify you on BEQ Load/Unload/Issues',
                type: 'checkbox',
                defaultValue: false,
                section: 'ezbeq'
            },
            {
                key: 'port',
                label: 'EzBEQ Port',
                description: 'EzBEQ Port',
                type: 'text',
                defaultValue: '8080',
                section: 'ezbeq'
            },
            {
                key: 'preferredauthor',
                label: 'BEQ Preferred Author',
                description: 'A whitelist of authors. Comma separated. Leave blank to use the first one returned.',
                type: 'text',
                section: 'ezbeq'
            },
            {
                key: 'slots',
                label: 'MiniDSP Slots',
                description: 'Which slot(s) to load into',
                type: 'numberArray',
                defaultValue: [1],
                section: 'ezbeq'
            },
            {
                key: 'stopplexifmismatch',
                label: 'Stop Plex On Mismatch',
                description: 'Send a Stop to Plex if its transcoding incorrectly',
                type: 'checkbox',
                defaultValue: false,
                section: 'ezbeq'
            },
            {
                key: 'url',
                label: 'EzBEQ URL',
                description: 'EzBEQ URL - Must have http://',
                type: 'text',
                placeholder: 'http://x.x.x.x',
                section: 'ezbeq'
            },
            {
                key: 'useavrcodecsearch',
                label: 'Use AVR For Codec Lookup',
                description: 'Use a supported AVR to get the codec instead of Plex metadata. Could be more accurate.',
                type: 'checkbox',
                defaultValue: false,
                section: 'ezbeq'
            },
            {
                key: 'avrbrand',
                label: 'Source',
                description: 'Supported AVR brands - Currently only supports "denon" for all Denon and Marantz',
                type: 'select',
                options: [{ label: 'Denon', value: 'denon' }],
                section: 'ezbeq'
            },
            {
                key: 'avrip',
                label: 'AVR IP address',
                description: 'IP Address for your AVR - "x.x.x.x"',
                type: 'text',
                placeholder: 'x.x.x.x',
                section: 'ezbeq'
            }
        ]
    },
    {
        name: 'homeassistant',
        enabled: false,
        options: [
            {
                key: 'enabled',
                label: 'Enabled',
                description: 'Use Home Assistant features',
                type: 'checkbox',
                defaultValue: false,
                section: 'homeassistant'
            },
            {
                key: 'url',
                label: 'Home Assistant URL',
                description: 'URL - Must have http://',
                type: 'text',
                placeholder: 'http://x.x.x.x',
                section: 'homeassistant'
            },
            {
                key: 'port',
                label: 'Home Assistant port',
                description: 'port - "8123"',
                type: 'text',
                defaultValue: '8123',
                section: 'homeassistant'
            },
            {
                key: 'token',
                label: 'Home Assistant token',
                description: 'token - "ey.xyz" get a token from your user profile',
                type: 'text',
                section: 'homeassistant'
            },
            {
                key: 'notifyendpointname',
                label: 'Notify Endpoint',
                description: 'Name of the Home Assistant notify endpoint (like your phone)',
                type: 'text',
                section: 'homeassistant'
            }
        ]
    },
    {
        name: 'plex',
        enabled: false,
        options: [
            {
                key: 'enabled',
                label: 'Enabled',
                description: 'Use plex',
                type: 'checkbox',
                defaultValue: false,
                section: 'plex'
            },
            {
                key: 'url',
                label: 'Plex URL',
                description: 'IP or domain of your plex server',
                type: 'text',
                placeholder: 'x.x.x.x',
                section: 'plex'
            },
            {
                key: 'port',
                label: 'Plex port',
                description: 'port - "32400"',
                type: 'text',
                defaultValue: '32400',
                section: 'plex'
            },
            {
                key: 'ownernamefilter',
                label: 'Owner Name Filter',
                description: 'Your primary account, will filter webhooks so others don\'t trigger. Leave blank if you don\'t want to filter on accounts',
                type: 'text',
                section: 'plex'
            },
            {
                key: 'deviceuuidfilter',
                label: 'Device UUID Filter',
                description: 'The client identifier of the device to filter webhooks for. You should set this to your theater player UUID, or a comma separated string (filter1,filter2) if you have multiple devices.',
                type: 'text',
                section: 'plex'
            },
            {
                key: 'enabletrailersupport',
                label: 'Enable Trailer Support',
                description: 'If you enable trailers before movies, it can process it like turn off lights. No beq obviously.',
                type: 'checkbox',
                defaultValue: false,
                section: 'plex'
            }
        ]
    },
    {
        name: 'Generic Media Player',
        enabled: false,
        options: [
            {
                key: 'enabled',
                label: 'Enabled',
                description: 'Use jellyfin',
                type: 'checkbox',
                defaultValue: false,
                section: 'jellyfin'
            },
            {
                key: 'skiptmdb',
                label: 'Allow Skipping TMDB Check',
                description: 'Jellyfin doesn\'t have robust metadata services like Plex so some items can just be entirely missing TMDB. If checked, BEQ matching will allow skipping TMDB if its not found from Jellyfin.',
                type: 'checkbox',
                defaultValue: false,
                section: 'jellyfin'
            },
            {
                key: 'url',
                label: 'Jellyfin URL',
                description: 'IP or domain of your server',
                type: 'text',
                placeholder: 'x.x.x.x',
                section: 'jellyfin'
            },
            {
                key: 'port',
                label: 'jellyfin port',
                description: 'port - "8096"',
                type: 'text',
                defaultValue: '8096',
                section: 'jellyfin'
            },
            {
                key: 'ownernamefilter',
                label: 'Owner Name Filter',
                description: 'Your main owner account, will filter webhooks so others don\'t trigger. Leave blank if you don\'t want to filter on accounts',
                type: 'text',
                section: 'jellyfin'
            },
            {
                key: 'deviceuuidfilter',
                label: 'Device UUID Filter',
                description: 'The client identifier of the device to filter webhooks for.',
                type: 'text',
                section: 'jellyfin'
            },
            {
                key: 'userid',
                label: 'User ID',
                description: 'The GUID of your user. Get this from your jellyfin profile page -> check the url',
                type: 'text',
                section: 'jellyfin'
            },
            {
                key: 'apitoken',
                label: 'API Token',
                description: 'The Jellyfin API token from Dashboard -> Admin -> API Keys',
                type: 'text',
                section: 'jellyfin'
            }
        ]
    },
    {
        name: 'signal',
        enabled: false,
        options: [
            {
                key: 'enabled',
                label: 'Enable HDMI Signal Sync',
                description: 'if you want to pause plex until hdmi sync is done',
                type: 'checkbox',
                defaultValue: false,
                section: 'signal'
            },
            {
                key: 'source',
                label: 'Source',
                description: 'What source to use for HDMI sync info. Either wait X seconds or attributes from a MadVR Envy',
                type: 'select',
                options: [
                    { label: 'Madvr Envy', value: 'envy' },
                    { label: 'Time', value: 'time' }
                ],
                section: 'signal'
            },
            {
                key: 'time',
                label: 'Sync Time',
                description: 'Time to wait for HDMI sync to finish. In seconds e.g "15". Leave blank if using Madvr Envy',
                type: 'text',
                section: 'signal'
            },
            {
                key: 'envy',
                label: 'Madvr Envy Name',
                description: 'Entity name of the madvr envy like "envy". Leave blank if using time. NOTE: this doesnt\' work unless your sync time is really long. I recommend using time.',
                type: 'text',
                section: 'signal'
            },
            {
                key: 'playermachineidentifier',
                label: 'Player Machine Identifier',
                description: 'Optional if not using hdmi sync - get this from "http://(player ip):32500/resources". Required for pausing/playing signal.',
                type: 'text',
                section: 'signal'
            },
            {
                key: 'playerip',
                label: 'Player IP',
                description: 'Optional if not using hdmi sync - IP of your CLIENT device (like a shield)',
                type: 'text',
                placeholder: 'x.x.x.x',
                section: 'signal'
            }
        ]
    }
];