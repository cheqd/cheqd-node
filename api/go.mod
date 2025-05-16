module github.com/cheqd/cheqd-node/api/v2

go 1.23.8

require (
	cosmossdk.io/api v0.7.6
	github.com/cosmos/cosmos-proto v1.0.0-beta.5
	github.com/cosmos/gogoproto v1.7.0
	google.golang.org/genproto/googleapis/api v0.0.0-20240814211410-ddb44dafa142
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240930140551-af27646dc61f // indirect
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/cheqd/cosmos-sdk v0.50.13-height-mismatch-patched

	github.com/cosmos/iavl => github.com/cheqd/iavl v0.19.2-0.20250417102206-84370d5811fb
)
