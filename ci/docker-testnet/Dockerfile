FROM verim-node

COPY ci/local_net/gen_node_configs.sh ./ci/gen_node_configs.sh
COPY ci/docker-testnet/supervisord.conf /etc/supervisord.conf

EXPOSE 9701 9702 9703 9704 9705 9706 9707 9708

RUN ./ci/gen_node_configs.sh

CMD ["/usr/bin/supervisord"]
STOPSIGNAL SIGTERM
