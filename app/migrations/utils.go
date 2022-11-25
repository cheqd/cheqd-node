package migrations

import (
	"crypto/sha256"
	"strings"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"

	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/mr-tron/base58"
)

type StateValueData interface {
	proto.Message
}

type ByteStr []byte

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
	didDoc.DidDoc.Id = IndyStyleDid(didDoc.DidDoc.Id)
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

func ReadAllKeys(store types.KVStore, prefix []byte) []ByteStr {
	keys := []ByteStr{}

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer closeIteratorOrPanic(iterator)

	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, ByteStr(iterator.Key()))
	}

	return keys
}

func closeIteratorOrPanic(iterator sdk.Iterator) {
	err := iterator.Close()
	if err != nil {
		panic(err.Error())
	}
}

func ResourceV2MetadataKeyToDataKey(metadataKey []byte) []byte {
	return []byte(
		strings.Replace(
			string(metadataKey),
			string(resourcetypes.ResourceMetadataKey),
			string(resourcetypes.ResourceDataKey),
			1))
}
