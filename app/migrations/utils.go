package migrations

import (
	"crypto/sha256"
	"encoding/json"
	"strings"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/mr-tron/base58"
)

type StateValueData interface {
	proto.Message
}

type IteratorKey []byte

func StateValueToDIDDocWithMetadata(stateValue *didtypesv1.StateValue) (didtypes.DidDocWithMetadata, error) {

	var newDidDoc didtypes.DidDoc
	var newMetadata didtypes.Metadata

	didDoc, err := stateValue.UnpackDataAsDid()
	metadata := stateValue.Metadata
	if err != nil {
		return didtypes.DidDocWithMetadata{}, err
	}

	NewDidDocFromV1(didDoc, &newDidDoc)
	newMetadata = didtypes.Metadata{
		Created:     metadata.Created,
		Updated:     metadata.Updated,
		Deactivated: metadata.Deactivated,
		VersionId:   metadata.VersionId,
		// ToDo: should we make it self-linked?
		NextVersionId:     "",
		PreviousVersionId: "",
	}

	return didtypes.DidDocWithMetadata{
		DidDoc:   &newDidDoc,
		Metadata: &newMetadata}, nil
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

func NewDidDocFromV1(didV1 *didtypesv1.Did, newDidDoc *didtypes.DidDoc) {
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
	newDidDoc.Id = didV1.Id
	newDidDoc.VerificationMethod = vms
	newDidDoc.Service = srvs
	newDidDoc.Context = didV1.Context
	newDidDoc.Controller = didV1.Controller
	newDidDoc.Authentication = didV1.Authentication
	newDidDoc.AssertionMethod = didV1.AssertionMethod
	newDidDoc.CapabilityDelegation = didV1.CapabilityDelegation
	newDidDoc.CapabilityInvocation = didV1.CapabilityInvocation
	newDidDoc.KeyAgreement = didV1.KeyAgreement
	newDidDoc.AlsoKnownAs = didV1.AlsoKnownAs
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

func CollectAllKeys(
	ctx sdk.Context,
	storeKey *types.KVStoreKey,
	iteratorPrefixKey []byte) []IteratorKey {

	keys := []IteratorKey{}
	store := ctx.KVStore(storeKey)

	iterator := sdk.KVStorePrefixIterator(store, iteratorPrefixKey)
	closeIteratorOrPanic(iterator)

	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, IteratorKey(iterator.Key()))
	}
	return keys
}

func ResourceV1HeaderkeyToDataKey(headerKey []byte) []byte {
	return []byte(
		strings.Replace(
			string(headerKey),
			string(resourcetypesv1.ResourceHeaderKey),
			string(resourcetypesv1.ResourceDataKey),
			1))
}

func ResourceV2MetadataKeyToDataKey(metadataKey []byte) []byte {
	return []byte(
		strings.Replace(
			string(metadataKey),
			string(resourcetypes.ResourceMetadataKey),
			string(resourcetypes.ResourceDataKey),
			1))
}

func closeIteratorOrPanic(iterator sdk.Iterator) {
	err := iterator.Close()
	if err != nil {
		panic(err.Error())
	}
}
