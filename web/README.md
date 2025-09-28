# GoWatchIt Web Interface

React + TypeScript + Vite web interface for GoWatchIt configuration management.

## Prerequisites

- [Bun](https://bun.sh/) - Fast JavaScript runtime and package manager
- Node.js 18+ (for compatibility)

## Installation

```bash
# Install dependencies
bun install
```

## Development

```bash
# Start development server
bun run dev

# Build for production
bun run build

# Preview production build
bun run preview

# Run linting
bun run lint
```

## TypeSpec Integration

This web interface uses automatically generated types and schemas from TypeSpec models. The configuration UI is dynamically generated from JSON Schema files.

### Generated Files

- `src/types/typespec-types.ts` - TypeScript types generated from OpenAPI
- `src/schemas/*.json` - JSON Schema files for UI generation

### Development Workflow

1. Make changes to TypeSpec models in `../typespec/config/config.tsp`
2. Run `just typespec-generate` from project root to regenerate all artifacts
3. The web interface automatically uses the updated schemas

### Configuration Schema Generation

The configuration UI is automatically generated from TypeSpec JSON Schema files:

- **Field Types**: Automatically mapped from JSON Schema types
- **Validation**: Patterns, min/max values from schema
- **Descriptions**: Documentation from TypeSpec `@doc` decorators
- **Examples**: Placeholder values from TypeSpec `@example` decorators
- **Default Values**: From TypeSpec model defaults

### Manual Schema Updates

If you need to manually update schemas (not recommended):

```bash
# Copy latest schemas from TypeSpec
just typespec-copy-schemas
```

## Project Structure

```
web/
├── src/
│   ├── components/
│   │   ├── config/         # Configuration form components
│   │   ├── layout/         # Layout components
│   │   └── ui/            # Reusable UI components
│   ├── lib/
│   │   └── schema-loader.ts # JSON Schema to UI transformation
│   ├── pages/
│   │   └── ConfigurationPage.tsx # Main config page
│   ├── schemas/           # Generated JSON Schema files
│   └── types/
│       ├── config.ts      # Configuration types
│       └── typespec-types.ts # Generated TypeScript types
├── package.json
└── vite.config.ts
```

## Configuration Components

### ConfigurationPage

Main page that loads and displays configuration sections dynamically from TypeSpec schemas.

### ConfigSection

Renders a configuration section with enable/disable toggle and fields.

### Schema Loader

Transforms JSON Schema files into UI-friendly configuration objects with proper field types, validation, and metadata.

## Type Safety

All configuration types are generated from TypeSpec models ensuring:

- **Compile-time type checking** between frontend and backend
- **Automatic UI updates** when schema changes
- **Single source of truth** for all configuration models

## API Integration

The configuration interface communicates with the Go backend using TypeSpec-generated types:

```typescript
// All types are automatically generated and type-safe
const config: AppConfig = await fetch("/api/config").then((r) => r.json());
```

## Troubleshooting

### Schema Changes Not Reflecting

Run `just typespec-generate` from the project root to regenerate all schemas and types.

### TypeScript Errors

Ensure TypeSpec generation completed successfully and all generated files are present in `src/schemas/` and `src/types/`.

### Build Issues

1. Clear node_modules: `rm -rf node_modules && bun install`
2. Clear Vite cache: `rm -rf node_modules/.vite`
3. Regenerate schemas: `just typespec-generate`
