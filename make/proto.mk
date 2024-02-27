#!/usr/bin/make -f

###############################################################################
###                                Protobuf                                 ###
###############################################################################

DOCKER := $(shell which docker)
containerProtoVer=0.14.0
containerProtoImage=ghcr.io/cosmos/proto-builder:$(containerProtoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage)

proto-all: proto-gen proto-swagger-gen

proto-gen:
	sudo find ../ -type d -exec chmod 777 {} +
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh;
	go mod tidy
	cd api
	go mod tidy

proto-format:
	@echo "Formatting Protobuf files"
	@$(protoImage) find .  -name "*.proto" -not -path "./third_party/*" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-swagger-gen:
	sudo find ../ -type d -exec chmod 777 {} +
	sudo chmod 666 ./app/client/docs/swagger.yaml
	@echo "Generating Protobuf Swagger"
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh;

proto-pulsar-gen:
	@echo "Generating Pulsar"
	@$(protoImage) sh ./scripts/protoc-pulsar-gen.sh;

.PHONY: proto-all proto-gen proto-format proto-lint proto-swagger-gen