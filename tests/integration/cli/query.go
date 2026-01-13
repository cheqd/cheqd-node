package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"cosmossdk.io/errors"
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	globalfeetypes "github.com/noble-assets/globalfee/types"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

var CLIQueryParams = []string{
	"--chain-id", network.ChainID,
	"--output", OutputFormat,
}

var KeyParams = []string{
	"--output", OutputFormat,
	"--keyring-backend", KeyringBackend,
}

func Query(module, query string, queryArgs ...string) (string, error) {
	args := []string{"query", module, query}

	// Common params
	args = append(args, CLIQueryParams...)

	// Other args
	args = append(args, queryArgs...)

	return Exec(args...)
}

func QueryBalance(address, denom string) (sdk.Coin, error) {
	res, err := Query("bank", "balance", address, denom)
	if err != nil {
		return sdk.Coin{}, err
	}

	var resp banktypes.QueryBalanceResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return sdk.Coin{}, err
	}

	return *resp.Balance, nil
}

func QueryDidDoc(did string) (didtypes.QueryDidDocResponse, error) {
	res, err := Query("cheqd", "did-document", did)
	if err != nil {
		return didtypes.QueryDidDocResponse{}, err
	}

	var resp didtypes.QueryDidDocResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return didtypes.QueryDidDocResponse{}, err
	}

	return resp, nil
}

func QueryDidParams() (didtypes.QueryParamsResponse, error) {
	res, err := Query("cheqd", "params")
	if err != nil {
		return didtypes.QueryParamsResponse{}, err
	}

	var resp didtypes.QueryParamsResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return didtypes.QueryParamsResponse{}, err
	}

	return resp, nil
}

func QueryResource(collectionID string, resourceID string) (resourcetypes.QueryResourceResponse, error) {
	res, err := Query("resource", "specific-resource", collectionID, resourceID)
	if err != nil {
		return resourcetypes.QueryResourceResponse{}, err
	}

	var resp resourcetypes.QueryResourceResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryResourceResponse{}, err
	}

	return resp, nil
}

func QueryResourceMetadata(collectionID string, resourceID string) (resourcetypes.QueryResourceMetadataResponse, error) {
	res, err := Query("resource", "metadata", collectionID, resourceID)
	if err != nil {
		return resourcetypes.QueryResourceMetadataResponse{}, err
	}

	var resp resourcetypes.QueryResourceMetadataResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryResourceMetadataResponse{}, err
	}

	return resp, nil
}

func QueryResourceCollection(collectionID string) (resourcetypes.QueryCollectionResourcesResponse, error) {
	res, err := Query("resource", "collection-metadata", collectionID)
	if err != nil {
		return resourcetypes.QueryCollectionResourcesResponse{}, err
	}

	var resp resourcetypes.QueryCollectionResourcesResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryCollectionResourcesResponse{}, err
	}

	return resp, nil
}

func QueryResourceParams() (resourcetypes.QueryParamsResponse, error) {
	res, err := Query("resource", "params")
	if err != nil {
		return resourcetypes.QueryParamsResponse{}, err
	}

	var resp resourcetypes.QueryParamsResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryParamsResponse{}, err
	}

	return resp, nil
}

func QueryTxn(hash string) (sdk.TxResponse, error) {
	res, err := Query("tx", hash)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	var resp sdk.TxResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return resp, nil
}

func QueryProposal(container, id string) (govtypesv1.Proposal, error) {
	fmt.Println("Querying proposal from", container)
	args := append([]string{
		CliBinaryName,
		"query", "gov", "proposal", id,
	}, QueryParamsConst...)

	out, err := LocalnetExecExec(container, args...)
	if err != nil {
		return govtypesv1.Proposal{}, err
	}

	// FIX: getting type instead of @type in messages struct when querying proposal via cli
	convertedJSON, err := convertProposalJSON(out)
	if err != nil {
		return govtypesv1.Proposal{}, err
	}

	var resp govtypesv1.QueryProposalResponse
	err = MakeCodecWithExtendedRegistry().UnmarshalJSON(convertedJSON, &resp)
	if err != nil {
		return govtypesv1.Proposal{}, err
	}
	return *resp.Proposal, nil
}

func QueryFeemarketGasPrice(denom string) (feemarkettypes.GasPriceResponse, error) {
	res, err := Query("feemarket", "gas-price", denom)
	if err != nil {
		return feemarkettypes.GasPriceResponse{}, err
	}

	var resp feemarkettypes.GasPriceResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return feemarkettypes.GasPriceResponse{}, err
	}

	return resp, nil
}

func QueryFeemarketGasPrices() (feemarkettypes.GasPricesResponse, error) {
	res, err := Query("feemarket", "gas-prices")
	if err != nil {
		return feemarkettypes.GasPricesResponse{}, err
	}

	var resp feemarkettypes.GasPricesResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return feemarkettypes.GasPricesResponse{}, err
	}

	return resp, nil
}

