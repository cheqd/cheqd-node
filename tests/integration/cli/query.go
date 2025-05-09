package cli

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
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
