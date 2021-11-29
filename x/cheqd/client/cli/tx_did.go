package cli

import (
	"bufio"
	"crypto/ed25519"
	"encoding/base64"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

const FlagVerKey = "ver-key"

func CmdCreateDid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-did [payload-json] [verification-method-id]",
		Short: "Creates a new DID.",
		Long: "Creates a new DID. [payload-json] is JSON encoded MsgCreateDidPayload. " +
			"Key to sign identity message (verKey) will be taken either from " + FlagVerKey + " flag or interactively." +
			"[verification-method-id] is the DID fragment that points to the verKey.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			payloadJson := args[0]
			verificationMethodId := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var payload types.MsgCreateDidPayload
			err = clientCtx.Codec.UnmarshalJSON([]byte(payloadJson), &payload)
			if err != nil {
				return err
			}

			// Get key if not provided
			verKeyPrivBase64, err := cmd.Flags().GetString(FlagVerKey)
			if err != nil {
				return err
			}

			if verKeyPrivBase64 == "" {
				inBuf := bufio.NewReader(clientCtx.Input)
				verKeyPrivBase64, err = input.GetString("Enter base64 encoded verification key", inBuf)
				if err != nil {
					return err
				}
			}

			// Decode key
			verKeyPrivBytes, err := base64.StdEncoding.DecodeString(verKeyPrivBase64)
			if err != nil {
				return err
			}

			verKeyPriv := ed25519.PrivateKey(verKeyPrivBytes)

			// Sign identity message
			signBytes := payload.GetSignBytes()
			signatureBytes := ed25519.Sign(verKeyPriv, signBytes)

			signInfo := types.SignInfo{
				VerificationMethodId: verificationMethodId,
				Signature: base64.StdEncoding.EncodeToString(signatureBytes),
			}

			// Construct message
			msg := types.MsgCreateDid{
				Payload: 	&payload,
				Signatures: []*types.SignInfo{ &signInfo },
			}

			//Set fee-payer if not set
			err = setFeePayerFromSigner(&clientCtx)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagVerKey, "", "Base64 encoded ed25519 private key to sign identity message with. " +
		"Use for testing purposes only because the key will be stored in shell history. Prefer interactive mode.")

	return cmd
}

func CmdUpdateDid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-did [did]",
		Short: "Update a did",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsDidBase64 := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var msg types.MsgUpdateDid
			data, err := base64.StdEncoding.DecodeString(argsDidBase64)
			if err != nil {
				return err
			}

			if err := msg.Unmarshal(data); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
