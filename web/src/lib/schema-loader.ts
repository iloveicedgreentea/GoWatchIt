// Import JSON Schema files
import ezbeqSchema from '../schemas/EZBEQConfig.json';
import homeAssistantSchema from '../schemas/HomeAssistantConfig.json';
import playerSchema from '../schemas/PlayerConfig.json';
import mainSchema from '../schemas/MainConfig.json';
import hdmiSyncSchema from '../schemas/HDMISyncConfig.json';
import { ConfigOption, ConfigSection } from '../types/config';

// TODO: add input validation based on schema (e.g., pattern for IP address)

// JSON Schema interface
interface JSONSchemaProperty {
  type?: string;
  description?: string;
  default?: any;
  examples?: any[];
  format?: string;
  pattern?: string;
  enum?: string[];
  anyOf?: Array<{ const: string }>;
  minItems?: number;
  maxItems?: number;
  minimum?: number;
  maximum?: number;
  items?: { type: string };
  $ref?: string;
}

interface JSONSchema {
  type: string;
  description?: string;
  properties: Record<string, JSONSchemaProperty>;
  required?: string[];
  $defs?: Record<string, any>;
}

// Map JSON Schema types to UI field types
function getUIFieldType(property: JSONSchemaProperty, schema: JSONSchema): string {
  // Handle password format
  if (property.format === 'password') return 'password';

  // Handle enum/select fields
  if (property.enum) return 'select';
  if (property.anyOf?.some(item => item.const)) return 'select';

  // Handle $ref enums (like Player type)
  if (property.$ref && schema.$defs) {
    const refKey = property.$ref.replace('#/$defs/', '');
    const refSchema = schema.$defs[refKey];
    if (refSchema?.enum) return 'select';
  }

  // Handle basic types
  if (property.type === 'boolean') return 'checkbox';
  if (property.type === 'array') return 'numberArray'; // Assuming number arrays for now
  if (property.type === 'integer' || property.type === 'number') return 'number';

  return 'text';
}

// Convert camelCase to human-readable label
function toLabel(str: string): string {
  return str
    .replace(/([A-Z])/g, ' $1')
    .replace(/^./, (str) => str.toUpperCase())
    .trim();
}

// Get select options from schema property
function getSelectOptions(property: JSONSchemaProperty, schema: JSONSchema): Array<{ label: string; value: string }> | undefined {
  // Direct enum
  if (property.enum) {
    return property.enum.map(value => ({ label: value, value }));
  }

  // anyOf const pattern
  if (property.anyOf) {
    const options = property.anyOf
      .filter(item => item.const)
      .map(item => item.const);
    if (options.length > 0) {
      return options.map(value => ({ label: value, value }));
    }
  }

  // $ref enum
  if (property.$ref && schema.$defs) {
    const refKey = property.$ref.replace('#/$defs/', '');
    const refSchema = schema.$defs[refKey];
    if (refSchema?.enum) {
      return refSchema.enum.map((value: string) => ({ label: value, value }));
    }
  }

  return undefined;
}

// Convert JSON Schema property to ConfigOption
function schemaPropertyToConfigOption(
  key: string,
  property: JSONSchemaProperty,
  schema: JSONSchema,
  section: string
): ConfigOption {
  const fieldType = getUIFieldType(property, schema);
  const options = getSelectOptions(property, schema);

  return {
    key,
    label: toLabel(key),
    description: property.description || '',
    type: fieldType as ConfigOption['type'],
    defaultValue: property.default,
    placeholder: property.examples?.[0] ? String(property.examples[0]) : undefined,
    options,
    section,
  };
}

// Convert JSON Schema to ConfigSection
function jsonSchemaToConfigSection(schema: JSONSchema, sectionName: string): ConfigSection {
  const options: ConfigOption[] = [];

  if (schema.properties) {
    for (const [key, property] of Object.entries(schema.properties)) {
      // Skip internal ID field
      if (key === 'id') continue;

      const configOption = schemaPropertyToConfigOption(key, property, schema, sectionName);
      options.push(configOption);
    }
  }

  return {
    name: sectionName,
    enabled: false, // Will be determined dynamically based on data
    options,
  };
}

// Generate all configuration sections from schemas
export function generateConfigSchema(): ConfigSection[] {
  const schemas = [
    { schema: ezbeqSchema as JSONSchema, name: 'ezbeq' },
    { schema: homeAssistantSchema as JSONSchema, name: 'homeassistant' },
    { schema: playerSchema as JSONSchema, name: 'player' },
    { schema: mainSchema as JSONSchema, name: 'main' },
    { schema: hdmiSyncSchema as JSONSchema, name: 'hdmisync' },
  ];

  return schemas.map(({ schema, name }) => jsonSchemaToConfigSection(schema, name));
}

// Get individual schema by name
export function getSchemaByName(name: string): JSONSchema | undefined {
  const schemaMap: Record<string, JSONSchema> = {
    ezbeq: ezbeqSchema as JSONSchema,
    homeassistant: homeAssistantSchema as JSONSchema,
    player: playerSchema as JSONSchema,
    main: mainSchema as JSONSchema,
    hdmisync: hdmiSyncSchema as JSONSchema,
  };

  return schemaMap[name];
}