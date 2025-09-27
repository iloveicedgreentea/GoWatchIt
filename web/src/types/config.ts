export interface ConfigOption {
    key: string;
    label: string;
    description: string;
    type: 'text' | 'checkbox' | 'number' | 'select' | 'password' | 'numberArray';
    defaultValue?: string | boolean | number | number[];
    options?: { label: string; value: string }[];
    section: string;
    placeholder?: string;
  }
  
  export interface ConfigValue {
    [key: string]: {
      [key: string]: string | boolean | number | number[];
    };
  }
  
  export interface ConfigSection {
    name: string;
    enabled: boolean;
    options: ConfigOption[];
  }