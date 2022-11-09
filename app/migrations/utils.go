package migrations

import (
	"encoding/json"
	"reflect"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

type StateValueData interface {
	proto.Message
}

type PubKeyMultibase struct {
	publicKeyMultibase string
}

func UnpackData(stateValue didtypesv1.StateValue) (StateValueData, error) {
	value, isOk := stateValue.Data.GetCachedValue().(StateValueData)
	if !isOk {
		return nil, didtypes.ErrUnpackStateValue.Wrapf("invalid type url: %s", stateValue.Data.TypeUrl)
	}

	return value, nil
}

func UnpackDataAsDid(stateValue didtypesv1.StateValue) (*didtypesv1.Did, error) {
	data, err := UnpackData(stateValue)
	if err != nil {
		return nil, err
	}

	value, isValue := data.(*didtypesv1.Did)
	if !isValue {
		return nil, didtypes.ErrUnpackStateValue.Wrap(reflect.TypeOf(data).String())
	}

	return value, nil
}

func StateValueToDIDDocWithMetadata(stateValue didtypesv1.StateValue) (didtypes.DidDocWithMetadata, error) {
	var err error
	didDoc, err := UnpackDataAsDid(stateValue)
	metadata := stateValue.Metadata
	if err != nil {
		return didtypes.DidDocWithMetadata{}, err
	}
	newDidDoc := NewDidDocFromV1(didDoc)
	newMetadata := &didtypes.Metadata{
		Created: metadata.Created,
		Updated: metadata.Updated,
		Deactivated: metadata.Deactivated,
		VersionId: metadata.VersionId,
		// ToDo: should we make it self-linked?
		NextVersionId: metadata.VersionId,
		PreviousVersionId: metadata.VersionId,
	}

	return didtypes.NewDidDocWithMetadata(newDidDoc, newMetadata), err
}

func GetVerificationMaterial(vm *didtypesv1.VerificationMethod) string{
	if len(vm.PublicKeyJwk) != 0 {
		jwk := make(map[string]string)
		for _, kv := range(vm.PublicKeyJwk) {
			jwk[kv.Key] = kv.Value
		}
		res, _ := json.Marshal(jwk)
		return string(res)
	}
	pk_multi := PubKeyMultibase{
		publicKeyMultibase: vm.PublicKeyMultibase,
	}
	res, _ := json.Marshal(pk_multi)
	return string(res)
}

func NewDidDocFromV1(didV1 *didtypesv1.Did) *didtypes.DidDoc {
	vms := []*didtypes.VerificationMethod{}
	for _, vm := range didV1.VerificationMethod { 
		vms = append(
			vms,
			&didtypes.VerificationMethod{
				Id: vm.Id,
				Type: vm.Type,
				Controller: vm.Controller,
				VerificationMaterial: GetVerificationMaterial(vm),
			})
	}
	srvs := []*didtypes.Service{}
	for _, srv := range didV1.Service {
		srvs = append(
			srvs,
			&didtypes.Service{
				Id: srv.Id,
				Type: srv.Type,
				ServiceEndpoint: []string{srv.ServiceEndpoint},
			})
	}
	return &didtypes.DidDoc{
		Context:              didV1.Context,
		Id:                   didV1.Id,
		Controller:           didV1.Controller,
		VerificationMethod:   vms,
		Authentication:       didV1.Authentication,
		AssertionMethod:      didV1.AssertionMethod,
		CapabilityInvocation: didV1.CapabilityInvocation,
		CapabilityDelegation: didV1.CapabilityDelegation,
		KeyAgreement:         didV1.KeyAgreement,
		Service:              srvs,
		AlsoKnownAs:          didV1.AlsoKnownAs,
	}
}


func closeIteratorOrPanic(iterator sdk.Iterator) {
	err := iterator.Close()
	if err != nil {
		panic(err.Error())
	}
}