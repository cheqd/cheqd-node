package keeper

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DidDoc = struct {
	Authentication     []string
	VerificationMethod []*types.VerificationMethod
}

func (k Keeper) Verify(context sdk.Context, request *types.MsgWriteRequest) bool {
	metadataBytes, _ := json.Marshal(&request.Metadata)
	dataBytes := request.Data.Value

	signingInput := base64.StdEncoding.EncodeToString(metadataBytes) + "." + base64.StdEncoding.EncodeToString(dataBytes)
	signingInputBytes, _ := base64.StdEncoding.DecodeString(signingInput)

	result := true

	for didUrl := range request.Signatures {
		did, _ := utils.SplitDidUrlIntoDidAndFragment(didUrl)
		didDoc := k.GetDid(context, did)

		pubKey := FindPublicKey(didDoc.Authentication, didDoc.VerificationMethod, didUrl)
		signature, _ := base64.StdEncoding.DecodeString(request.Signatures[didUrl])
		result = result && ed25519.Verify(pubKey, signingInputBytes, signature)
	}

	return result
}

func (k Keeper) VerifyMsgCreateDid(request *types.MsgWriteRequest, didDoc *types.MsgCreateDid) bool {
	metadataBytes, _ := json.Marshal(&request.Metadata)
	dataBytes := request.Data.Value

	signingInput := base64.StdEncoding.EncodeToString(metadataBytes) + "." + base64.StdEncoding.EncodeToString(dataBytes)
	signingInputBytes := []byte(base64.StdEncoding.EncodeToString([]byte(signingInput)))

	result := true

	for did := range request.Signatures {
		pubKey := FindPublicKey(didDoc.Authentication, didDoc.VerificationMethod, did)
		signature, _ := base64.StdEncoding.DecodeString(request.Signatures[did])
		result = result && ed25519.Verify(pubKey, signingInputBytes, signature)
	}

	return result
}

// FindAuthentication improve error handling
func FindPublicKey(authentication []string, verificationMethods []*types.VerificationMethod, id string) ed25519.PublicKey {
	for _, authentication := range authentication {
		if authentication == id {
			for _, vm := range verificationMethods {
				if vm.Id == id {
					return base58.Decode(vm.PublicKeyMultibase[1:])
				}
			}
		}
	}

	return nil
}
