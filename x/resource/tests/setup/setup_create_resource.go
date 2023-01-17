package setup

import (
	"crypto/ed25519"

	"github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/google/uuid"
)

func (s *TestSetup) CreateResource(payload *types.MsgCreateResourcePayload, signInputs []setup.SignInput) (*types.MsgCreateResourceResponse, error) {
	signBytes := payload.GetSignBytes()
	signatures := make([]*didtypes.SignInfo, 0, len(signInputs))

	for _, input := range signInputs {
		signature := ed25519.Sign(input.Key, signBytes)

		signatures = append(signatures, &didtypes.SignInfo{
			VerificationMethodId: input.VerificationMethodID,
			Signature:            signature,
		})
	}

	msg := &types.MsgCreateResource{
		Payload:    payload,
		Signatures: signatures,
	}

	return s.ResourceMsgServer.CreateResource(s.StdCtx, msg)
}

func (s *TestSetup) BuildSimpleResource(collectionID, data, name, _type string) types.MsgCreateResourcePayload {
	return types.MsgCreateResourcePayload{
		Id:           uuid.NewString(),
		CollectionId: collectionID,
		Data:         []byte(data),
		Name:         name,
		ResourceType: _type,
	}
}

func (s *TestSetup) CreateSimpleResource(collectionID, data, name, _type string, signInputs []setup.SignInput) *types.MsgCreateResourceResponse {
	resource := s.BuildSimpleResource(collectionID, data, name, _type)
	res, err := s.CreateResource(&resource, signInputs)
	if err != nil {
		panic(err)
	}

	return res
}
