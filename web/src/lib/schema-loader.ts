// Import JSON Schema files
import ezbeqSchema from '../schemas/EZBEQConfig.json';
import homeAssistantSchema from '../schemas/HomeAssistantConfig.json';
import playerSchema from '../schemas/PlayerConfig.json';
import hdmiSyncSchema from '../schemas/HDMISyncConfig.json';
import { ConfigOption, ConfigSection } from '../types/config';
import { validatePort, validateRange, getValidatorFromPattern, type Validator } from './validators';


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
  if (property.type === 'array') {
    // Check if it's a string array or number array
    if (property.items?.type === 'string') return 'stringArray';
    return 'numberArray';
  }
  if (property.type === 'integer' || property.type === 'number') return 'number';

  return 'text';
}

// Convert camelCase to human-readable label
// Handles both regular camelCase (myVariable -> My Variable) and acronyms (avrURL -> Avr URL, useAVRCodec -> Use AVR Codec)
function toLabel(str: string): string {
  return str
    // Insert space before uppercase letter that follows a lowercase letter (camelCase boundaries)
    .replace(/([a-z])([A-Z])/g, '$1 $2')
    // Insert space before uppercase letter followed by lowercase (preserves acronyms as groups)
    .replace(/([A-Z])([A-Z][a-z])/g, '$1 $2')
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

// Get validator for a JSON Schema property
function getValidator(key: string, property: JSONSchemaProperty): Validator | undefined {
  // Port validation for fields ending in "Port" or named "listenPort"
  if (key.endsWith('Port') || key === 'listenPort') {
    return validatePort;
  }

  // Range validation for min/max values
  if (property.minimum !== undefined || property.maximum !== undefined) {
    const min = property.minimum ?? -Infinity;
    const max = property.maximum ?? Infinity;
    return validateRange(min, max);
  }

  // Pattern validation from JSON Schema
  if (property.pattern) {
    return getValidatorFromPattern(property.pattern);
  }

  return undefined;
}

// Fetch BEQ authors from API
async function fetchBEQAuthors(): Promise<Array<{ label: string; value: string }>> {
  try {
    // Use window.location to construct the full URL to the backend
    const apiUrl = `${window.location.protocol}//${window.location.hostname}:9999/api/v1/config/authors`;
    const response = await fetch(apiUrl);
    if (!response.ok) {
      console.error('Failed to fetch BEQ authors:', response.status, response.statusText);
      return [];
    }

    // Check if response is JSON
    const contentType = response.headers.get('content-type');
    if (!contentType || !contentType.includes('application/json')) {
      console.error('Expected JSON response but got:', contentType);
      return [];
    }

    const authors: string[] = await response.json();
    if (!Array.isArray(authors)) {
      console.error('Expected array of authors but got:', typeof authors);
      return [];
    }

    return authors.map(author => ({ label: author, value: author }));
  } catch (error) {
    console.error('Error fetching BEQ authors:', error);
    return [];
  }
}

// Cache for authors to avoid repeated fetches
let authorsCache: Array<{ label: string; value: string }> | null = null;

// Get options for string array fields
async function getStringArrayOptions(key: string): Promise<Array<{ label: string; value: string }> | undefined> {
  // Handle preferredAuthors field specifically
  if (key === 'preferredAuthors') {
    if (!authorsCache) {
      authorsCache = await fetchBEQAuthors();
    }
    // Return undefined if no authors (allows all authors)
    return authorsCache.length > 0 ? authorsCache : undefined;
  }
  return undefined;
}

// Convert JSON Schema property to ConfigOption
async function schemaPropertyToConfigOption(
  key: string,
  property: JSONSchemaProperty,
  schema: JSONSchema,
  section: string
): Promise<ConfigOption> {
  const fieldType = getUIFieldType(property, schema);
  let options = getSelectOptions(property, schema);

  // For string arrays, get specific options if available
  if (fieldType === 'stringArray') {
    options = await getStringArrayOptions(key);
  }

  const validator = getValidator(key, property);

  return {
    key,
    label: toLabel(key),
    description: property.description || '',
    type: fieldType as ConfigOption['type'],
    defaultValue: property.default,
    placeholder: property.examples?.[0] ? String(property.examples[0]) : undefined,
    options,
    section,
    validator,
  };
}

// Convert JSON Schema to ConfigSection
async function jsonSchemaToConfigSection(schema: JSONSchema, sectionName: string): Promise<ConfigSection> {
  const options: ConfigOption[] = [];

  if (schema.properties) {
    for (const [key, property] of Object.entries(schema.properties)) {
      // Skip internal ID field
      if (key === 'id') continue;

      const configOption = await schemaPropertyToConfigOption(key, property, schema, sectionName);
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
export async function generateConfigSchema(): Promise<ConfigSection[]> {
  const schemas = [
    { schema: ezbeqSchema as JSONSchema, name: 'ezbeq' },
    { schema: homeAssistantSchema as JSONSchema, name: 'homeassistant' },
    { schema: playerSchema as JSONSchema, name: 'player' },
    { schema: hdmiSyncSchema as JSONSchema, name: 'hdmisync' },
  ];

  return Promise.all(schemas.map(({ schema, name }) => jsonSchemaToConfigSection(schema, name)));
}

// Get individual schema by name
export function getSchemaByName(name: string): JSONSchema | undefined {
  const schemaMap: Record<string, JSONSchema> = {
    ezbeq: ezbeqSchema as JSONSchema,
    homeassistant: homeAssistantSchema as JSONSchema,
    player: playerSchema as JSONSchema,
    hdmisync: hdmiSyncSchema as JSONSchema,
  };

  return schemaMap[name];
}