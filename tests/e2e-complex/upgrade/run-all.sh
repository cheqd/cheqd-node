#!/bin/bash

set -euo pipefail

BASE_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# shellcheck disable=SC1091
. "${BASE_DIR}/common.sh"


# Test modules
MODULES=("cheqd-1" "resource-2")


function run_module_script() {
    MODULE=$1
    STAGE=$2

    if [[ -f "${BASE_DIR}/${MODULE}/${STAGE}.sh" ]]; then
        echo "=> Run ${STAGE} handler for ${MODULE} module"
        "${BASE_DIR}/${MODULE}/${STAGE}.sh"
    else
        echo "=> Skip ${STAGE} handler for ${MODULE} module"
    fi
}


echo "===> Run setup"

bash "setup.sh"

for MODULE in "${MODULES[@]}"; do
    run_module_script "${MODULE}" "setup"
done


echo "===> Run before upgrade handlers"

for MODULE in "${MODULES[@]}"; do
    run_module_script "${MODULE}" "before-upgrade"
done


echo "===> Run upgrade"

bash "upgrade.sh"


echo "===> Run after upgrade handlers"

for MODULE in "${MODULES[@]}"; do
    run_module_script "${MODULE}" "after-upgrade"
done


echo "===> Run cleanup"

for MODULE in "${MODULES[@]}"; do
    run_module_script "${MODULE}" "cleanup"
done

bash "cleanup.sh"
