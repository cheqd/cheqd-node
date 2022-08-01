#!/bin/bash

set -euox pipefail

# Get protoc-gen-gocosmos
go get github.com/cosmos/gogoproto 2>/dev/null

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find ./ -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  find "${dir}" -maxdepth 1 -name '*.proto' | while IFS= read -r -d '' file; do
    if grep go_package $file &> /dev/null ; then
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

cd ..

# move proto files to the right places
cp -r github.com/cheqd/cheqd-node/* ./
rm -rf github.com

go mod tidy

./scripts/protocgen2.sh
