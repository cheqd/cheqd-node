package setup

import (
	"context"
	"encoding/base64"

	didtypesv2 "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgServerV1 struct {
	Keeper KeeperV1
}

func NewMsgServerV1(keeper KeeperV1) *MsgServerV1 {
	return &MsgServerV1{Keeper: keeper}
}

func (k MsgServerV1) CreateDidDocV1(goCtx context.Context, msg *didtypesv1.MsgCreateDid) (*didtypesv1.MsgCreateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate DID doesn't exist
	if k.Keeper.HasDid(&ctx, msg.Payload.Id) {
		return nil, didtypesv2.ErrDidDocExists.Wrap(msg.Payload.Id)
	}

	// Validate namespaces
	namespace := k.Keeper.GetDidNamespace(&ctx)
	err := msg.Validate([]string{namespace})
	if err != nil {
		return nil, didtypesv2.ErrNamespaceValidation.Wrap(err.Error())
	}

	// Build metadata and stateValue
	did := msg.Payload.ToDid()
	metadata := didtypesv1.NewMetadataFromContext(ctx)
	stateValue, err := didtypesv1.NewStateValue(&did, &metadata)
	if err != nil {
		return nil, err
	}

	// Consider did that we are going to create during did resolutions
	inMemoryDids := map[string]didtypesv1.StateValue{did.Id: stateValue}

	// Check controllers' existence
	controllers := did.AllControllerDids()
	for _, controller := range controllers {
		_, err := MustFindDid(&k.Keeper, &ctx, inMemoryDids, controller)
		if err != nil {
			return nil, err
		}
	}

	// Verify signatures
	signers := GetSignerDIDsForDIDCreation(did)
	err = VerifyAllSignersHaveAllValidSignatures(&k.Keeper, &ctx, inMemoryDids, msg.Payload.GetSignBytes(), signers, msg.Signatures)
	if err != nil {
		return nil, err
	}

	// Build state value
	value, err := didtypesv1.NewStateValue(&did, &metadata)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.SetDid(&ctx, &value)
	if err != nil {
		return nil, didtypesv2.ErrInternal.Wrapf(err.Error())
	}

	// Build and return response
	return &didtypesv1.MsgCreateDidResponse{
		Id: did.Id,
	}, nil
}

func FindDid(k *KeeperV1, ctx *sdk.Context, inMemoryDIDs map[string]didtypesv1.StateValue, did string) (res didtypesv1.StateValue, found bool, err error) {
	// Look in inMemory dict
	value, found := inMemoryDIDs[did]
	if found {
		return value, true, nil
	}

	// Look in state
	if k.HasDid(ctx, did) {
		value, err := k.GetDid(ctx, did)
		if err != nil {
			return didtypesv1.StateValue{}, false, err
		}

		return value, true, nil
	}

	return didtypesv1.StateValue{}, false, nil
}

func MustFindDid(k *KeeperV1, ctx *sdk.Context, inMemoryDIDs map[string]didtypesv1.StateValue, did string) (res didtypesv1.StateValue, err error) {
	res, found, err := FindDid(k, ctx, inMemoryDIDs, did)
	if err != nil {
		return didtypesv1.StateValue{}, err
	}

	if !found {
		return didtypesv1.StateValue{}, didtypesv2.ErrDidDocNotFound.Wrap(did)
	}

	return res, nil
}

func GetSignerDIDsForDIDCreation(did didtypesv1.Did) []string {
	res := did.GetControllersOrSubject()
	res = append(res, did.GetVerificationMethodControllers()...)

	return didutils.UniqueSorted(res)
}

func FindVerificationMethod(k *KeeperV1, ctx *sdk.Context, inMemoryDIDs map[string]didtypesv1.StateValue, didUrl string) (res didtypesv1.VerificationMethod, found bool, err error) {
	did, _, _, _ := didutils.MustSplitDIDUrl(didUrl)

	stateValue, found, err := FindDid(k, ctx, inMemoryDIDs, did)
	if err != nil || !found {
		return didtypesv1.VerificationMethod{}, found, err
	}

	didDoc, err := stateValue.UnpackDataAsDid()
	if err != nil {
		return didtypesv1.VerificationMethod{}, false, err
	}

	for _, vm := range didDoc.VerificationMethod {
		if vm.Id == didUrl {
			return *vm, true, nil
		}
	}

	return didtypesv1.VerificationMethod{}, false, nil
}

func MustFindVerificationMethod(k *KeeperV1, ctx *sdk.Context, inMemoryDIDs map[string]didtypesv1.StateValue, didUrl string) (res didtypesv1.VerificationMethod, err error) {
	res, found, err := FindVerificationMethod(k, ctx, inMemoryDIDs, didUrl)
	if err != nil {
		return didtypesv1.VerificationMethod{}, err
	}

	if !found {
		return didtypesv1.VerificationMethod{}, didtypesv2.ErrVerificationMethodNotFound.Wrap(didUrl)
	}

	return res, nil
}

func VerifySignature(k *KeeperV1, ctx *sdk.Context, inMemoryDIDs map[string]didtypesv1.StateValue, message []byte, signature didtypesv1.SignInfo) error {
	verificationMethod, err := MustFindVerificationMethod(k, ctx, inMemoryDIDs, signature.VerificationMethodId)
	if err != nil {
		return err
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signature.Signature)
	if err != nil {
		return err
	}

	err = didtypesv1.VerifySignature(verificationMethod, message, signatureBytes)
	if err != nil {
		return didtypesv2.ErrInvalidSignature.Wrapf("method id: %s", signature.VerificationMethodId)
	}

	return nil
}

func VerifyAllSignersHaveAllValidSignatures(k *KeeperV1, ctx *sdk.Context, inMemoryDIDs map[string]didtypesv1.StateValue, message []byte, signers []string, signatures []*didtypesv1.SignInfo) error {
	for _, signer := range signers {
		signatures := didtypesv1.FindSignInfosBySigner(signatures, signer)

		if len(signatures) == 0 {
			return didtypesv2.ErrSignatureNotFound.Wrapf("signer: %s", signer)
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
