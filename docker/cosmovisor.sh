#!/bin/bash
set -euo pipefail

COSMOVISOR_ROOT_DIR=${HOME}/.cheqdnode/cosmovisor

mkdir -p ${HOME}/.cheqdnode
mkdir -p ${COSMOVISOR_ROOT_DIR}
mkdir -p ${COSMOVISOR_ROOT_DIR}/genesis
mkdir -p ${COSMOVISOR_ROOT_DIR}/genesis/bin/
mkdir -p ${COSMOVISOR_ROOT_DIR}/upgrades

cp /bin/cheqd-noded ${COSMOVISOR_ROOT_DIR}/genesis/bin/

cosmovisor "$@"