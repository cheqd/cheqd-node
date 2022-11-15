package setup

import (
	didtestssetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
)

type MinimalDidDocInfoV1 struct {
	Msg       *didtypesv1.MsgCreateDidPayload
	SignInput didtestssetup.SignInput
}
