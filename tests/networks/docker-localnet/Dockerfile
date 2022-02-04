FROM cheqd-cli

USER root

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    supervisor \
    jq

USER cheqd

COPY --chown=cheqd:cheqd tests/networks/docker-localnet/gen-node-configs.sh ./
COPY tests/networks/docker-localnet/supervisord.conf /etc

RUN chmod +x ./gen-node-configs.sh
RUN ./gen-node-configs.sh

# Update persistent peers list for at least one node
RUN NODE_0_ID=$(cheqd-noded tendermint show-node-id); \
    PERSISTENT_PEERS="$NODE_0_ID@127.0.0.1:26656"; \
    sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/g" .cheqdnode/config/config.toml

EXPOSE 26657 26659

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
STOPSIGNAL SIGTERM
