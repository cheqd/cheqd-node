package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
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

func ResolveDid(k *Keeper, ctx *sdk.Context, did string, inMemoryDIDs map[string]types.StateValue) (types.StateValue, error) {
	value, found := inMemoryDIDs[did]
	if found {
		return value, nil
	}

	stateValue, err := k.GetDid(ctx, did)
	if err != nil {
		return types.StateValue{}, err
	}

	return stateValue, nil
}

func ValidateDIDHasSupportedAuthKey() {

}


//func ValidateVerificationMethod(namespace string, vm *VerificationMethod) error {
//	switch utils.GetVerificationMethodType(vm.Type) {
//	case utils.PublicKeyJwk:
//		if len(vm.PublicKeyJwk) == 0 {
//			return ErrBadRequest.Wrapf("%s: should contain `PublicKeyJwk` verification material property", vm.Type)
//		}
//	case utils.PublicKeyMultibase:
//		if len(vm.PublicKeyMultibase) == 0 {
//			return ErrBadRequest.Wrapf("%s: should contain `PublicKeyMultibase` verification material property", vm.Type)
//		}
//	default:
//		return ErrBadRequest.Wrapf("%s: unsupported verification method type", vm.Type)
//	}
//
//	if len(vm.PublicKeyMultibase) == 0 && vm.PublicKeyJwk == nil {
//		return ErrBadRequest.Wrap("The verification method must contain either a PublicKeyMultibase or a PublicKeyJwk")
//	}
//
//	if len(vm.Controller) == 0 {
//		return ErrBadRequestIsRequired.Wrap("Controller")
//	}
//
//	return nil
//}



//func AppendSignerIfNeed(signers []types.Signer, controller string, msg *types.MsgUpdateDidPayload) []types.Signer {
//	for _, signer := range signers {
//		if signer.Signer == controller {
//			return signers
//		}
//	}
//
//	signer := types.Signer{
//		Signer: controller,
//	}
//
//	if controller == msg.Id {
//		signer.VerificationMethod = msg.VerificationMethod
//		signer.Authentication = msg.Authentication
//	}
//
//	return append(signers, signer)
//}
//
//func (k msgServer) ValidateDidControllers(ctx *sdk.Context, id string, controllers []string, verMethods []*types.VerificationMethod) error {
//
//	for _, verificationMethod := range verMethods {
//		if err := k.ValidateController(ctx, id, verificationMethod.Controller); err != nil {
//			return err
//		}
//	}
//
//	for _, didController := range controllers {
//		if err := k.ValidateController(ctx, id, didController); err != nil {
//			return err
//		}
//	}
//	return nil
//}
