#!/usr/bin/make -f

###############################################################################
###                                Protobuf                                 ###
###############################################################################

containerProtoVer=0.13.2
containerProtoImage=ghcr.io/cosmos/proto-builder:$(containerProtoVer)
containerProtoGen=cheqd-node-proto-gen-$(containerProtoVer)
containerProtoFmt=cheqd-node-proto-fmt-$(containerProtoVer)
containerProtoGenSwagger=cheqd-node-proto-gen-swagger-$(containerProtoVer)
containerPulsar=cheqd-node-pulsar-gen-$(containerProtoVer)

proto-all: proto-gen proto-swagger-gen

proto-gen:
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protocgen.sh; fi

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoFmt}$$"; then docker start -a $(containerProtoFmt); else docker run --name $(containerProtoFmt) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		find .  -name "*.proto" -not -path "./third_party/*" -exec clang-format -i {} \; ; fi

DOCKER_BUF := docker run -v $(shell pwd):/workspace --workdir /workspace bufbuild/buf:1.21.0

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGenSwagger}$$"; then docker start -a $(containerProtoGenSwagger); else docker run --name $(containerProtoGenSwagger) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protoc-swagger-gen.sh; fi

proto-pulsar-gen:
	@echo "Generating Pulsar"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerPulsar}$$"; then docker start -a $(containerPulsar); else docker run --name $(containerPulsar) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protoc-pulsar-gen.sh; fi

.PHONY: proto-all proto-gen proto-format proto-lint proto-swagger-gen