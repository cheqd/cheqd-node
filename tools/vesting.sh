#!/bin/bash

# A script for shareholder vesting accounts creation. Creates post genesis vesting accounts with adjusted vested tokens amount.
# bc is used instead of bash arithmetic because balances are big numbers.
# Inputs: source_address, target_address (with either cheqd of cosmos prefix), total_amount, vesting_amount, vesting_start, vesting_end

set -euo pipefail

# sed in macos requires extra argument
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi


# Constants
RPC_ENDPOINT="http://localhost:26657"
GAS_PRICES="25ncheq"
CHAIN_ID="cheqd"
DENOM="ncheq"

CHEQD_PREFIX="cheqd"
COSMOS_PREFIX="cosmos"


# Read source and target addresses
echo "Source account address or key name:"
read -r SOURCE_ACCOUNT

echo "Target account address:"
read -r TARGET_ADDRESS


# Convert target address to cheqd if needed
if [[ $TARGET_ADDRESS == "${CHEQD_PREFIX}"* ]]
then
  echo "Target address prefix is ${CHEQD_PREFIX}. No need to convert."
elif [[ $TARGET_ADDRESS == "${COSMOS_PREFIX}"* ]]
then
  echo "Target address prefix is ${COSMOS_PREFIX}. Need to convert to ${CHEQD_PREFIX}"

  BYTES=$(cheqd-noded keys parse "${TARGET_ADDRESS}" --output json | jq -r ".bytes")
  TARGET_ADDRESS=$(cheqd-noded keys parse "${BYTES}" --output json | jq -r ".formats[0]")

  echo "Conversion result: ${TARGET_ADDRESS}"
else
  echo  "Unrecognised account prefix. Must be either ${COSMOS_PREFIX} or ${CHEQD_PREFIX}."
  exit 1
fi


# Check that source account exists
if ! cheqd-noded keys auth account "${SOURCE_ACCOUNT}"  > /dev/null 2>&1 ;
then
  echo "Source account doesn't exist"
  exit 1
fi

# Check that target account doesn't exist
if cheqd-noded query auth account "${TARGET_ADDRESS}"  > /dev/null 2>&1 ;
then
  echo "Target account already exists"
  exit 1
fi


# Read vesting parameters
echo "Total amount to send (ncheq):"
read -r TOTAL_AMOUNT

echo "Vested amount (ncheq):"
read -r VESTED_AMOUNT

NOW=$(date +%s)
echo "Current epoch time: ${NOW}"

echo "Vesting start time (epoch):"
read -r VESTING_START

echo "Vesting end time (epoch):"
read -r VESTING_END


# Sanity checks
if [[ $(echo "${VESTING_START} < ${NOW}" | bc) == 0 ]]
then
  echo "Vesting start time must be less then now"
  exit 1
fi

if [[ $(echo "${VESTING_END} > ${VESTING_START}" | bc) == 0 ]]
then
  echo "Vesting start time must be less then vesting end time"
  exit 1
fi


# Calculate adjustments
UNVESTED_AMOUNT=$(echo "${TOTAL_AMOUNT} - ${VESTED_AMOUNT}" | bc)
VESTING_PERIOD=$(echo "${VESTING_END} - ${VESTING_START}" | bc)

VESTING_PERIOD_ADJUSTED=$(echo "${VESTING_END} - ${NOW}" | bc)
VESTED_AMOUNT_ADJUSTED=$(echo "${TOTAL_AMOUNT} - ${UNVESTED_AMOUNT} * ${VESTING_PERIOD_ADJUSTED} / ${VESTING_PERIOD}" | bc)
UNVESTED_AMOUNT_ADJUSTED=$(echo "${TOTAL_AMOUNT} - ${VESTED_AMOUNT_ADJUSTED}" | bc)

echo "Initial values:"
echo "  Vesting period: ${VESTING_PERIOD}"
echo "  Vested tokens: ${VESTED_AMOUNT}"
echo "  Unvested tokens: ${UNVESTED_AMOUNT}"

echo "Adjusted values values:"
echo "  Vesting period: ${VESTING_PERIOD_ADJUSTED}"
echo "  Vested tokens: ${VESTED_AMOUNT_ADJUSTED}"
echo "  Unvested tokens: ${UNVESTED_AMOUNT_ADJUSTED}"


# Send actual transactions
cheqd-noded tx vesting create-vesting-account "${TARGET_ADDRESS}" "${UNVESTED_AMOUNT_ADJUSTED}${DENOM}" "${VESTING_END}" \
  --from "${SOURCE_ACCOUNT}" --gas-prices "${GAS_PRICES}" --chain-id "${CHAIN_ID}" --node "${RPC_ENDPOINT}"

cheqd-noded tx bank send "${SOURCE_ACCOUNT}" "${TARGET_ADDRESS}" "${VESTED_AMOUNT_ADJUSTED}${DENOM}" \
  --gas-prices "${GAS_PRICES}" --chain-id "${CHAIN_ID}" --node "${RPC_ENDPOINT}"


# Check the result
ACTUAL_BALANCE=$(cheqd-noded query bank balances "${TARGET_ADDRESS}" --output json | jq -r ".balances[0].amount")

echo "Expected balance: ${TOTAL_AMOUNT}"
echo "Actual balance: ${ACTUAL_BALANCE}"

if [[ $(echo "${TOTAL_AMOUNT} == ${ACTUAL_BALANCE}" | bc) == 1 ]] ;
then
  echo "Balances match"
else
  echo "Balances doesn't match"
  exit 1
fi
