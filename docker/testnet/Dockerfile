FROM cheqd-node

USER root

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    supervisor \
    jq

USER cheqd

COPY --chown=cheqd:cheqd docker/testnet/gen_node_configs.sh .
COPY docker/testnet/supervisord.conf /etc

EXPOSE 26657 26659

RUN ./gen_node_configs.sh

# Update persistent peers list for at least one node
RUN NODE_0_ID=$(cheqd-noded tendermint show-node-id); \
    PERSISTENT_PEERS="$NODE_0_ID@127.0.0.1:26656"; \
    sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/g" .cheqdnode/config/config.toml

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
STOPSIGNAL SIGTERM
