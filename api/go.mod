module github.com/cheqd/cheqd-node/api/v2

go 1.22.0

toolchain go1.23.4

require (
	cosmossdk.io/api v0.3.1
	github.com/cosmos/cosmos-proto v1.0.0-beta.5
	github.com/cosmos/gogoproto v1.5.0
	google.golang.org/genproto/googleapis/api v0.0.0-20240610135401-a8a62080eff3
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.36.1
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240709173604-40e1e62336c5 // indirect
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/cheqd/cosmos-sdk v0.47.10-height-mismatch

	github.com/cosmos/iavl => github.com/cheqd/iavl v0.20.1-uneven-heights
)
