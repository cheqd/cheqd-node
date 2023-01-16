#!/usr/bin/env bash

set -euox pipefail

# Define directories
SWAGGER_DIR=./api/docs
SWAGGER_UI_DIR=${SWAGGER_DIR}/swagger-ui

# Define URLs
SWAGGER_UI_VERSION=4.15.5
SWAGGER_UI_DOWNLOAD_URL=https://github.com/swagger-api/swagger-ui/archive/refs/tags/v${SWAGGER_UI_VERSION}.zip
SWAGGER_UI_PACKAGE_NAME=${SWAGGER_DIR}/swagger-ui-${SWAGGER_UI_VERSION}

cd ./proto

# Find all proto files but exclude "v1" paths
proto_dirs=$(find ./cheqd -type f -path '*/v1/*' -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

# loop through proto directories
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' \))
  if [[ -n "$query_file" ]]; then
    buf generate --template buf.gen.swagger.yaml "$query_file"
  fi
done

cd ..

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.yaml` for merging
swagger-combine ${SWAGGER_DIR}/config.json -o ${SWAGGER_DIR}/swagger.yaml -f yaml --continueOnConflictingPaths true

# Remove individual swagger files
rm -rf ${SWAGGER_DIR}/cheqd

# if swagger-ui does not exist locally, download swagger-ui and move dist directory to
# swagger-ui directory, then remove zip file and unzipped swagger-ui directory
if [ ! -d "${SWAGGER_UI_DIR}" ]; then
  # download swagger-ui
  wget -O ${SWAGGER_UI_PACKAGE_NAME}.zip ${SWAGGER_UI_DOWNLOAD_URL}
  # unzip swagger-ui package
  unzip ${SWAGGER_UI_PACKAGE_NAME}.zip -d ${SWAGGER_DIR}
  # move swagger-ui dist directory to swagger-ui directory
  mv ${SWAGGER_UI_PACKAGE_NAME}/dist ${SWAGGER_UI_DIR}
  # remove swagger-ui zip file and unzipped swagger-ui directory
  rm -rf ${SWAGGER_UI_PACKAGE_NAME}.zip ${SWAGGER_UI_PACKAGE_NAME}
fi

# Copy generated swagger yaml file to swagger-ui directory
cp ${SWAGGER_DIR}/swagger.yaml ${SWAGGER_DIR}/swagger-ui/

# update swagger initializer to default to swagger.yaml
# Note: using -i.bak makes this compatible with both GNU and BSD/Mac
sed -i.bak "s|https://petstore.swagger.io/v2/swagger.json|swagger.yaml|" ${SWAGGER_UI_DIR}/swagger-initializer.js

# install statik
go install github.com/rakyll/statik@latest

# generate statik golang code using updated swagger-ui directory
statik -src=${SWAGGER_DIR}/swagger-ui -dest=${SWAGGER_DIR} -f -m

# log whether or not the swagger directory was updated
if [ -n "$(git status ${SWAGGER_DIR} --porcelain)" ]; then
  echo "Swagger statik file updated"
else
  echo "Swagger statik file already in sync"
fi
