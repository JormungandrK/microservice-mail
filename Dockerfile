### Multi-stage build
FROM golang:1.10-alpine3.7 as build

RUN apk --no-cache add git

RUN go get -u -v gopkg.in/gomail.v2 && \
    go get -u -v github.com/Microkubes/microservice-tools/...

COPY . /go/src/github.com/Microkubes/microservice-mail

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install github.com/Microkubes/microservice-mail


### Main
FROM scratch

ENV API_GATEWAY_URL="http://localhost:8001"

COPY --from=build /go/src/github.com/Microkubes/microservice-mail/config.json /config.json
COPY --from=build /go/bin/microservice-mail /microservice-mail

COPY public /public

EXPOSE 8080

CMD ["/microservice-mail"]
