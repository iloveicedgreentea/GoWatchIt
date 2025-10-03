// Validation functions for configuration form fields
// Each validator returns true if valid, or an error message string if invalid

export type ValidationResult = true | string;
export type Validator = (value: string) => ValidationResult;

/**
 * Validates IPv4 addresses (e.g., 192.168.1.1)
 */
export function validateIPv4(value: string): ValidationResult {
  if (!value) return true; // Allow empty for optional fields

  const parts = value.split('.');

  // Must have exactly 4 parts
  if (parts.length !== 4) {
    return 'Invalid IPv4 address (e.g., 192.168.1.1)';
  }

  // Each part must be a number between 0 and 255
  for (const part of parts) {
    const num = parseInt(part, 10);

    // Check if it's a valid number
    if (isNaN(num) || String(num) !== part) {
      return 'Invalid IPv4 address (e.g., 192.168.1.1)';
    }

    // Check if it's in valid range
    if (num < 0 || num > 255) {
      return 'Invalid IPv4 address - each octet must be 0-255';
    }
  }

  return true;
}

/**
 * Validates hostnames (e.g., example.com, my-server.local)
 */
export function validateHostname(value: string): ValidationResult {
  if (!value) return true; // Allow empty for optional fields

  const hostnameRegex = /^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$/i;

  if (!hostnameRegex.test(value)) {
    return 'Invalid hostname (e.g., example.com, server.local)';
  }

  return true;
}

/**
 * Validates IPv4 address OR hostname
 */
export function validateIPOrHostname(value: string): ValidationResult {
  if (!value) return true; // Allow empty for optional fields

  const ipResult = validateIPv4(value);
  if (ipResult === true) return true;

  const hostnameResult = validateHostname(value);
  if (hostnameResult === true) return true;

  return 'Invalid IP address or hostname';
}

/**
 * Validates port numbers (1-65535)
 */
export function validatePort(value: string): ValidationResult {
  if (!value) return true; // Allow empty for optional fields

  const port = parseInt(value, 10);

  if (isNaN(port)) {
    return 'Port must be a number';
  }

  if (port < 1 || port > 65535) {
    return 'Port must be between 1 and 65535';
  }

  return true;
}

/**
 * Validates URLs with http or https protocol
 */
export function validateHttpUrl(value: string): ValidationResult {
  if (!value) return true; // Allow empty for optional fields

  if (!value.startsWith('http://') && !value.startsWith('https://')) {
    return 'URL must start with http:// or https://';
  }

  try {
    new URL(value);
    return true;
  } catch {
    return 'Invalid URL format';
  }
}

/**
 * Validates that value is a non-negative integer
 */
export function validateNonNegativeInteger(value: string): ValidationResult {
  if (!value) return true; // Allow empty for optional fields

  const num = parseInt(value, 10);

  if (isNaN(num)) {
    return 'Must be a number';
  }

  if (num < 0) {
    return 'Must be a non-negative number';
  }

  if (!Number.isInteger(num)) {
    return 'Must be an integer';
  }

  return true;
}

/**
 * Validates that value is within a specified range
 */
export function validateRange(min: number, max: number): Validator {
  return (value: string): ValidationResult => {
    if (!value) return true; // Allow empty for optional fields

    const num = parseFloat(value);

    if (isNaN(num)) {
      return 'Must be a number';
    }

    if (num < min || num > max) {
      return `Must be between ${min} and ${max}`;
    }

    return true;
  };
}

/**
 * Validates using a custom regex pattern
 */
export function validatePattern(pattern: string | RegExp, errorMessage?: string): Validator {
  const regex = typeof pattern === 'string' ? new RegExp(pattern) : pattern;

  return (value: string): ValidationResult => {
    if (!value) return true; // Allow empty for optional fields

    if (!regex.test(value)) {
      return errorMessage || 'Invalid format';
    }

    return true;
  };
}

/**
 * Combines multiple validators - all must pass
 */
export function combineValidators(...validators: Validator[]): Validator {
  return (value: string): ValidationResult => {
    for (const validator of validators) {
      const result = validator(value);
      if (result !== true) {
        return result;
      }
    }
    return true;
  };
}

/**
 * Get validator based on JSON Schema pattern
 */
export function getValidatorFromPattern(pattern: string): Validator | undefined {
  // IP or hostname pattern
  if (pattern === '^[a-z0-9.-]+$|^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$') {
    return validateIPOrHostname;
  }

  // IP only pattern
  if (pattern === '^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$') {
    return validateIPv4;
  }

  // HTTP/HTTPS URL pattern
  if (pattern === '^https?://') {
    return validateHttpUrl;
  }

  // Generic pattern fallback
  return validatePattern(pattern);
}
