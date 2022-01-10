#!/bin/bash

# Ensures that all bash scripts in the repository use `set -euox pipefail` or `set -euo pipefail` statement at the beginning.

set -euo pipefail

INVALID_FILES_FOUND=0

for BASH_SCRIPT in $(find . -type f -name "*.sh")
do
    if ( ! grep -q "set -euo pipefail" "${BASH_SCRIPT}" ) && ( ! grep -q "set -euox pipefail" "${BASH_SCRIPT}" )
    then
        echo "${BASH_SCRIPT}"
        INVALID_FILES_FOUND=1
    fi
done

if [[ INVALID_FILES_FOUND ]]
then
    echo ""
    echo "The bash scripts above must include either 'set -euo pipefail' or 'set -euox pipefail."
fi
