package cli

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"os"

	errorsmod "cosmossdk.io/errors"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

const (
	FlagVersionID = "version-id"
)

type DIDDocument struct {
	Context              []string             `json:"context"`
	ID                   string               `json:"id"`
	Controller           []string             `json:"controller,omitempty"`
	VerificationMethod   []VerificationMethod `json:"verificationMethod,omitempty"`
	Authentication       []string             `json:"authentication,omitempty"`
	AssertionMethod      []string             `json:"assertionMethod,omitempty"`
	CapabilityInvocation []string             `json:"capabilityInvocation,omitempty"`
	CapabilityDelegation []string             `json:"capabilityDelegation,omitempty"`
	KeyAgreement         []string             `json:"keyAgreement,omitempty"`
	Service              []Service            `json:"service,omitempty"`
	AlsoKnownAs          []string             `json:"alsoKnownAs,omitempty"`
}

type VerificationMethod map[string]any

type Service struct {
	ID              string   `json:"id"`
	Type            string   `json:"type"`
	ServiceEndpoint []string `json:"serviceEndpoint"`
}

type PayloadWithSignInputs struct {
	Payload    json.RawMessage
	SignInputs []SignInput
}

type SignInput struct {
	VerificationMethodID string
	PrivKey              ed25519.PrivateKey
}

// AddTxFlagsToCmd adds common flags to a module tx command.
func AddTxFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "json", "Output format (text|json)")
	cmd.Flags().String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	cmd.Flags().String(flags.FlagFrom, "", "Name or address of private key with which to sign")
	cmd.Flags().Uint64P(flags.FlagAccountNumber, "a", 0, "The account number of the signing account (offline mode only)")
	cmd.Flags().Uint64P(flags.FlagSequence, "s", 0, "The sequence number of the signing account (offline mode only)")
	cmd.Flags().String(flags.FlagNote, "", "Note to add a description to the transaction (previously --memo)")
	cmd.Flags().String(flags.FlagGasPrices, "", "Gas prices in decimal format to determine the transaction fee (e.g. 50ncheq)")
	cmd.Flags().String(flags.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	cmd.Flags().Bool(flags.FlagUseLedger, false, "Use a connected Ledger device")
	cmd.Flags().StringP(flags.FlagBroadcastMode, "b", flags.BroadcastSync, "Transaction broadcasting mode (sync|async|block)")
	cmd.Flags().Bool(flags.FlagDryRun, false, "ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)")
	cmd.Flags().Bool(flags.FlagGenerateOnly, false, "Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)")
	cmd.Flags().Bool(flags.FlagOffline, false, "Offline mode (does not allow any online functionality)")
	cmd.Flags().BoolP(flags.FlagSkipConfirmation, "y", false, "Skip tx broadcasting prompt confirmation")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test|memory)")
	cmd.Flags().String(flags.FlagSignMode, "", "Choose sign mode (direct|amino-json|direct-aux), this is an advanced feature")
	cmd.Flags().Uint64(flags.FlagTimeoutHeight, 0, "Set a block timeout height to prevent the tx from being committed past a certain height")
	cmd.Flags().String(flags.FlagFeePayer, "", "Fee payer pays fees for the transaction instead of deducting from the signer")
	cmd.Flags().String(flags.FlagFeeGranter, "", "Fee granter grants fees for the transaction")
	cmd.Flags().String(flags.FlagTip, "", "Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator")
	cmd.Flags().Bool(flags.FlagAux, false, "Generate aux signer data instead of sending a tx")

	// overrides
	cmd.Flags().String(flags.FlagGas, flags.GasFlagAuto, fmt.Sprintf("Gas limit to set per-transaction; set to %q to calculate sufficient gas automatically", flags.GasFlagAuto))
	cmd.Flags().Float64(flags.FlagGasAdjustment, 1.8, "adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored ")

	// flags --fees added by each module's tx command
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
	signatures := make([]*types.SignInfo, 0, len(signInputs))

	for _, signInput := range signInputs {
		signatureBytes := ed25519.Sign(signInput.PrivKey, signBytes)

		signInfo := types.SignInfo{
			VerificationMethodId: signInput.VerificationMethodID,
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

	if !errorsmod.IsOf(err, sdkerrors.ErrIO, sdkerrors.ErrKeyNotFound) {
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

func GetFromSpecCompliantPayload(specPayload DIDDocument) ([]*types.VerificationMethod, []*types.Service, error) {
	verificationMethod := make([]*types.VerificationMethod, 0, len(specPayload.VerificationMethod))
	for i, vm := range specPayload.VerificationMethod {
		var verificationMethodType string
		if value, ok := vm["type"].(string); !ok {
			return nil, nil, fmt.Errorf("%d: verification method type is not specified", i)
		} else {
			verificationMethodType = value
		}

		switch verificationMethodType {
		case "Ed25519VerificationKey2020":
			_, ok := vm["publicKeyMultibase"]
			if !ok {
				return nil, nil, fmt.Errorf("%d: publicKeyMultibase is not specified", i)
			}

			verificationMethod = append(verificationMethod, &types.VerificationMethod{
				Id:                     vm["id"].(string),
				VerificationMethodType: vm["type"].(string),
				Controller:             vm["controller"].(string),
				VerificationMaterial:   vm["publicKeyMultibase"].(string),
			})
		case "Ed25519VerificationKey2018":
			_, ok := vm["publicKeyBase58"]
			if !ok {
				return nil, nil, fmt.Errorf("%d: publicKeyBase58 is not specified", i)
			}

			verificationMethod = append(verificationMethod, &types.VerificationMethod{
				Id:                     vm["id"].(string),
				VerificationMethodType: vm["type"].(string),
				Controller:             vm["controller"].(string),
				VerificationMaterial:   vm["publicKeyBase58"].(string),
			})
		case "JsonWebKey2020":
			_, ok := vm["publicKeyJwk"]
			if !ok {
				return nil, nil, fmt.Errorf("%d: publicKeyJwk is not specified", i)
			}

			jwk, err := json.Marshal(vm["publicKeyJwk"])
			if err != nil {
				return nil, nil, err
			}

			verificationMethod = append(verificationMethod, &types.VerificationMethod{
				Id:                     vm["id"].(string),
				VerificationMethodType: vm["type"].(string),
				Controller:             vm["controller"].(string),
				VerificationMaterial:   string(jwk),
			})
		default:
			return nil, nil, fmt.Errorf("%d: verification method type is not supported", i)
		}
	}

	service := make([]*types.Service, 0, len(specPayload.Service))
	for _, s := range specPayload.Service {
		service = append(service, &types.Service{
			Id:              s.ID,
			ServiceType:     s.Type,
			ServiceEndpoint: s.ServiceEndpoint,
		})
	}

	return verificationMethod, service, nil
}
