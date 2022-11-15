package migrations

import (
	"crypto/sha256"
	"encoding/json"
	"reflect"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/mr-tron/base58"
)

type StateValueData interface {
	proto.Message
}

// TODO: Deprecate utils file after exporting did types v1 package as expected.
// They are used in the migration script, but a better solution is to export the types v1 package,
// including utils and load it in the migration v1.0.0 script.
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
		Created:     metadata.Created,
		Updated:     metadata.Updated,
		Deactivated: metadata.Deactivated,
		VersionId:   metadata.VersionId,
		// ToDo: should we make it self-linked?
		NextVersionId:     metadata.VersionId,
		PreviousVersionId: metadata.VersionId,
	}

	return didtypes.NewDidDocWithMetadata(newDidDoc, newMetadata), err
}

func GetVerificationMaterial(vm *didtypesv1.VerificationMethod) string {
	if len(vm.PublicKeyJwk) != 0 {
		jwk := make(map[string]string)
		for _, kv := range vm.PublicKeyJwk {
			jwk[kv.Key] = kv.Value
		}
		res, err := json.Marshal(jwk)
		if err != nil {
			panic(err.Error())
		}

		jwk2020 := didtypes.JsonWebKey2020{
			PublicKeyJwk: res,
		}
		res, err = json.Marshal(jwk2020)
		if err != nil {
			panic(err.Error())
		}

		return string(res)
	}
	pk_multi := didtypes.Ed25519VerificationKey2020{
		PublicKeyMultibase: vm.PublicKeyMultibase,
	}
	res, err := json.Marshal(pk_multi)
	if err != nil {
		panic(err.Error())
	}
	return string(res)
}

func NewDidDocFromV1(didV1 *didtypesv1.Did) *didtypes.DidDoc {
	vms := []*didtypes.VerificationMethod{}
	for _, vm := range didV1.VerificationMethod {
		vms = append(
			vms,
			&didtypes.VerificationMethod{
				Id:                   vm.Id,
				Type:                 vm.Type,
				Controller:           vm.Controller,
				VerificationMaterial: GetVerificationMaterial(vm),
			})
	}
	srvs := []*didtypes.Service{}
	for _, srv := range didV1.Service {
		srvs = append(
			srvs,
			&didtypes.Service{
				Id:              srv.Id,
				Type:            srv.Type,
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

// Make Indy-Style identifiers

func IndyStyleId(id string) string {
	// If id is UUID it should not be changed
	if didutils.IsValidUUID(id) {
		return id
	}

	// Get Hash from current id to make a 32-symbol string
	hash := sha256.Sum256([]byte(id))
	// Indy-style identifier is 16-byte base58 string
	return base58.Encode(hash[:16])
}

func IndyStyleIdList(keys []string) []string {
	if keys == nil {
		return nil
	}
	newKeys := []string{}
	for _, id := range keys {
		newKeys = append(newKeys, IndyStyleId(id))
	}
	return newKeys
}

func IndyStyleDid(did string) string {
	method, namespace, id := didutils.MustSplitDID(did)
	id = IndyStyleId(id)
	return didutils.JoinDID(method, namespace, id)
}

func IndyStyleDidList(didList []string) []string {
	if didList == nil {
		return nil
	}
	newDIDs := []string{}
	for _, did := range didList {
		newDIDs = append(newDIDs, IndyStyleDid(did))
	}
	return newDIDs
}

func IndyStyleDidUrl(didUrl string) string {
	did, path, query, fragment := didutils.MustSplitDIDUrl(didUrl)
	did = IndyStyleDid(did)
	return didutils.JoinDIDUrl(did, path, query, fragment)
}

func IndyStyleDidUrlList(didUrls []string) []string {
	if didUrls == nil {
		return nil
	}
	newDIDUrls := []string{}
	for _, id := range didUrls {
		newDIDUrls = append(newDIDUrls, IndyStyleDidUrl(id))
	}
	return newDIDUrls
}

func MoveToIndyStyleIds(didDoc *didtypes.DidDocWithMetadata) {
	didDoc.DidDoc.Id = IndyStyleId(didDoc.DidDoc.Id)
	for _, vm := range didDoc.DidDoc.VerificationMethod {
		vm.Id = IndyStyleDidUrl(vm.Id)
		vm.Controller = IndyStyleDid(vm.Controller)
	}
	for _, s := range didDoc.DidDoc.Service {
		s.Id = IndyStyleDidUrl(s.Id)
	}

	didDoc.DidDoc.Controller = IndyStyleDidList(didDoc.DidDoc.Controller)
	didDoc.DidDoc.Authentication = IndyStyleDidUrlList(didDoc.DidDoc.Authentication)
	didDoc.DidDoc.AssertionMethod = IndyStyleDidUrlList(didDoc.DidDoc.AssertionMethod)
	didDoc.DidDoc.CapabilityInvocation = IndyStyleDidUrlList(didDoc.DidDoc.CapabilityInvocation)
	didDoc.DidDoc.CapabilityDelegation = IndyStyleDidUrlList(didDoc.DidDoc.CapabilityDelegation)
	didDoc.DidDoc.KeyAgreement = IndyStyleDidUrlList(didDoc.DidDoc.KeyAgreement)
}

func closeIteratorOrPanic(iterator sdk.Iterator) {
	err := iterator.Close()
	if err != nil {
		panic(err.Error())
	}
}
