# build stage
FROM golang:alpine AS build-env

ADD . $GOPATH/src/github.com/ninnemana/pi-server
ADD $CERT_FILE .
ADD $KEY_FILE .
RUN ls -l
WORKDIR $GOPATH/src/github.com/ninnemana/pi-server

RUN go build -o /server

# final stage
FROM alpine

WORKDIR /app
COPY --from=build-env /server /app/

EXPOSE 80
EXPOSE 443

ENTRYPOINT ./server
