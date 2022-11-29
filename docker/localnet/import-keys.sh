#!/bin/bash

set -euox pipefail

KEYRING_BACKEND="test"

function import_key() {
    ALIAS=${1}
    MNEMONIC=${2}

    echo "Importing key: ${ALIAS}"

    if cheqd-noded keys show "${ALIAS}" --keyring-backend ${KEYRING_BACKEND}
    then
      echo "Key ${ALIAS} already exists"
      return 0
    fi

    echo "${MNEMONIC}" | cheqd-noded keys add "${ALIAS}" --keyring-backend ${KEYRING_BACKEND} --recover
}

import_key "base_account_1" "sketch mountain erode window enact net enrich smoke claim kangaroo another visual write meat latin bacon pulp similar forum guilt father state erase bright"
import_key "base_account_2" "ugly dirt sorry girl prepare argue door man that manual glow scout bomb pigeon matter library transfer flower clown cat miss pluck drama dizzy"
import_key "base_account_3" "fix wheel picnic about army scan table fence device trust alter erupt wear donkey wood slender gold reunion grant quiz absurd tragic reform attitude"
import_key "base_account_4" "horn slim pigeon winner capable piano soul find ignore crawl arrow genuine magnet nasty basic lamp scissors treat stick arm dress elbow trash naive"
import_key "base_account_5" "blue town hobby lens hawk deputy father tissue state choose another liquid license start push iron limb visa taste mother cause history tackle fiber"
import_key "base_account_6" "gallery hospital vicious demand orient piano melody vanish remind pistol elephant bracket olive kitten caution apart capital protect junior endorse run drama tiny patrol"
import_key "base_vesting_account" "coach index fence broken very cricket someone casino dial truth fitness stay habit such three jump exotic spawn planet fragile walk enact angry great"
import_key "continuous_vesting_account" "phone worry flame safe panther dirt picture pepper purchase tiny search theme issue genre orange merit stove spoil surface color garment mind chuckle image"
import_key "delayed_vesting_account" "pilot text keen deal economy donkey use artist divide foster walk pink breeze proud dish brown icon shaft infant level labor lift will tomorrow"
import_key "periodic_vesting_account" "want merge flame plate trouble moral submit wing whale sick meat lonely yellow lens enable oyster slight health vast weird radar mesh grab olive"
