FROM golang:1.23 as build

WORKDIR /go/src/app
COPY . .
RUN go mod download
WORKDIR /go/src/app/cmd/gowatchit
RUN go vet -v

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM alpine:20240923

RUN apk add --no-cache tzdata supervisor
COPY docker/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY docker/watch.py /watch.py

COPY --from=build /go/bin/app /
COPY ./web /web
EXPOSE 9999

ENV TZ=America/New_York
ENV GIN_MODE=release
ENV LOG_TO_FILE=true

# CMD ["/app"]
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