func QueryFeemarketParams() (feemarkettypes.Params, error) {
	res, err := Query("feemarket", "params")
	if err != nil {
		return feemarkettypes.Params{}, err
	}

	var resp feemarkettypes.Params
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return feemarkettypes.Params{}, err
	}

	return resp, nil
}

func GetProposalID(events []abcitypes.Event) (string, error) {
	// Iterate over events
	for _, event := range events {
		// Look for the "submit_proposal" event type
		if event.Type == "submit_proposal" {
			for _, attr := range event.Attributes {
				// Look for the "proposal_id" attribute
				if attr.Key == "proposal_id" {
					return attr.Value, nil
				}
			}
		}
	}

	return "", fmt.Errorf("proposal_id not found")
}

// QueryKeys retrieves the key information and extracts the address
func QueryKeys(name string) (string, error) {
	args := []string{"keys", "show", name}

	args = append(args, KeyParams...)

	output, err := Exec(args...)
	if err != nil {
		return "", err
	}

	var result struct {
		Address string `json:"address"`
	}

	err = json.Unmarshal([]byte(output), &result)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal JSON")
	}

	return result.Address, nil
}

func QueryOracleParams() (oracletypes.QueryParamsResponse, error) {
	res, err := Query("oracle", "params")
	if err != nil {
		return oracletypes.QueryParamsResponse{}, err
	}

	var resp oracletypes.QueryParamsResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return oracletypes.QueryParamsResponse{}, err
	}

	return resp, nil
}

// QueryAggregateVote queries the aggregate vote for a validator
func QueryAggregateVote(validatorAddr string) (*oracletypes.QueryAggregateVoteResponse, error) {
	res, err := Query("oracle", "aggregate-votes", validatorAddr)
	if err != nil {
		return nil, err
	}
	var resp oracletypes.QueryAggregateVoteResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling aggregate vote response: %w", err)
	}

	return &resp, nil
}

// QueryBypassMessages queries the global fee bypass messages
func QueryBypassMessages() ([]string, error) {
	res, err := Query("globalfee", "bypass-messages")
	if err != nil {
		return nil, err
	}

	var resp globalfeetypes.QueryBypassMessagesResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return nil, err
	}

	return resp.BypassMessages, nil
}

// QueryAggregatePrevote queries the aggregate prevote for a validator
func QueryAggregatePrevote(validatorAddr string) (*oracletypes.QueryAggregatePrevoteResponse, error) {
	res, err := Query("oracle", "aggregate-prevotes", validatorAddr)
	if err != nil {
		return nil, err
	}

	var resp oracletypes.QueryAggregatePrevoteResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling aggregate prevote response: %w", err)
	}

	return &resp, nil
}

// QueryExchangeRates queries all exchange rates
func QueryExchangeRates() (*oracletypes.QueryExchangeRatesResponse, error) {
	res, err := Query("oracle", "exchange-rates")
	if err != nil {
		return nil, err
	}

	var resp oracletypes.QueryExchangeRatesResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling exchange rates response: %w", err)
	}

	return &resp, nil
}

// QueryExchangeRate queries the exchange rate for a specific denom
func QueryExchangeRate(denom string) (*oracletypes.QueryExchangeRatesResponse, error) {
	res, err := Query("oracle", "exchange-rate", denom)
	if err != nil {
		return nil, err
	}

	var resp oracletypes.QueryExchangeRatesResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling exchange rate response: %w", err)
	}

	return &resp, nil
}

// QueryFeederDelegation queries the feeder delegation for a validator
func QueryFeederDelegation(validatorAddr string) (*oracletypes.QueryFeederDelegationResponse, error) {
	res, err := Query("oracle", "feeder-delegation", validatorAddr)
	if err != nil {
		return nil, err
	}

	var resp oracletypes.QueryFeederDelegationResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling feeder delegation response: %w", err)
	}

	return &resp, nil
}

// QueryMissCounter queries the miss counter for a validator
func QueryMissCounter(validatorAddr string) (*oracletypes.QueryMissCounterResponse, error) {
	res, err := Query("oracle", "miss-counter", validatorAddr)
	if err != nil {
		return nil, err
	}

	var resp oracletypes.QueryMissCounterResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling miss counter response: %w", err)
	}

	return &resp, nil
}

// QuerySlashWindow queries the current slash window progress
func QuerySlashWindow() (*oracletypes.QuerySlashWindowResponse, error) {
	res, err := Query("oracle", "slash-window")
	if err != nil {
		return nil, err
	}

	var resp oracletypes.QuerySlashWindowResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling slash window response: %w", err)
	}

	return &resp, nil
}

// FundAccount sends tokens to the specified account address to allow it to pay for transaction fees
func FundAccount(recipientAddr, fromAddr string, amount string, feeParams []string) (sdk.TxResponse, error) {
	// Execute a bank send transaction to fund the account
	return Tx("bank", "send", fromAddr, feeParams, fromAddr, recipientAddr, amount)
}

