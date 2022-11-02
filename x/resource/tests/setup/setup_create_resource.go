package setup

import (
	"crypto/ed25519"

	"github.com/cheqd/cheqd-node/x/cheqd/tests/setup"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/google/uuid"
)

func (s *TestSetup) CreateResource(payload *types.MsgCreateResourcePayload, signInputs []setup.SignInput) (*types.MsgCreateResourceResponse, error) {
	signBytes := payload.GetSignBytes()
	var signatures []*cheqdtypes.SignInfo

	for _, input := range signInputs {
		signature := ed25519.Sign(input.Key, signBytes)

		signatures = append(signatures, &cheqdtypes.SignInfo{
			VerificationMethodId: input.VerificationMethodId,
			Signature:            signature,
		})
	}

	msg := &types.MsgCreateResource{
		Payload:    payload,
		Signatures: signatures,
	}

	return s.ResourceMsgServer.CreateResource(s.StdCtx, msg)
}

func (s *TestSetup) BuildSimpleResource(collectionId, data, name, _type string) types.MsgCreateResourcePayload {
	return types.MsgCreateResourcePayload{
		Id:           uuid.NewString(),
		CollectionId: collectionId,
		Data:         []byte(data),
		Name:         name,
		ResourceType: _type,
	}
}

func (s *TestSetup) CreateSimpleResource(collectionId, data, name, _type string, signInputs []setup.SignInput) *types.MsgCreateResourceResponse {
	resource := s.BuildSimpleResource(collectionId, data, name, _type)
	res, err := s.CreateResource(&resource, signInputs)
	if err != nil {
		panic(err)
	}

	return res
}
