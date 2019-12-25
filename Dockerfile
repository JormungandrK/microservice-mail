### Multi-stage build
FROM golang:1.13.5-alpine3.10 as build

RUN apk --no-cache add git

COPY . /go/src/github.com/Microkubes/microservice-mail

RUN cd /go/src/github.com/Microkubes/microservice-mail && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install


### Main
FROM alpine:3.10

ENV API_GATEWAY_URL="http://localhost:8001"

COPY --from=build /go/src/github.com/Microkubes/microservice-mail/config.json /config.json
COPY --from=build /go/bin/microservice-mail /microservice-mail

COPY public /public

EXPOSE 8080

CMD ["/microservice-mail"]
