#####  Build container  #####

FROM verim-node

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    # Common
    curl \
    # Protoc
    protobuf-compiler \
    libprotobuf-dev \
	git \
	wget \
	python3.5 \
	python3-pip \
	python-setuptools \
	python3-nacl \
	apt-transport-https \
	ca-certificates \
	supervisor

RUN pip3 install -U \
	pip==9.0.3 \
	setuptools

RUN useradd -ms /bin/bash -u 1000 verim

# Starport
RUN curl https://get.starport.network/starport! | bash

# App
WORKDIR /app

COPY app ./app
COPY cmd ./cmd
COPY proto ./proto
COPY vue ./vue
COPY x ./x
COPY ci/local_net/gen_node_configs.sh ./ci/gen_node_configs.sh
COPY ci/docker-testnet/supervisord.conf /etc/supervisord.conf
COPY go.mod .
COPY go.sum .
RUN mkdir node_configs

RUN starport build

EXPOSE 9701 9702 9703 9704 9705 9706 9707 9708

#####  Run container  #####

RUN ./ci/gen_node_configs.sh

CMD ["/usr/bin/supervisord"]
STOPSIGNAL SIGTERM
