FROM golang:1.19 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM alpine:latest

COPY --from=build /go/bin/app /
CMD ["/app"]