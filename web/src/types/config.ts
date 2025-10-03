import type { components } from './typespec-types';
import type { Validator } from '../lib/validators';

export interface ConfigOption {
    key: string;
    label: string;
    description: string;
    type: 'text' | 'checkbox' | 'number' | 'select' | 'password' | 'numberArray' | 'stringArray';
    defaultValue?: string | boolean | number | number[] | string[];
    options?: { label: string; value: string }[];
    section: string;
    placeholder?: string;
    validator?: Validator;
  }

  export interface ConfigSection {
    name: string;
    enabled: boolean;
    options: ConfigOption[];
  }

  // Use TypeSpec-generated types for actual configuration data
  export type AppConfig = components['schemas']['AppConfig'];
  export type EZBEQConfig = components['schemas']['EZBEQConfig'];
  export type HomeAssistantConfig = components['schemas']['HomeAssistantConfig'];
  export type PlayerConfig = components['schemas']['PlayerConfig'];
  export type MainConfig = components['schemas']['MainConfig'];
  export type HDMISyncConfig = components['schemas']['HDMISyncConfig'];
  export type Player = components['schemas']['Player'];

  // Configuration value type that matches the API structure
  export interface ConfigValue {
    ezbeq?: EZBEQConfig;
    homeAssistant?: HomeAssistantConfig;
    player?: PlayerConfig;
    main?: MainConfig;
    hdmiSync?: HDMISyncConfig;
    // Allow string indexing for dynamic access
    [key: string]: any;
  }