FROM golang:1.21 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
WORKDIR /go/src/app/cmd
RUN go vet -v

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian11


COPY --from=build /go/bin/app /
CMD ["/app"]
