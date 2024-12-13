export interface ConfigOption {
    key: string;
    label: string;
    description: string;
    type: 'text' | 'checkbox' | 'number' | 'select' | 'password';
    defaultValue?: string | boolean | number;
    options?: { label: string; value: string }[];
    section: string;
    placeholder?: string;
}

export interface ConfigValue {
    [key: string]: {
        [key: string]: string | boolean | number;
    };
}

export interface ConfigSection {
    name: string;
    enabled: boolean;
    options: ConfigOption[];
}


export const CONFIG_SCHEMA: ConfigSection[] = [
    {
        name: 'ezbeq',
        enabled: false,
        options: [
            {
                key: 'adjustmastervolumewithprofile',
                label: 'Adjust Master Volume With Profile',
                description: 'Most BEQ profiles have a Master Volume adjustment. Most lower the MV but some raise it. Recommended on because the BEQ is created with MV in mind.',
                type: 'checkbox',
                defaultValue: true,
                section: 'ezbeq'
            },
            {
                key: 'enabled',
                label: 'Enabled',
                description: 'Use EzBEQ',
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
                key: 'avrbrand',
                label: 'Source',
                description: 'Supported AVR brands - Currently only supports "denon" for all Denon and Marantz',
                type: 'select',
                options: [{ label: 'Denon', value: 'denon' }],
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
            }
        ]
    }
];