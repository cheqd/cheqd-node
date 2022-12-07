package integration

import (
	"encoding/json"
	"os"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/x/did/client/cli"
	didtypesv2 "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetypesv2 "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
)

func Loader(path string, payload interface{}) ([]cli.SignInput, error) {
	var tc cli.PayloadWithSignInputs

	file, err := os.ReadFile(path)
	if err != nil {
		return []cli.SignInput{}, err
	}

	err = json.Unmarshal(file, &tc)
	if err != nil {
		return []cli.SignInput{}, err
	}

	switch payload := payload.(type) {
	case *didtypesv1.MsgCreateDidPayload:
		err = integrationhelpers.Codec.UnmarshalJSON(tc.Payload, payload)
	case *didtypesv1.MsgUpdateDidPayload:
		err = integrationhelpers.Codec.UnmarshalJSON(tc.Payload, payload)
	case *didtypesv1.Did:
		err = integrationhelpers.Codec.UnmarshalJSON(tc.Payload, payload)
	case *didtypesv2.MsgCreateDidDocPayload:
		err = integrationhelpers.Codec.UnmarshalJSON(tc.Payload, payload)
	case *didtypesv2.MsgUpdateDidDocPayload:
		err = integrationhelpers.Codec.UnmarshalJSON(tc.Payload, payload)
	case *didtypesv2.DidDoc:
		err = integrationhelpers.Codec.UnmarshalJSON(tc.Payload, payload)
	case *resourcetypesv2.MsgCreateResourcePayload:
		err = integrationhelpers.Codec.UnmarshalJSON(tc.Payload, payload)
	case *resourcetypesv2.Metadata:
		err = integrationhelpers.Codec.UnmarshalJSON(tc.Payload, payload)
	case *resourcetypesv2.ResourceWithMetadata:
		err = integrationhelpers.Codec.UnmarshalJSON(tc.Payload, payload)
	case *resourcetypesv1.MsgCreateResourcePayload:
		err = integrationhelpers.Codec.UnmarshalJSON(tc.Payload, payload)
	case *resourcetypesv1.ResourceHeader:
		err = integrationhelpers.Codec.UnmarshalJSON(tc.Payload, payload)
	default:
		err = json.Unmarshal(tc.Payload, payload)
	}
	if err != nil {
		return []cli.SignInput{}, err
	}
	return tc.SignInputs, err
}
