package setup

import (
	"crypto/ed25519"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

func (s *TestSetup) DeactivateDid(payload *types.MsgDeactivateDidDocPayload, signInputs []SignInput) (*types.MsgDeactivateDidDocResponse, error) {
	signBytes := payload.GetSignBytes()
	var signatures []*types.SignInfo

	for _, input := range signInputs {
		signature := ed25519.Sign(input.Key, signBytes)

		signatures = append(signatures, &types.SignInfo{
			VerificationMethodId: input.VerificationMethodId,
			Signature:            signature,
		})
	}

	msg := &types.MsgDeactivateDidDoc{
		Payload:    payload,
		Signatures: signatures,
	}

	return s.MsgServer.DeactivateDidDoc(s.StdCtx, msg)
}
