module github.com/cheqd/cheqd-node/api/v2

go 1.24.0

require (
	cosmossdk.io/api v0.7.6
	github.com/cosmos/cosmos-proto v1.0.0-beta.5
	github.com/cosmos/gogoproto v1.7.0
	google.golang.org/genproto/googleapis/api v0.0.0-20251202230838-ff82c1b0f217
	google.golang.org/grpc v1.79.3
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/cheqd/cosmos-sdk v0.50.14-height-mismatch-iavl.0.20250808071119-3b33570d853b

	github.com/cosmos/iavl => github.com/cheqd/iavl v1.2.2-uneven-heights.0.20250808065519-2c3d5a9959cc
)