// EnsureAccountFunded checks if an account has sufficient funds and funds it if needed
func EnsureAccountFunded(accountAddr, funderAddr string, minAmount string, feeParams []string) error {
	// First check the account balance
	balance, err := QueryBankBalance(accountAddr, didtypes.BaseMinimalDenom)
	if err != nil {
		// If we can't query the balance, try funding anyway
		fundingAmount := "5000000000ncheq" // 5000 CHEQ in ncheq
		_, err = FundAccount(accountAddr, funderAddr, fundingAmount, feeParams)
		return err
	}

	// Parse the minimum required amount
	requiredAmount, err := parseAmount(minAmount)
	if err != nil {
		return err
	}

	// If balance is less than required, fund the account
	if balance < requiredAmount {
		fundingAmount := fmt.Sprintf("%dncheq", requiredAmount*2) // Fund with twice the minimum needed
		_, err = FundAccount(accountAddr, funderAddr, fundingAmount, feeParams)
		return err
	}

	// Account has sufficient funds
	return nil
}

// QueryBankBalance queries the balance of an account for a specific denom
func QueryBankBalance(accountAddr, denom string) (int64, error) {
	// Construct the command to query the balance
	output, err := Exec("query", "bank", "balances", accountAddr, "--denom", denom, "--output", "json")
	if err != nil {
		return 0, err
	}

	// Parse the output to extract the balance
	// This is a simplified example - the actual parsing would depend on the JSON structure
	if strings.Contains(output, "amount") {
		var amount int64
		_, err := fmt.Sscanf(output, `{"amount":"%d"`, &amount)
		if err != nil {
			return 0, err
		}
		return amount, nil
	}

	// If no balance found, return 0
	return 0, nil
}

// parseAmount parses an amount string like "1000000ncheq" into a numeric value
func parseAmount(amount string) (int64, error) {
	var value int64
	_, err := fmt.Sscanf(amount, "%d", &value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func QueryEMA(denom string) (*oracletypes.QueryWMAResponse, error) {
	res, err := Query("oracle", "ema", denom)
	if err != nil {
		return nil, err
	}

	var resp oracletypes.QueryWMAResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling exchange rate response: %w", err)
	}

	return &resp, nil
}

func QuerySMA(denom string) (*oracletypes.QuerySMAResponse, error) {
	res, err := Query("oracle", "sma", denom)
	if err != nil {
		return nil, err
	}

	var resp oracletypes.QuerySMAResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling exchange rate response: %w", err)
	}

	return &resp, nil
}

func QueryWMA(denom string, strategy string, customWeights []int64) (*oracletypes.QueryWMAResponse, error) {
	// Base CLI args
	args := []string{denom, "--strategy", strategy}

	// Only add weights if strategy is CUSTOM
	if strings.ToUpper(strategy) == "CUSTOM" {
		if len(customWeights) == 0 {
			return nil, fmt.Errorf("custom weights must be provided for CUSTOM strategy")
		}

		var weightsStr []string
		for _, w := range customWeights {
			weightsStr = append(weightsStr, fmt.Sprintf("%d", w))
		}
		args = append(args, "--weights", strings.Join(weightsStr, ","))
	}

	// Run CLI query
	res, err := Query("oracle", "wma", args...)
	if err != nil {
		return nil, fmt.Errorf("failed to run WMA query: %w", err)
	}
	var resp oracletypes.QueryWMAResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling WMA response: %w", err)
	}

	return &resp, nil
}

func QueryConvertUSDCtoCHEQ(amount string, maType string, wmaStrategy string, customWeights []int64) (*oracletypes.ConvertUSDCtoCHEQResponse, error) {
	// Base CLI args
	args := []string{amount, maType}

	if strings.ToLower(maType) == "wma" {
		if wmaStrategy == "" {
			return nil, fmt.Errorf("wma_strategy must be provided when ma_type is 'wma'")
		}
		args = append(args, wmaStrategy)

		if strings.ToUpper(wmaStrategy) == "CUSTOM" {
			if len(customWeights) == 0 {
				return nil, fmt.Errorf("custom weights must be provided for CUSTOM strategy")
			}
			var weightsStr []string
			for _, w := range customWeights {
				weightsStr = append(weightsStr, fmt.Sprintf("%d", w))
			}
			args = append(args, strings.Join(weightsStr, ","))
		}
	}

	// Run CLI query
	res, err := Query("oracle", "convert-usdc-to-cheq", args...)
	if err != nil {
		return nil, fmt.Errorf("failed to run convert-usdc-to-cheq query: %w", err)
	}

	var resp oracletypes.ConvertUSDCtoCHEQResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling ConvertUSDCtoCHEQ response: %w", err)
	}

	return &resp, nil
}
