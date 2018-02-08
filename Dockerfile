FROM alpine:3.6
ADD ./build/bin /bin
ENTRYPOINT ["/bin/code-annotation"]
