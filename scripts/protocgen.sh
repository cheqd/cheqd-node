#!/bin/bash

set -euox pipefail

# Get protoc executions
go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos 2>/dev/null

# Get cosmos sdk from github
go get github.com/cosmos/cosmos-sdk 2>/dev/null

# Get the path of the cosmos-sdk repo from go/pkg/mod
cosmos_sdk_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)

proto_dirs=$(find . -path ./third_party -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

for dir in $proto_dirs; do
  read -ra proto_files < <(find "${dir}" -name '*.proto')

  # Generate protobuf bind
  buf protoc \
  -I "proto" \
  -I "$cosmos_sdk_dir/third_party/proto" \
  -I "$cosmos_sdk_dir/proto" \
  --gocosmos_out=plugins=interfacetype+grpc,\
Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  "${proto_files[@]}"

  read -ra proto_files < <(find "${dir}" -maxdepth 1 -name '*.proto')

  # Generate grpc gateway
  buf protoc \
  -I "proto" \
  -I "$cosmos_sdk_dir/third_party/proto" \
  -I "$cosmos_sdk_dir/proto" \
  --grpc-gateway_out=logtostderr=true,allow_colon_final_segments=true:. \
  "${proto_files[@]}"
done

cp -r ./github.com/cheqd/cheqd-node/* ./
rm -rf ./github.com
