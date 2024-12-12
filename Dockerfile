FROM golang:1.22 as build

WORKDIR /go/src/app
COPY . .
RUN go mod download
WORKDIR /go/src/app/cmd
RUN go vet -v

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM alpine:20240923

RUN apk add  supervisor
COPY docker/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY docker/watch.py /watch.py

COPY --from=build /go/bin/app /
COPY --from=build /go/src/app/web /web
EXPOSE 9999

# CMD ["/app"]
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
