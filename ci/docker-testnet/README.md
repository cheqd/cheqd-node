# Docker

To build docker image:

- Go to the repository root
- Run `docker build -f ci/docker/Dockerfile -t verim-node .` to build verim node
- Run `docker build -f ci/docker-testnet/Dockerfile -t testnet .` to build testnet
