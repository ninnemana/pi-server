# build stage
FROM golang:alpine AS build-env

RUN mkdir -p /etc/ssl
RUN touch /etc/ssl/cert.key
RUN touch /etc/ssl/private.key

ADD . $GOPATH/src/github.com/ninnemana/pi-server
WORKDIR $GOPATH/src/github.com/ninnemana/pi-server

RUN go build -o /server

# final stage
FROM alpine

RUN mkdir -p /etc/ssl

WORKDIR /app
COPY --from=build-env /server /app/
COPY --from=build-env /etc/ssl/cert.key /etc/ssl/cert.key
COPY --from=build-env /etc/ssl/private.key /etc/ssl/private.key

EXPOSE 80
EXPOSE 443

ENTRYPOINT ["/app/server", "-cert=/etc/ssl/cert.key", "-key=/etc/ssl/private.key"]
