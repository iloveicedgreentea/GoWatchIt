# Backend build stage
FROM golang:1.25 AS backend-build
WORKDIR /go/src/app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build from the service directory
WORKDIR /go/src/app/services/gowatchit
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app main.go

# Frontend build stage
FROM oven/bun:1 AS frontend-build
WORKDIR /usr/src/app
COPY web/package.json web/bun.lockb ./
RUN bun install --frozen-lockfile
COPY web/ ./
RUN bun run build

# Final stage - using golang alpine for CGO compatibility
FROM debian:stable-slim

# Install required packages
RUN apt-get update && apt-get install -y --no-install-recommends \
    tzdata \
    nginx \
    sqlite3 \
    libsqlite3-dev \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create necessary directories
RUN mkdir -p /data \
    /run/nginx \
    /var/lib/nginx/tmp \
    /var/lib/nginx/logs \
    /var/lib/nginx/tmp/client_body \
    /var/lib/nginx/tmp/proxy \
    /var/lib/nginx/tmp/fastcgi \
    /var/lib/nginx/tmp/uwsgi \
    /var/lib/nginx/tmp/scgi

# Copy backend
COPY --from=backend-build /go/bin/app /gowatchit
RUN chmod +x /gowatchit

# Copy frontend build files
COPY --from=frontend-build /usr/src/app/dist /var/www/html/

# Copy config files
COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 9999
EXPOSE 3000

ENV TZ=America/New_York
ENV GIN_MODE=release
ENV LOG_FILE=true
ENV BASE_DIR=/data

CMD ["/entrypoint.sh"]