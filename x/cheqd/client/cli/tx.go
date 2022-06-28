package cli

import (
	"bufio"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

type SignInput struct {
	verificationMethodId string
	privKey              ed25519.PrivateKey
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

	cmd.AddCommand(CmdCreateDid())
	cmd.AddCommand(CmdUpdateDid())

	return cmd
}

func GetPayloadAndSignInputs(clientCtx client.Context, args []string) (string, []SignInput, error) {
	// Check for args count
	if len(args)%2 != 1 {
		return "", []SignInput{}, fmt.Errorf("invalid number of arguments: %d. must be an odd number", len(args))
	}

	// Get payload json
	payloadJson := args[0]

	// Get signInputs
	signInputs, err := GetSignInputs(clientCtx, args[1:])
	if err != nil {
		return "", []SignInput{}, err
	}

	return payloadJson, signInputs, nil
}

func GetSignInputs(clientCtx client.Context, args []string) ([]SignInput, error) {
	// Check for args count
	if len(args)%2 != 0 {
		return []SignInput{}, fmt.Errorf("can't read sign inputs. invalid number of arguments: %d", len(args))
	}

	// Get signInputs
	var signInputs []SignInput

	for i := 0; i < len(args); i += 2 {
		vmId := args[i]
		privKey := args[i+1]

		if privKey == "interactive" {
			inBuf := bufio.NewReader(clientCtx.Input)

			var err error
			privKey, err = input.GetString("Enter base64 encoded verification key", inBuf)

			if err != nil {
				return nil, err
			}
		}

		privKeyBytes, err := base64.StdEncoding.DecodeString(privKey)
		if err != nil {
			return nil, fmt.Errorf("unable to decode private key: %s", err.Error())
		}

		signInput := SignInput{
			verificationMethodId: vmId,
			privKey:              privKeyBytes,
		}

		signInputs = append(signInputs, signInput)
	}

	return signInputs, nil
}

func SignWithSignInputs(signBytes []byte, signInputs []SignInput) []*types.SignInfo {
	var signatures []*types.SignInfo

	for _, signInput := range signInputs {
		signatureBytes := ed25519.Sign(signInput.privKey, signBytes)

		signInfo := types.SignInfo{
			VerificationMethodId: signInput.verificationMethodId,
			Signature:            base64.StdEncoding.EncodeToString(signatureBytes),
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
		return info.GetAddress(), nil
	}

	if !sdkerrors.IsOf(err, sdkerrors.ErrIO, sdkerrors.ErrKeyNotFound) {
		return nil, err
	}

	// Fallback: convert keyref to address
	return sdk.AccAddressFromBech32(keyRef)
}
