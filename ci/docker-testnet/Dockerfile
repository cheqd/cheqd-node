FROM verim-node

USER root

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    # supervisord
    supervisor

USER verim

COPY ci/local_net/gen_node_configs.sh gen_node_configs.sh
COPY ci/docker-testnet/supervisord.conf /etc/supervisord.conf

EXPOSE 26657 26660 26663 26666

RUN ./gen_node_configs.sh

# Update persistent peers list for at least one node
RUN NODE_0_ID=$(verim-noded tendermint show-node-id --home node_configs/node0); \
    NODE_1_ID=$(verim-noded tendermint show-node-id --home node_configs/node1); \
    NODE_2_ID=$(verim-noded tendermint show-node-id --home node_configs/node2); \
    NODE_3_ID=$(verim-noded tendermint show-node-id --home node_configs/node3); \
    PERSISTENT_PEERS="$NODE_0_ID@localhost:26656,$NODE_1_ID@localhost:26658,$NODE_2_ID@localhost:26660,$NODE_3_ID@localhost:26662"; \
    sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/g" node_configs/node0/config/config.toml; \
    sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/g" node_configs/node1/config/config.toml; \
    sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/g" node_configs/node2/config/config.toml; \
    sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/g" node_configs/node3/config/config.toml

# ...
RUN sed -i $extension "s/addr_book_strict = true/addr_book_strict = false/g" node_configs/node0/config/config.toml; \
    sed -i $extension "s/addr_book_strict = true/addr_book_strict = false/g" node_configs/node1/config/config.toml; \
    sed -i $extension "s/addr_book_strict = true/addr_book_strict = false/g" node_configs/node2/config/config.toml; \
    sed -i $extension "s/addr_book_strict = true/addr_book_strict = false/g" node_configs/node3/config/config.toml

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
STOPSIGNAL SIGTERM
