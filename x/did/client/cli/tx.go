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

type VerificationMethod map[string]interface{}

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
