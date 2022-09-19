#!/bin/bash

set -euo pipefail

BASE_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# shellcheck disable=SC1091
. "${BASE_DIR}/common.sh"


# Test modules
MODULES=("ibc")


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

function run_module_scripts() {
    STAGE=$1

    for MODULE in "${MODULES[@]}"; do
        run_module_script "${MODULE}" "${STAGE}"
    done
}


echo "===> Run cleanup"
run_module_scripts "cleanup"
"${BASE_DIR}/cleanup.sh"


echo "===> Run setup"
"${BASE_DIR}/setup.sh"
run_module_scripts "setup"


echo "===> Run before upgrade handlers"
run_module_scripts "before-upgrade"


echo "===> Run upgrade"
"${BASE_DIR}/upgrade.sh"


echo "===> Run after upgrade handlers"
run_module_scripts "after-upgrade"


echo "===> Run cleanup"
run_module_scripts "cleanup"
"${BASE_DIR}/cleanup.sh"
