FROM alpine:3.6
ADD ./build/bin /bin

RUN apk --update upgrade && \
    apk add --no-cache ca-certificates

ENTRYPOINT ["/bin/server"]
