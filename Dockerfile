# build stage
FROM golang:alpine AS build-env

ADD . $GOPATH/src/github.com/ninnemana/pi-server
WORKDIR $GOPATH/src/github.com/ninnemana/pi-server

RUN go build -o /server

# final stage
FROM alpine

WORKDIR /app
COPY --from=build-env /server /app/

EXPOSE 8080

ENTRYPOINT ["/app/server"]
