package setup

import (
	didtestssetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/resource/types/v1"
)

func (s *ExtendedTestSetup) CreateResourceV1(payload *v1.MsgCreateResourcePayload, signInputs []didtestssetup.SignInput) (*v1.MsgCreateResourceResponse, error) {
	// TODO: Implement this method
	return &v1.MsgCreateResourceResponse{}, nil
}
