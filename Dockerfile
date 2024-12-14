# Backend build stage
FROM golang:1.23 AS backend-build
WORKDIR /go/src/app
COPY . .
RUN go mod download
WORKDIR /go/src/app/cmd/gowatchit

RUN go build -o /go/bin/app

# Frontend build stage
FROM oven/bun:1 AS frontend-build
WORKDIR /usr/src/app
COPY web/package.json web/bun.lockb ./
RUN bun install --frozen-lockfile
COPY web/ ./
RUN bun run build

# Final stage
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y \
    tzdata \
    supervisor \
    nginx \
    sqlite3 \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*      

# Create necessary directories
RUN mkdir -p /data \
    /run/supervisor \
    /var/log/supervisor \
    /run/nginx \
    /var/lib/nginx \
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
COPY docker/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 9999
EXPOSE 3000

ENV TZ=America/New_York
ENV GIN_MODE=release
ENV LOG_TO_FILE=true

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]