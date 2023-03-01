module github.com/cheqd/cheqd-node/api/v2

go 1.18

require (
	github.com/cosmos/cosmos-proto v1.0.0-alpha8
	github.com/cosmos/cosmos-sdk/api v0.1.0
	github.com/cosmos/gogoproto v1.4.6
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	golang.org/x/exp v0.0.0-20230131160201-f062dba9d201 // indirect
	golang.org/x/net v0.6.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
)

replace (
	// Keyring replacement from Cosmos SDK v0.46.8
	github.com/99designs/keyring => github.com/cosmos/keyring v1.2.0

	// cosmos-sdk state sync allow fast forward to latest height version
	github.com/cosmos/cosmos-sdk => github.com/cheqd/cosmos-sdk v0.46.10-state-sync

	// iavl allow pruning of uneven heights
	github.com/cosmos/iavl => github.com/cheqd/iavl v0.19.5-cheqd

	// dgrijalva/jwt-go is deprecated and doesn't receive security updates.
	// TODO: remove it: https://github.com/cosmos/cosmos-sdk/issues/13134
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.4.2

	// Fix upstream GHSA-h395-qcrw-5vmq vulnerability.
	// TODO Remove it: https://github.com/cosmos/cosmos-sdk/issues/10409
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.7.0

	// From Cosmos SDK v0.46.8 upstream
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

	// From Cosmos SDK v0.46.8 upstream
	github.com/jhump/protoreflect => github.com/jhump/protoreflect v1.9.0

	// use informal systems fork of tendermint
	github.com/tendermint/tendermint => github.com/informalsystems/tendermint v0.34.26
)
