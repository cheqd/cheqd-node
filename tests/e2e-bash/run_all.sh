#!/bin/bash

set -euox pipefail

function log() {
  MESSAGE=$1

  echo "*****************************************************************************************"
  echo "* $MESSAGE"
  echo "*****************************************************************************************"
}

# Run all shell tests
TESTS=$(find tests -type f -name "test_*")

for TEST in ${TESTS}
do
  log "Running $TEST"

  if bash "$TEST"
  then
    log "$TEST finished successfully"
  else
    log "$TEST failed"
    exit 1
  fi
done
