package cmd

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multibase"
	"github.com/spf13/cobra"
)

func extendDebug(debugCmd *cobra.Command) *cobra.Command {
	debugCmd.AddCommand(ed25519Cmd(),
		encodingCmd())

	return debugCmd
}

// ed25519Cmd returns cobra Command.
func ed25519Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ed25519",
		Short: "ed25519 tools",
	}

	cmd.AddCommand(ed25519RandomCmd(), ed25519publicKeyBase64ToJwkCmd())

	return cmd
}

// ed25519Cmd returns cobra Command.
func ed25519RandomCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "random",
		Short: "Generate random ed25519 keypair. Output is in JSON format, with base64 encoded public and private keys.",
		RunE: func(cmd *cobra.Command, args []string) error {
			publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
			if err != nil {
				return err
			}

			keyInfo := struct {
				PublicKeyBase64  string `json:"public_key_base_64"`
				PrivateKeyBase64 string `json:"private_key_base_64"`
			}{
				PublicKeyBase64:  base64.StdEncoding.EncodeToString(publicKey),
				PrivateKeyBase64: base64.StdEncoding.EncodeToString(privateKey),
			}

			keyInfoJSON, err := json.Marshal(keyInfo)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(keyInfoJSON))
			return err
		},
	}

	return cmd
}

// ed25519publicKeyBase64ToJwk returns cobra Command.
func ed25519publicKeyBase64ToJwkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "base64-jwk",
		Short: `Convert ed25519 public key from base64 to Json Web Key, according to JsonWebKey2020 spec.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			publicKeyBase64 := args[0]
			publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
			if err != nil {
				return err
			}

			publicKey := ed25519.PublicKey(publicKeyBytes)

			publicKeyJwk, err := jwk.New(publicKey)
			if err != nil {
				return err
			}

			publicKeyJwkJSON, err := json.Marshal(publicKeyJwk)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(publicKeyJwkJSON))
			return err
		},
	}

	return cmd
}

// encoding returns cobra Command.
func encodingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encoding",
		Short: "Encoding tools",
	}

	cmd.AddCommand(base64toMultibaseCmd())
	cmd.AddCommand(base64toBase58Cmd())
	return cmd
}

// base64toMultibaseCmd returns cobra Command.
func base64toMultibaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "base64-multibase [input]",
		Short: `Convert public key from base64 to multibase, according to Ed25519Signature2020 spec.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			base64Str := args[0]
			bytes, err := base64.StdEncoding.DecodeString(base64Str)
			if err != nil {
				return err
			}

			publicKeyMultibaseBytes := []byte{0xed, 0x01}
			publicKeyMultibaseBytes = append(publicKeyMultibaseBytes, bytes...)

			multibaseStr, err := multibase.Encode(multibase.Base58BTC, publicKeyMultibaseBytes)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), multibaseStr)
			return err
		},
	}

	return cmd
}

// base64toBase58Cmd returns cobra Command.
func base64toBase58Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "base64-base58 [input]",
		Short: `Convert public key from base64 to base58, according to Ed25519VerificationKey2018 spec.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			base64Str := args[0]
			bytes, err := base64.StdEncoding.DecodeString(base64Str)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), base58.Encode(bytes))
			return err
		},
	}

	return cmd
}
