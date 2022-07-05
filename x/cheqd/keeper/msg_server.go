package keeper

import (
	"encoding/base64"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func FindDid(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.StateValue, did string) (res types.StateValue, found bool, err error) {
	// Look in inMemory dict
	value, found := inMemoryDIDs[did]
	if found {
		return value, true, nil
	}

	// Look in state
	if k.HasDid(ctx, did) {
		value, err := k.GetDid(ctx, did)
		if err != nil {
			return types.StateValue{}, false, err
		}

		return value, true, nil
	}

	return types.StateValue{}, false, nil
}

func MustFindDid(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.StateValue, did string) (res types.StateValue, err error) {
	res, found, err := FindDid(k, ctx, inMemoryDIDs, did)
	if err != nil {
		return types.StateValue{}, err
	}

	if !found {
		return types.StateValue{}, types.ErrDidDocNotFound.Wrap(did)
	}

	return res, nil
}

func FindVerificationMethod(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.StateValue, didUrl string) (res types.VerificationMethod, found bool, err error) {
	did, _, _, _ := utils.MustSplitDIDUrl(didUrl)

	stateValue, found, err := FindDid(k, ctx, inMemoryDIDs, did)
	if err != nil || !found {
		return types.VerificationMethod{}, found, err
	}

	didDoc, err := stateValue.UnpackDataAsDid()
	if err != nil {
		return types.VerificationMethod{}, false, err
	}

	for _, vm := range didDoc.VerificationMethod {
		if vm.Id == didUrl {
			return *vm, true, nil
		}
	}

	return types.VerificationMethod{}, false, nil
}

func MustFindVerificationMethod(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.StateValue, didUrl string) (res types.VerificationMethod, err error) {
	res, found, err := FindVerificationMethod(k, ctx, inMemoryDIDs, didUrl)
	if err != nil {
		return types.VerificationMethod{}, err
	}

	if !found {
		return types.VerificationMethod{}, types.ErrVerificationMethodNotFound.Wrap(didUrl)
	}

	return res, nil
}

func VerifySignature(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.StateValue, message []byte, signature types.SignInfo) error {
	verificationMethod, err := MustFindVerificationMethod(k, ctx, inMemoryDIDs, signature.VerificationMethodId)
	if err != nil {
		return err
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signature.Signature)
	if err != nil {
		return err
	}

	err = types.VerifySignature(verificationMethod, message, signatureBytes)
	if err != nil {
		return types.ErrInvalidSignature.Wrapf("method id: %s", signature.VerificationMethodId)
	}

	return nil
}

func VerifyAllSignersHaveAllValidSignatures(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.StateValue, message []byte, signers []string, signatures []*types.SignInfo) error {
	for _, signer := range signers {
		signature, found := types.FindSignInfoBySigner(signatures, signer)

		if !found {
			return types.ErrSignatureNotFound.Wrapf("signer: %s", signer)
		}

		err := VerifySignature(&k.Keeper, &ctx, inMemoryDids, message, signature)
		if err != nil {
			return err
		}
	}
	return nil
}

// VerifyAllSignersHaveAtLeastOneValidSignature verifies that all signers have at least one valid signature.
// Omit DIDtoBeUpdated and updatedDID if not updating a DID. Otherwise those values will be used to better format error messages.
func VerifyAllSignersHaveAtLeastOneValidSignature(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.StateValue,
	message []byte, signers []string, signatures []*types.SignInfo, DIDToBeUpdated string, updatedDID string,
) error {
	for _, signer := range signers {
		signaturesBySigner := types.FindSignInfosBySigner(signatures, signer)
		signerForErrorMessage := GetSignerIdForErrorMessage(signer, DIDToBeUpdated, updatedDID)

		if len(signaturesBySigner) == 0 {
			return types.ErrSignatureNotFound.Wrapf("there should be at least one signature by %s", signerForErrorMessage)
		}

		found := false
		for _, signature := range signaturesBySigner {
			err := VerifySignature(k, ctx, inMemoryDIDs, message, signature)
			if err == nil {
				found = true
				break
			}
		}

		if !found {
			return types.ErrInvalidSignature.Wrapf("there should be at least one valid signature by %s", signerForErrorMessage)
		}
	}

	return nil
}
