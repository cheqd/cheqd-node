package setup

import (
	"crypto/ed25519"

	"github.com/cheqd/cheqd-node/x/did/types"
)

func (s *TestSetup) DeactivateDidDoc(payload *types.MsgDeactivateDidDocPayload, signInputs []SignInput) (*types.MsgDeactivateDidDocResponse, error) {
	signBytes := payload.GetSignBytes()
	signatures := make([]*types.SignInfo, 0, len(signInputs))

	for _, input := range signInputs {
		signature := ed25519.Sign(input.Key, signBytes)

		signatures = append(signatures, &types.SignInfo{
			VerificationMethodId: input.VerificationMethodID,
			Signature:            signature,
		})
	}

	msg := &types.MsgDeactivateDidDoc{
		Payload:    payload,
		Signatures: signatures,
	}

	return s.MsgServer.DeactivateDidDoc(s.StdCtx, msg)
}
