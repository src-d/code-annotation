FROM debian:buster-slim
ADD ./build/bin /bin

RUN apt-get update && \
    apt-get install -y --no-install-recommends --no-install-suggests \
    ca-certificates \
    && apt-get clean
RUN mkdir /var/code-annotation

ENTRYPOINT ["/bin/server"]
