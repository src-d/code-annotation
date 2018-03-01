FROM alpine:3.7
ADD ./build/bin /bin

RUN apk --update upgrade && \
    apk add --no-cache ca-certificates
RUN mkdir /var/code-annotation

ENTRYPOINT ["/bin/server"]
