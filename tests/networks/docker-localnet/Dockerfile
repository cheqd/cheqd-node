FROM cheqd-node

USER root

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    supervisor \
    jq \
    nano

USER cheqd

COPY --chown=cheqd:cheqd tests/networks/tools/gen-network-config.sh ./
COPY tests/networks/docker-localnet/supervisord.conf /etc

RUN chmod +x ./gen-network-config.sh
RUN ./gen-network-config.sh "cheqd" 2 0 0

# Tweak ports for val-1
RUN sed -i.bak 's|laddr = "tcp://0.0.0.0:26656"|laddr = "tcp://0.0.0.0:26666"|g' network-config/validator-1/config/config.toml; \
    sed -i.bak 's|laddr = "tcp://0.0.0.0:26657"|laddr = "tcp://0.0.0.0:26667"|g' network-config/validator-1/config/config.toml; \
    sed -i.bak 's|pprof_laddr = "localhost:6060"|pprof_laddr = "localhost:6070"|g' network-config/validator-1/config/config.toml; \
    sed -i.bak 's|address = "0.0.0.0:9090"|address = "0.0.0.0:9100"|g' network-config/validator-1/config/app.toml; \
    sed -i.bak 's|address = "0.0.0.0:9091"|address = "0.0.0.0:9101"|g' network-config/validator-1/config/app.toml; \
    sed -i.bak 's|address = "tcp://0.0.0.0:1317"|address = "tcp://0.0.0.0:1327"|g' network-config/validator-1/config/app.toml; \
    sed -i.bak 's|address = ":8080"|address = ":8090"|g' network-config/validator-1/config/app.toml

# Tweak peers
RUN VAL_0_ID=$(cheqd-noded tendermint show-node-id --home "network-config/validator-0"); \
    PERSISTENT_PEERS="$VAL_0_ID@127.0.0.1:26656"; \
    sed -i.bak "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/g" network-config/validator-1/config/config.toml

EXPOSE 26656 26657 9090

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
STOPSIGNAL SIGTERM
