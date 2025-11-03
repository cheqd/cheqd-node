package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/util"
	"github.com/cheqd/cheqd-node/x/did/client/cli"
	"github.com/cheqd/cheqd-node/x/did/types"
	oraclekeeper "github.com/cheqd/cheqd-node/x/oracle/keeper"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	FlagVersionID = "--version-id"
)

var CLITxParams = []string{
	"--chain-id", network.ChainID,
	"--keyring-backend", KeyringBackend,
	"--output", OutputFormat,
	"--yes",
}

var CliGasParams = []string{
	"--gas", Gas,
	"--gas-adjustment", GasAdjustment,
	"--gas-prices", GasPrices,
}

type jsonTx struct {
	Body       txBody     `json:"body"`
	AuthInfo   txAuthInfo `json:"auth_info"`
	Signatures []string   `json:"signatures,omitempty"`
}

type txBody struct {
	Messages                    []json.RawMessage `json:"messages"`
	Memo                        string            `json:"memo,omitempty"`
	TimeoutHeight               string            `json:"timeout_height,omitempty"`
	ExtensionOptions            []json.RawMessage `json:"extension_options,omitempty"`
	NonCriticalExtensionOptions []json.RawMessage `json:"non_critical_extension_options,omitempty"`
}

type txAuthInfo struct {
	SignerInfos []txSignerInfo `json:"signer_infos"`
	Fee         txFee          `json:"fee"`
}

type txSignerInfo struct {
	PublicKey *json.RawMessage `json:"public_key"`
	ModeInfo  txModeInfo       `json:"mode_info"`
	Sequence  string           `json:"sequence"`
}

type txModeInfo struct {
	Single *txSingleModeInfo `json:"single,omitempty"`
}

type txSingleModeInfo struct {
	Mode string `json:"mode"`
}

type txFee struct {
	Amount   sdk.Coins `json:"amount"`
	GasLimit string    `json:"gas_limit"`
	Payer    string    `json:"payer,omitempty"`
	Granter  string    `json:"granter,omitempty"`
}

