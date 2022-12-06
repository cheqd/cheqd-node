#!/usr/bin/make -f

###############################################################################
###                                Protobuf                                 ###
###############################################################################

containerProtoVer=v0.7
containerProtoImage=tendermintdev/sdk-proto-gen:$(containerProtoVer)
containerProtoGen=cheqd-node-proto-gen-$(containerProtoVer)
containerProtoFmt=cheqd-node-proto-fmt-$(containerProtoVer)
containerProtoGenSwagger=cheqd-node-proto-gen-swagger-$(containerProtoVer)

proto-all: proto-lint proto-format proto-gen proto-swagger-gen

proto-gen:
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protocgen.sh; fi

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoFmt}$$"; then docker start -a $(containerProtoFmt); else docker run --name $(containerProtoFmt) -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find .  -name "*.proto" -not -path "./third_party/*" -exec clang-format -i {} \; ; fi

DOCKER_BUF := docker run -v $(shell pwd):/workspace --workdir /workspace bufbuild/buf:1.7.0

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGenSwagger}$$"; then docker start -a $(containerProtoGenSwagger); else docker run --name $(containerProtoGenSwagger) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protoc-swagger-gen.sh; fi

.PHONY: proto-all proto-gen proto-format proto-lint proto-swagger-gen