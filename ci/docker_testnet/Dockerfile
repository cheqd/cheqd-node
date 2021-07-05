FROM verim-node

USER root

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    supervisor \
    jq

USER verim

COPY ci/docker-testnet/gen_node_configs.sh .
COPY ci/docker-testnet/supervisord.conf /etc

EXPOSE 26657 26659

RUN ./gen_node_configs.sh

# Update persistent peers list for at least one node
RUN NODE_0_ID=$(verim-noded tendermint show-node-id --home node_configs/node0); \
    NODE_1_ID=$(verim-noded tendermint show-node-id --home node_configs/node1); \
    PERSISTENT_PEERS="$NODE_0_ID@127.0.0.1:26656,$NODE_1_ID@127.0.0.1:26658"; \
    sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/g" node_configs/node0/config/config.toml

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
STOPSIGNAL SIGTERM
