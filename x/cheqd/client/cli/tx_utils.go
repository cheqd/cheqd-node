package cli

import (
	"bufio"
	"crypto/ed25519"
	"encoding/base64"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)


func setFeePayerFromSigner(ctx *client.Context) error {
	if ctx.FromAddress != nil {
		ctx.FeePayer = ctx.FromAddress
		return nil
	}

	signerAccAddr, err := accAddrByKeyRef(ctx.Keyring, ctx.From)
	if err != nil {
		return err
	}

	ctx.FeePayer = signerAccAddr
	return nil
}

func accAddrByKeyRef(keyring keyring.Keyring, keyRef string) (sdk.AccAddress, error) {
	// Firstly check if the keyref is a key name of a key registered in a keyring
	info, err := keyring.Key(keyRef)

	if err == nil {
		return info.GetAddress(), nil
	}

	if !sdkerr.IsOf(err, sdkerr.ErrIO, sdkerr.ErrKeyNotFound) {
		return nil, err
	}

	// Fallback: convert keyref to address
	return sdk.AccAddressFromBech32(keyRef)
}

func getVerKey(cmd *cobra.Command, err error, clientCtx client.Context) (ed25519.PrivateKey, error) {
	// Try getting from arg
	verKeyPrivBase64, err := cmd.Flags().GetString(FlagVerKey)
	if err != nil {
		return nil, err
	}

	// Get interactively instead if the flag isn't provided
	if verKeyPrivBase64 == "" {
		inBuf := bufio.NewReader(clientCtx.Input)
		verKeyPrivBase64, err = input.GetString("Enter base64 encoded verification key", inBuf)
		if err != nil {
			return nil, err
		}
	}

	// Decode key
	verKeyPrivBytes, err := base64.StdEncoding.DecodeString(verKeyPrivBase64)
	if err != nil {
		return nil, err
	}

	return verKeyPrivBytes, nil
}
