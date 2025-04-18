#!/usr/bin/env bash

# this script is for generating protobuf files for the new google.golang.org/protobuf API
set -euox pipefail

protoc_gen_install() {
  go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest #2>/dev/null
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest #2>/dev/null
  go install cosmossdk.io/orm/cmd/protoc-gen-go-cosmos-orm@latest #2>/dev/null
}

protoc_gen_install

echo "Generating API module"

# Remove old generated API/Pulsar files
(cd api; find ./ -type f \( -iname \*.pulsar.go -o -iname \*.pb.go -o -iname \*.cosmos_orm.go -o -iname \*.pb.gw.go \) -delete; find . -empty -type d -delete; cd ..)
cd proto

# Find all proto files but exclude "v1" paths
proto_dirs=$(find ./ -type f -path '*/v1/*' -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for proto_dir in $proto_dirs; do
  proto_files=$(find "${proto_dir}" -maxdepth 1 -name '*.proto')
  for f in $proto_files; do
    buf generate --template buf.gen.pulsar.yaml "$f"

    # replace github.com/cheqd/cheqd-node/api/v2/cosmos/base with cosmossdk.io/api/cosmos/base
    sed -i 's/github.com\/cheqd\/cheqd-node\/api\/v2\/cosmos\/base/cosmossdk.io\/api\/cosmos\/base/g' "../api/${f%.*}.pulsar.go"

    # replace amino
    sed -i 's/github.com\/cheqd\/cheqd-node\/api\/v2\/amino/cosmossdk.io\/api\/amino/g' "../api/${f%.*}.pulsar.go"

    # replace msg/v1
    sed -i 's/github.com\/cheqd\/cheqd-node\/api\/v2\/cosmos\/msg\/v1/cosmossdk.io\/api\/cosmos\/msg\/v1/g' "../api/${f%.*}.pulsar.go"
  done
done
