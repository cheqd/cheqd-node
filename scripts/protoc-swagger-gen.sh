#!/usr/bin/env bash

set -euox pipefail

cd ./proto


# find all proto directories
proto_dirs=$(find ./cheqd -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

# loop through proto directories
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' \))
  if [[ -n "$query_file" ]]; then
    buf generate --template buf.gen.swagger.yaml "$query_file"
  fi
done