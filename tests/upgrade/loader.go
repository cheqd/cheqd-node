//go:build upgrade

package upgrade

import (
	"encoding/json"
	"os"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	didtypesv2 "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetypesv2 "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
)

func Loader(path string, msg any) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	switch msg := msg.(type) {
	case *didtypesv1.MsgCreateDidPayload:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	case *didtypesv1.MsgUpdateDidPayload:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	case *didtypesv1.Did:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	case *didtypesv2.MsgCreateDidDocPayload:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	case *didtypesv2.MsgUpdateDidDocPayload:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	case *didtypesv2.DidDocWithMetadata:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	case *resourcetypesv2.MsgCreateResourcePayload:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	case *resourcetypesv2.Metadata:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	case *resourcetypesv1.MsgCreateResourcePayload:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	case *resourcetypesv1.ResourceHeader:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	default:
		err = json.Unmarshal(file, msg)
	}
	if err != nil {
		return err
	}
	return nil
}
