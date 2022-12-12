#!/bin/bash

set -euox pipefail

DIR_="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

bash "$DIR_/cleanup.sh"



bash "$DIR_/setup.sh"

ginkgo -r --tags upgrade --race --tags upgrade_integration --focus-file pre_test.go

bash "$DIR_/upgrade.sh"

ginkgo -r --tags upgrade --race --tags upgrade_integration --focus-file post_test.go

bash "$DIR_/cleanup.sh"