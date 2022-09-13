#!/bin/bash

set -euox pipefail

source common.sh

info "Tear down" # ---
docker compose down --timeout 20 --volumes --remove-orphans
