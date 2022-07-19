FROM ubuntu:22.04

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    nano \
    curl \
    wget \
    netcat

# Node binary
COPY --from=osmolabs/osmosis:10 /bin/osmosisd /bin/osmosisd

ARG UID=1000
ARG GID=1000

ARG USER=osmosis
ARG GROUP=osmosis

ARG HOME=/home/$USER

# User
RUN groupadd --system --gid $GID $USER && \
    useradd --system --create-home --home-dir $HOME --shell /bin/bash --gid $GROUP --uid $UID $USER

WORKDIR $HOME

RUN chown -R $USER $HOME
USER $USER


ENTRYPOINT [ "osmosisd", "start" ]
