### Multi-stage build
FROM golang:1.8.3-alpine3.6 as build

RUN apk add --no-cache git

COPY . /go/src/github.com/JormungandrK/microservice-email

WORKDIR /go/src/github.com/JormungandrK/microservice-email

RUN go get -u -v github.com/golang/dep/cmd/dep
RUN dep ensure -v
RUN go install github.com/JormungandrK/microservice-email


### Main
FROM alpine:3.6

COPY --from=build /go/bin/microservice-email /usr/local/bin/microservice-email
EXPOSE 8080

ENV API_GATEWAY_URL="http://localhost:8001"

CMD ["/usr/local/bin/microservice-email"]
