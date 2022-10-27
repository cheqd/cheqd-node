package setup

import (
	"crypto/ed25519"
	"encoding/base64"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

func (s *TestSetup) UpdateDid(payload *types.MsgUpdateDidPayload, signInputs []SignInput) (*types.MsgUpdateDidResponse, error) {
	signBytes := payload.GetSignBytes()
	var signatures []*types.SignInfo

	for _, input := range signInputs {
		signature := ed25519.Sign(input.Key, signBytes)

		signatures = append(signatures, &types.SignInfo{
			VerificationMethodId: input.VerificationMethodId,
			Signature:            base64.StdEncoding.EncodeToString(signature),
		})
	}

	msg := &types.MsgUpdateDid{
		Payload:    payload,
		Signatures: signatures,
	}

	return s.MsgServer.UpdateDid(s.StdCtx, msg)
}
