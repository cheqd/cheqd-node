#!/usr/bin/env bash

set -euox pipefail

echo "Generating gogo proto code"
cd proto

proto_dirs=$(find ./ -type f -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for proto_dir in $proto_dirs; do
  proto_files=$(find "${proto_dir}" -maxdepth 1 -name '*.proto')
  for f in $proto_files; do
    if grep go_package "$f" &>/dev/null; then
      buf generate --template buf.gen.gogo.yaml "$f"
    fi
  done
done

ls -l
cd ..
ls -l


# move proto files to the right places
cp -r github.com/cheqd/cheqd-node/* ./
rm -rf github.com

sh ./scripts/protoc-pulsar-gen.sh
