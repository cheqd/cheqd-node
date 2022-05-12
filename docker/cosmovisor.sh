#!/bin/bash

COSMOVISOR_ROOT_DIR=${HOME}/.cheqdnode/cosmovisor

mkdir ${HOME}/.cheqdnode
mkdir ${COSMOVISOR_ROOT_DIR}
mkdir ${COSMOVISOR_ROOT_DIR}/genesis
mkdir ${COSMOVISOR_ROOT_DIR}/genesis/bin/
mkdir ${COSMOVISOR_ROOT_DIR}/upgrades

cp /bin/cheqd-noded ${COSMOVISOR_ROOT_DIR}/genesis/bin/

cosmovisor "$@"