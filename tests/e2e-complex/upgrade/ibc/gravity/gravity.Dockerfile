FROM golang:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    curl \
    git \
    libprotobuf-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN git clone --depth 1 --branch v1.6.5 --single-branch https://github.com/Gravity-Bridge/Gravity-Bridge

WORKDIR /app/Gravity-Bridge/module

# Installing the protobuf tooling
RUN make proto-tools

# Install protobufs plugins
RUN go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos && \
    go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos  && \
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0 && \
    go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0

# generate new protobuf files from the definitions, this makes sure the previous instructions worked
# you will need to run this any time you change a proto file
RUN make proto-gen

# build all code, including your newly generated go protobuf file
RUN make

# run all the unit tests
RUN make test


FROM debian:buster as runtime

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    curl \
    git \
    libprotobuf-dev \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Node binary
COPY --from=builder /go/bin/gravity /bin/gravity

ARG UID=1000
ARG GID=1000

ARG USER=gravity
ARG GROUP=gravity

ARG HOME=/home/$USER

# User
RUN groupadd --system --gid $GID $USER && \
    useradd --system --create-home --home-dir $HOME --shell /bin/bash --gid $GROUP --uid $UID $USER

WORKDIR $HOME

RUN chown -R $USER $HOME
USER $USER

COPY ./gravity_init.sh .
RUN bash gravity_init.sh
RUN gravity start

ENTRYPOINT [ "gravity", "start" ]
