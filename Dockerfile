### Multi-stage build
FROM jormungandrk/goa-build as build

COPY . /go/src/github.com/JormungandrK/microservice-mail
RUN go install github.com/JormungandrK/microservice-mail

### Main
FROM alpine:3.6

COPY --from=build /go/bin/microservice-mail /usr/local/bin/microservice-mail
COPY public /public
COPY config.json config.json
EXPOSE 8080

ENV API_GATEWAY_URL="http://localhost:8001"

CMD ["/usr/local/bin/microservice-mail"]
