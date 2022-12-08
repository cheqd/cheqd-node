package keeper

import (
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/did/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgServer struct {
	Keeper
}

// NewMsgServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewMsgServer(keeper Keeper) types.MsgServer {
	return &MsgServer{Keeper: keeper}
}

var _ types.MsgServer = MsgServer{}

func FindDidDoc(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.DidDocWithMetadata, did string) (res types.DidDocWithMetadata, found bool, err error) {
	// Look in inMemory dict
	value, found := inMemoryDIDs[did]
	if found {
		return value, true, nil
	}

	// Look in state
	if k.HasDidDoc(ctx, did) {
		value, err := k.GetLatestDidDoc(ctx, did)
		if err != nil {
			return types.DidDocWithMetadata{}, false, err
		}

		return value, true, nil
	}

	return types.DidDocWithMetadata{}, false, nil
}

func MustFindDidDoc(k *Keeper, ctx *sdk.Context, inMemoryDIDDocs map[string]types.DidDocWithMetadata, did string) (res types.DidDocWithMetadata, err error) {
	res, found, err := FindDidDoc(k, ctx, inMemoryDIDDocs, did)
	if err != nil {
		return types.DidDocWithMetadata{}, err
	}

	if !found {
		return types.DidDocWithMetadata{}, types.ErrDidDocNotFound.Wrap(did)
	}

	return res, nil
}

func FindVerificationMethod(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.DidDocWithMetadata, didUrl string) (res types.VerificationMethod, found bool, err error) {
	did, _, _, _ := utils.MustSplitDIDUrl(didUrl)

	didDoc, found, err := FindDidDoc(k, ctx, inMemoryDIDs, did)
	if err != nil || !found {
		return types.VerificationMethod{}, found, err
	}

	for _, vm := range didDoc.DidDoc.VerificationMethod {
		if vm.Id == didUrl {
			return *vm, true, nil
		}
	}

	return types.VerificationMethod{}, false, nil
}

func MustFindVerificationMethod(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.DidDocWithMetadata, didUrl string) (res types.VerificationMethod, err error) {
	res, found, err := FindVerificationMethod(k, ctx, inMemoryDIDs, didUrl)
	if err != nil {
		return types.VerificationMethod{}, err
	}

	if !found {
		return types.VerificationMethod{}, types.ErrVerificationMethodNotFound.Wrap(didUrl)
	}

	return res, nil
}

func VerifySignature(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.DidDocWithMetadata, message []byte, signature types.SignInfo) error {
	verificationMethod, err := MustFindVerificationMethod(k, ctx, inMemoryDIDs, signature.VerificationMethodId)
	if err != nil {
		return err
	}

	err = types.VerifySignature(verificationMethod, message, signature.Signature)
	if err != nil {
		return types.ErrInvalidSignature.Wrapf("method id: %s", signature.VerificationMethodId)
	}

	return nil
}

func VerifyAllSignersHaveAllValidSignatures(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.DidDocWithMetadata, message []byte, signers []string, signatures []*types.SignInfo) error {
	for _, signer := range signers {
		signatures := types.FindSignInfosBySigner(signatures, signer)

		if len(signatures) == 0 {
			return types.ErrSignatureNotFound.Wrapf("signer: %s", signer)
		}

		for _, signature := range signatures {
			err := VerifySignature(k, ctx, inMemoryDIDs, message, signature)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// VerifyAllSignersHaveAtLeastOneValidSignature verifies that all signers have at least one valid signature.
// Omit DIDtoBeUpdated and updatedDID if not updating a DID. Otherwise those values will be used to better format error messages.
func VerifyAllSignersHaveAtLeastOneValidSignature(k *Keeper, ctx *sdk.Context, inMemoryDIDs map[string]types.DidDocWithMetadata,
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
