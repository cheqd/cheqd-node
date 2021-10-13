package keeper

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) Verify(context sdk.Context, request *types.MsgWriteRequest) bool {
	metadataBytes, _ := json.Marshal(&request.Metadata)
	dataBytes := request.Data.Value

	signingInput := ([]byte)(base64.StdEncoding.EncodeToString(metadataBytes) + base64.StdEncoding.EncodeToString(dataBytes))

	result := true

	for didUrl := range request.Signatures {
		pubKey, err := k.FindPublicKey(context, request, didUrl)
		if err != nil {
			errMsg := fmt.Sprintf("Cannot get public key. Cause: %T", err)
			k.Logger(context).Error(errMsg)
			return false
		}

		signature, _ := base64.StdEncoding.DecodeString(request.Signatures[didUrl])
		result = result && ed25519.Verify(pubKey, signingInput, signature)
	}

	return result
}

func (k Keeper) FindPublicKey(context sdk.Context, request *types.MsgWriteRequest, didUrl string) (ed25519.PublicKey, error) {
	did, _ := utils.SplitDidUrlIntoDidAndFragment(didUrl)

	var authentication []string
	var verificationMethod []*types.VerificationMethod

	if request.Data.TypeUrl == "/cheqdid.cheqdnode.cheqd.MsgCreateDid" {
		didDoc, isMsgIdentity := request.Data.GetCachedValue().(*types.MsgCreateDid)

		if !isMsgIdentity {
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, didDoc)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}

		authentication = didDoc.Authentication
		verificationMethod = didDoc.VerificationMethod
	} else {
		didDoc, _ := k.GetDid(context, did)
		authentication = didDoc.Authentication
		verificationMethod = didDoc.VerificationMethod
	}

	return utils.FindPublicKey(authentication, verificationMethod, didUrl)
}
