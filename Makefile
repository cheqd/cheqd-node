PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=cheqd-node \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=cheqd-noded \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=cheqd-nodecli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

install: go.sum
		@echo "--> Installing cheqd-noded & cheqd-nodecli"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/cheqd-noded
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/cheqd-nodecli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test:
	@go test -mod=readonly $(PACKAGES)
