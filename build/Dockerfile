FROM alpine:3.8

RUN apk upgrade --update --no-cache

USER nobody

ADD build/_output/bin/node-network-operator /usr/local/bin/node-network-operator
