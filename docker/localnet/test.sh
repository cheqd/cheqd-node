#!/bin/bash

set -e


RUNNER_BIN_DIR=./../../cheqd-noded
# Checkout the repository
echo "Checking out the repository..."
git clone https://github.com/your-username/your-repository.git
cd your-repository

# Download binary artifact
echo "Downloading binary artifact..."
mkdir -p $RUNNER_BIN_DIR
wget -O $RUNNER_BIN_DIR/cheqd-noded https://example.com/cheqd-noded
chmod +x $RUNNER_BIN_DIR/cheqd-noded

# Download and load node Docker image
echo "Downloading and loading node Docker image..."
wget -O cheqd-node-build.tar https://example.com/cheqd-node-build.tar
docker load -i cheqd-node-build.tar

# Generate localnet configs
echo "Generating localnet configs..."
cd docker/localnet
bash gen-network-config.sh
sudo chown -R 1000:1000 network-config

# Set up Docker localnet
echo "Setting up Docker localnet..."
docker compose --env-file build-latest.env up --detach --no-build

# Import keys
echo "Importing keys..."
bash import-keys.sh

# Install ginkgo
echo "Installing ginkgo..."
cd ../..
go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Run Ginkgo integration tests
echo "Running Ginkgo integration tests..."
cd tests/integration
ginkgo -r --tags integration --race --randomize-suites --keep-going --trace --junit-report ../../report-integration.xml

# Show logs on failure
if [ $? -ne 0 ]; then
  echo "Test failed. Showing logs..."
  cd ../..
  docker compose --env-file build-latest.env logs --tail --follow
fi

# Upload integration tests result
echo "Uploading integration tests result..."
cd ../../
wget -O report-integration.xml https://example.com/report-integration.xml
