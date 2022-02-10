FROM cheqd-node

USER root

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    supervisor \
    jq \
    nano

USER cheqd

COPY --chown=cheqd:cheqd tests/networks/tools/gen-network-config-base.sh ./
COPY --chown=cheqd:cheqd tests/networks/docker-localnet/gen-network-config.sh ./
COPY tests/networks/docker-localnet/supervisord.conf /etc

RUN chmod +x ./gen-network-config.sh; \
    ./gen-network-config.sh

EXPOSE 1317 26656 26657 9090

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
STOPSIGNAL SIGTERM
