FROM debian:buster

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    nano \
    curl \
    wget \
    netcat

# Node binary
COPY --from=osmolabs/osmosis:10 /bin/osmosisd /bin/osmosisd

ARG USER=osmosis
ARG GROUP=osmosis

ARG HOME=/home/$USER

# User
RUN groupadd --system --gid 1000 $USER && \
    useradd --system --create-home --home-dir $HOME --shell /bin/bash --gid $GROUP --uid 1000 $USER

WORKDIR $HOME

RUN chown -R $USER $HOME
USER $USER


ENTRYPOINT [ "osmosisd", "start" ]
