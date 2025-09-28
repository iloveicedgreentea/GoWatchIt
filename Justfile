help:
	just --list

# TypeSpec commands

[group('generation')]
typespec-compile:
	@echo "Compiling TypeSpec models..."
	@cd typespec && bun run compile
	@echo "TypeSpec compilation complete"

[group('generation')]
typespec-generate-go:
	@echo "Generating Go structs from OpenAPI..."
	@cd typespec && bun run generate-go
	@echo "Go structs generated at typespec/generated/go/models.go"

[group('generation')]
typespec-generate-ts:
	@echo "Generating TypeScript types from OpenAPI..."
	@cd typespec && bun run generate-ts
	@echo "TypeScript types generated at typespec/generated/typescript/types.ts"

[group('generation')]
typespec-copy-schemas:
	@echo "Copying JSON schemas to web app..."
	@cp typespec/generated/json-schema/*.json web/src/schemas/
	@cp typespec/generated/typescript/types.ts web/src/types/typespec-types.ts
	@echo "Schemas copied to web app"

[group('generation')]
[doc('Use this command to generate all TypeSpec artifacts')]
typespec-generate: typespec-compile typespec-generate-go typespec-generate-ts typespec-copy-schemas
	@echo "All TypeSpec artifacts generated"

[group('generation')]
typespec-watch:
	@echo "Watching TypeSpec models for changes..."
	@cd typespec && bun run watch

[group('build')]
build: typespec-generate web-build
	@echo "Building Go application..."
	@cd cmd/gowatchit && go build -o ../../build/gowatchit
	@echo "Build complete"

[group('build')]
build-go: typespec-generate
	@echo "Building Go application only..."
	@cd cmd/gowatchit && go build -o ../../build/gowatchit
	@echo "Go build complete"

[group('test')]
test:
	./test_web.sh

[group('docker')]
docker-build:
	docker buildx build --platform linux/arm64 --load --tag gowatchit-local . -f services/gowatchit/Dockerfile

[group('docker')]
docker-push:
	docker buildx build --push --platform linux/amd64 --tag ghcr.io/iloveicedgreentea/gowatchit:test .

[group('docker')]
docker-run: docker-build
	LOG_FILE=true LOG_LEVEL=debug docker-compose -f services/gowatchit/docker-compose.yml up

[group('run')]
run:
	LOG_ENV=local LOG_FILE=false LOG_LEVEL=debug go run ./services/gowatchit

# Web development commands

[group('setup')]
web-install:
	@echo "Installing web dependencies with bun..."
	@cd web && bun install
	@echo "Web dependencies installed"

[group('web')]
web-dev:
	@echo "Starting web development server..."
	@cd web && bun run dev

[group('build')]
web-build:
	@echo "Building web application..."
	@cd web && bun run build
	@echo "Web build complete"

[group('web')]
web-preview:
	@echo "Starting preview server..."
	@cd web && bun run preview

[group('lint')]
web-lint:
	@echo "Running web linting..."
	@cd web && bun run lint

[group('clean')]
web-clean:
	@echo "Cleaning web build artifacts..."
	@cd web && rm -rf dist node_modules/.vite
	@echo "Web artifacts cleaned"

# Full development setup
[group('setup')]
dev-setup: web-install typespec-generate
	@echo "Development environment ready"

# Alias for backward compatibility
[group('web')]
run-ui: web-dev

[group('test')]
live-test:
	./test_web.sh