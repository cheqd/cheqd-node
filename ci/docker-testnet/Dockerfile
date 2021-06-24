FROM verim-node

COPY ci/local_net/gen_node_configs.sh gen_node_configs.sh
COPY ci/docker-testnet/supervisord.conf /etc/supervisord.conf

EXPOSE 26657 26660 26663 26666

RUN ./gen_node_configs.sh

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
STOPSIGNAL SIGTERM
