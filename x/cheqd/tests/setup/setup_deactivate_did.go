package setup

import (
	"crypto/ed25519"
	"encoding/base64"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	// "github.com/cheqd/cheqd-node/x/cheqd/utils"
)


// func (s *TestSetup) DeactivateDID(
// 	payload *types.MsgDeactivateDidPayload, 
// 	signInputs []SignInput) *types.MsgDeactivateDid {

// 	signBytes := payload.GetSignBytes()
// 	var signatures []*types.SignInfo

// 	for _, input := range signInputs {
// 		signature := ed25519.Sign(input.Key, signBytes)

// 		signatures = append(signatures, &types.SignInfo{
// 			VerificationMethodId: input.VerificationMethodId,
// 			Signature:            base64.StdEncoding.EncodeToString(signature),
// 		})
// 	}

// 	msg := &types.MsgDeactivateDid{
// 		Payload:    payload,
// 		Signatures: signatures,
// 	}
	
// }

func (s *TestSetup) DeactivateDid(payload *types.MsgDeactivateDidPayload, signInputs []SignInput) (*types.MsgDeactivateDidResponse, error) {
	signBytes := payload.GetSignBytes()
	var signatures []*types.SignInfo

	for _, input := range signInputs {
		signature := ed25519.Sign(input.Key, signBytes)

		signatures = append(signatures, &types.SignInfo{
			VerificationMethodId: input.VerificationMethodId,
			Signature:            base64.StdEncoding.EncodeToString(signature),
		})
	}

	msg := &types.MsgDeactivateDid{
		Payload:    payload,
		Signatures: signatures,
	}

	return s.MsgServer.DeactivateDid(s.StdCtx, msg)
}

// func (s *TestSetup) SendDeactivateDid(msg *types.MsgDeactivateDidPayload, keys map[string]ed25519.PrivateKey) (*types.Metadata, error) {
// 	_, err := s.Handler(s.Ctx, s.DeactivateDID(msg, keys))
// 	if err != nil {
// 		return nil, err
// 	}

// 	updated, _ := s.MsgServer.GetDid(&s.Ctx, msg.Id)
// 	return updated.Metadata, nil
// }


