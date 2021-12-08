export VERSION := $(shell echo $(shell git describe --always --match "v*") | sed 's/^v//')
export TMVERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::')
export COMMIT := $(shell git log -1 --format='%H')

NAME ?= cheqd-node
APPNAME ?= $(NAME)d
LEDGER_ENABLED ?= true

BINDIR ?= $(GOPATH)/bin
BUILDDIR ?= $(CURDIR)/build

export GO111MODULE = on


### Process build tags
build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    $(error Windows OS is not supported currently, exiting)
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))


### Process linker flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=$(NAME)\
		  -X github.com/cosmos/cosmos-sdk/version.AppName=$(APPNAME) \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
			-X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TMVERSION)


# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))


### Resulting build flags
BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# Check for debug option
ifeq (debug,$(findstring debug,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -gcflags "all=-N -l"
endif


###############################################################################
###                                  Build                                  ###
###############################################################################


all: install

build: go.sum $(BUILDDIR)
	go build -mod=readonly $(BUILD_FLAGS) -o $(BUILDDIR)/cheqd-noded ./cmd/cheqd-noded

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/cheqd-noded

$(BUILDDIR):
	mkdir -p $(BUILDDIR)

clean:
	rm -rf \
    $(BUILDDIR)/ \
    artifacts/ \
    tmp-swagger-gen/

.PHONY: all clean build install

go.sum: go.mod
	echo "Ensure dependencies have not been modified ..." >&2
	go mod verify
	go mod tidy


###############################################################################
###                                  Protobuf                               ###
###############################################################################

proto-gen:
	bash scripts/protocgen.sh
