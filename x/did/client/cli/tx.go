package cli

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

type PayloadWithSignInputs struct {
	Payload    json.RawMessage
	SignInputs []SignInput
}

type SignInput struct {
	VerificationMethodId string
	PrivKey              ed25519.PrivateKey
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transaction subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdCreateDidDoc())
	cmd.AddCommand(CmdUpdateDidDoc())
	cmd.AddCommand(CmdDeactivateDidDoc())

	return cmd
}

func SignWithSignInputs(signBytes []byte, signInputs []SignInput) []*types.SignInfo {
	var signatures []*types.SignInfo

	for _, signInput := range signInputs {
		signatureBytes := ed25519.Sign(signInput.PrivKey, signBytes)

		signInfo := types.SignInfo{
			VerificationMethodId: signInput.VerificationMethodId,
			Signature:            signatureBytes,
		}

		signatures = append(signatures, &signInfo)
	}

	return signatures
}

func SetFeePayerFromSigner(ctx *client.Context) error {
	if ctx.FromAddress != nil {
		ctx.FeePayer = ctx.FromAddress
		return nil
	}

	signerAccAddr, err := AccAddrByKeyRef(ctx.Keyring, ctx.From)
	if err != nil {
		return err
	}

	ctx.FeePayer = signerAccAddr
	return nil
}

func AccAddrByKeyRef(keyring keyring.Keyring, keyRef string) (sdk.AccAddress, error) {
	// Firstly check if the keyref is a key name of a key registered in a keyring
	info, err := keyring.Key(keyRef)

	if err == nil {
		return info.GetAddress()
	}

	if !sdkerrors.IsOf(err, sdkerrors.ErrIO, sdkerrors.ErrKeyNotFound) {
		return nil, err
	}

	// Fallback: convert keyref to address
	return sdk.AccAddressFromBech32(keyRef)
}

func ReadPayloadWithSignInputsFromFile(filePath string) (json.RawMessage, []SignInput, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	payloadWithSignInputs := &PayloadWithSignInputs{}
	err = json.Unmarshal(bytes, payloadWithSignInputs)
	if err != nil {
		return nil, nil, err
	}

	return payloadWithSignInputs.Payload, payloadWithSignInputs.SignInputs, nil
}
