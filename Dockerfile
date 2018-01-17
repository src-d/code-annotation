FROM golang:1.8-alpine3.6

# base deps
RUN apk --update upgrade && \
    apk add --no-cache make git curl ca-certificates bash \
    build-base libxml2-dev protobuf nodejs=6.10.3-r1 nodejs-npm && \
    npm install -g yarn

ADD . /go/src/code-annotation
WORKDIR /go/src/code-annotation

RUN make build && \
    make packages

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/code-annotation/build ./build
COPY --from=0 /go/src/code-annotation/bin/code-annotation .
CMD ["./code-annotation"]