func Tx(module, tx, from string, feeParams []string, txArgs ...string) (sdk.TxResponse, error) {
	args := []string{"tx", module, tx}

	// Common params
	args = append(args, CLITxParams...)

	// Fee params
	args = append(args, feeParams...)

	// Cosmos account
	args = append(args, "--from", from)

	// Other args
	args = append(args, txArgs...)

	fmt.Println("args", args)

	output, err := Exec(args...)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	// Skip 'gas estimate: xxx' string, trim 'Successfully migrated key' string
	output = helpers.TrimImportedStdout(output)

	var resp sdk.TxResponse

	err = helpers.Codec.UnmarshalJSON([]byte(output), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return resp, nil
}

func SignTx(tmpDir, from string, msg json.RawMessage, fee sdk.Coins, gasLimit uint64) (string, error) {
	if fee == nil {
		fee = sdk.NewCoins()
	}

	unsigned := jsonTx{
		Body: txBody{
			Messages:      []json.RawMessage{msg},
			TimeoutHeight: "0",
		},
		AuthInfo: txAuthInfo{
			SignerInfos: []txSignerInfo{},
			Fee: txFee{
				Amount:   fee,
				GasLimit: strconv.FormatUint(gasLimit, 10),
				Payer:    "",
				Granter:  "",
			},
		},
	}
	unsigned.Signatures = []string{}

	unsignedJSON, err := json.Marshal(&unsigned)
	if err != nil {
		return "", err
	}

	unsignedFile := helpers.MustWriteTmpFile(tmpDir, unsignedJSON)

	println("Unsigned tx written to:", unsignedFile)

	args := []string{
		"tx", "sign", unsignedFile,
		"--from", from,
		"--chain-id", network.ChainID,
		"--keyring-backend", KeyringBackend,
		"--output", OutputFormat,
	}

	output, err := Exec(args...)
	if err != nil {
		return "", err
	}

	return helpers.TrimImportedStdout(output), nil
}

func BroadcastTx(tmpDir, signedTxJSON string) (sdk.TxResponse, error) {
	signedFile := helpers.MustWriteTmpFile(tmpDir, []byte(signedTxJSON))

	args := []string{
		"tx", "broadcast", signedFile,
		"--output", OutputFormat,
	}

	output, err := Exec(args...)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	output = helpers.TrimImportedStdout(output)

	var resp sdk.TxResponse
	if err := helpers.Codec.UnmarshalJSON([]byte(output), &resp); err != nil {
		return sdk.TxResponse{}, err
	}

	return resp, nil
}

func GrantFees(granter, grantee string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("feegrant", "grant", granter, feeParams, granter, grantee)
}

func RevokeFeeGrant(granter, grantee string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("feegrant", "revoke", granter, feeParams, granter, grantee)
}

func CreateDidDoc(tmpDir string, payload cli.DIDDocument, signInputs []cli.SignInput, versionID, from string, feeParams []string) (sdk.TxResponse, error) {
	payloadJSON, err := json.Marshal(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	payloadWithSignInputs := cli.PayloadWithSignInputs{
		Payload:    payloadJSON,
		SignInputs: signInputs,
	}

	payloadWithSignInputsJSON, err := json.Marshal(&payloadWithSignInputs)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	payloadFile := helpers.MustWriteTmpFile(tmpDir, payloadWithSignInputsJSON)

	if versionID != "" {
		return Tx("cheqd", "create-did", from, feeParams, payloadFile, FlagVersionID, versionID)
	}

	return Tx("cheqd", "create-did", from, feeParams, payloadFile)
}

func UpdateDidDoc(tmpDir string, payload cli.DIDDocument, signInputs []cli.SignInput, versionID, from string, feeParams []string) (sdk.TxResponse, error) {
	payloadJSON, err := json.Marshal(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	payloadWithSignInputs := cli.PayloadWithSignInputs{
		Payload:    payloadJSON,
		SignInputs: signInputs,
	}

	payloadWithSignInputsJSON, err := json.Marshal(&payloadWithSignInputs)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	payloadFile := helpers.MustWriteTmpFile(tmpDir, payloadWithSignInputsJSON)

	if versionID != "" {
		return Tx("cheqd", "update-did", from, feeParams, payloadFile, FlagVersionID, versionID)
	}

	return Tx("cheqd", "update-did", from, feeParams, payloadFile)
}

func DeactivateDidDoc(tmpDir string, payload types.MsgDeactivateDidDocPayload, signInputs []cli.SignInput, versionID, from string, feeParams []string) (sdk.TxResponse, error) {
	payloadJSON, err := helpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	payloadWithSignInputs := cli.PayloadWithSignInputs{
		Payload:    payloadJSON,
		SignInputs: signInputs,
	}

	payloadWithSignInputsJSON, err := json.Marshal(&payloadWithSignInputs)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	payloadFile := helpers.MustWriteTmpFile(tmpDir, payloadWithSignInputsJSON)

	if versionID != "" {
		return Tx("cheqd", "deactivate-did", from, feeParams, payloadFile, FlagVersionID, versionID)
	}

	return Tx("cheqd", "deactivate-did", from, feeParams, payloadFile)
}

func CreateResource(tmpDir string, payload resourcetypes.MsgCreateResourcePayload, signInputs []cli.SignInput, dataFile, from string, feeParams []string) (sdk.TxResponse, error) {
	payloadJSON, err := helpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	payloadWithSignInputs := cli.PayloadWithSignInputs{
		Payload:    payloadJSON,
		SignInputs: signInputs,
	}

	payloadWithSignInputsJSON, err := json.Marshal(&payloadWithSignInputs)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	payloadFile := helpers.MustWriteTmpFile("", payloadWithSignInputsJSON)

	return Tx("resource", "create", from, feeParams, payloadFile, dataFile)
}

func BurnMsg(from string, coins string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("cheqd", "burn", from, feeParams, coins)
}

func SubmitProposalTx(from, pathToDir string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("gov", "submit-proposal", from, feeParams, pathToDir)
}

func VoteProposalTx(from, option, id string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("gov", "vote", from, feeParams, option, id)
}

func SendTokensTx(from, to, amount string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("bank", "send", from, feeParams, from, to, amount)
}

// DelegateFeederAddress delegates a feeder address for a validator
func DelegateFeedConsent(validatorAddr, feederAddr, account string, fees []string) (sdk.TxResponse, error) {
	return Tx("oracle", "delegate-feed-consent", account, fees, validatorAddr, feederAddr)
}

// // DelegateFeedConsent executes the delegate-feed-consent transaction command
// func DelegateFeedConsent(operatorAddr, feederAddr, from string, feeParams []string) (sdk.TxResponse, error) {
// 	return Tx(ModuleName, "delegate-feed-consent", from, feeParams, operatorAddr, feederAddr)
// }

// AggregateExchangeRatePrevote executes the exchange-rate-prevote transaction command
func AggregateExchangeRatePrevote(hash string, validatorAddr, from string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("oracle", "exchange-rate-prevote", from, feeParams, hash, validatorAddr)
}

// AggregateExchangeRateVote executes the exchange-rate-vote transaction command
func AggregateExchangeRateVote(salt string, exchangeRates string, validatorAddr, from string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("oracle", "exchange-rate-vote", from, feeParams, salt, exchangeRates, validatorAddr)
}

func ResolveFeeFromParams(feeRanges []types.FeeRange, useMin bool) (sdk.Coin, error) {
	getFeeAmount := func(fee types.FeeRange) *math.Int {
		if useMin {
			if fee.MinAmount != nil {
				return fee.MinAmount
			}
			return fee.MaxAmount
		} else {
			if fee.MaxAmount != nil {
				return fee.MaxAmount
			}
			return fee.MinAmount
		}
	}

	// 1. Try native (ncheq) fee
	for _, fee := range feeRanges {
		if fee.Denom != types.BaseMinimalDenom {
			continue
		}
		amount := getFeeAmount(fee)
		if amount == nil {
			return sdk.Coin{}, fmt.Errorf("both MinAmount and MaxAmount are nil for %s", fee.Denom)
		}
		return sdk.NewCoin(types.BaseMinimalDenom, *amount), nil
	}

	// 2. Try USD fallback
	for _, fee := range feeRanges {
		if fee.Denom != oracletypes.UsdDenom {
			continue
		}

		price, err := QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		if err != nil {
			return sdk.Coin{}, fmt.Errorf("failed to query WMA price: %w", err)
		}

		amount := getFeeAmount(fee)
		if amount == nil {
			return sdk.Coin{}, fmt.Errorf("both MinAmount and MaxAmount are nil for %s", fee.Denom)
		}

		converted, err := ConvertUsdToCheq(*amount, price.Price)
		if err != nil {
			return sdk.Coin{}, fmt.Errorf("failed to convert usd to ncheq: %w", err)
		}
		return sdk.NewCoin(types.BaseMinimalDenom, converted), nil
	}

	return sdk.Coin{}, fmt.Errorf("no valid fee param found with ncheq or usd denom")
}

func ConvertUsdToCheq(usdAmt math.Int, cheqPrice math.LegacyDec) (math.Int, error) {
	if cheqPrice.IsZero() {
		return math.ZeroInt(), fmt.Errorf("cheq price is zero")
	}

	// Convert: 1e18 usd → 1e0 → 1e9 ncheq
	usdDec := math.LegacyNewDecFromInt(usdAmt).Quo(math.LegacyNewDecFromInt(util.UsdExponent)) // convert from 1e18 scale
	ncheqDec := usdDec.Quo(cheqPrice).MulInt64(util.CheqScale.Int64())                         // convert to 1e9 scale
	return ncheqDec.TruncateInt(), nil
}

func IBCAcknowledgementTx(tmpDir string, from string, ack json.RawMessage, gasLimit uint64) (sdk.TxResponse, error) {
	signed, err := SignTx(tmpDir, from, ack, sdk.NewCoins(), gasLimit)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return BroadcastTx(tmpDir, signed)
}

func IBCUpdateClientTx(tmpDir string, from string, msg json.RawMessage, gasLimit uint64) (sdk.TxResponse, error) {
	signed, err := SignTx(tmpDir, from, msg, sdk.NewCoins(), gasLimit)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return BroadcastTx(tmpDir, signed)
}

func IBCRecvPacketTx(tmpDir string, from string, msg json.RawMessage, gasLimit uint64) (sdk.TxResponse, error) {
	signed, err := SignTx(tmpDir, from, msg, sdk.NewCoins(), gasLimit)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return BroadcastTx(tmpDir, signed)
}

func IBCTimeoutTx(tmpDir string, from string, msg json.RawMessage, gasLimit uint64) (sdk.TxResponse, error) {
	signed, err := SignTx(tmpDir, from, msg, sdk.NewCoins(), gasLimit)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return BroadcastTx(tmpDir, signed)
}
