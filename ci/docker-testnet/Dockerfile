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

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
STOPSIGNAL SIGTERM
