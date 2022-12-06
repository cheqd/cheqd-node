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

cd ..

# install swagger-combine if not already installed
npm list -g | grep swagger-combine > /dev/null || npm install -g swagger-combine --no-shrinkwrap

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./client/docs/config.json -o ./client/docs/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# clean swagger files
rm -rf ./tmp-swagger-gen