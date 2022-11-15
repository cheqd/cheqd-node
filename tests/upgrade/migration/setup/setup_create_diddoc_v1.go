package setup

import (
	didtestssetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	v1 "github.com/cheqd/cheqd-node/x/did/types/v1"
)

func (s *ExtendedTestSetup) CreateDidV1(payload *v1.MsgCreateDidPayload, signInputs []didtestssetup.SignInput) (*v1.MsgCreateDidResponse, error) {
	// TODO: Implement this method
	return &v1.MsgCreateDidResponse{}, nil
}
