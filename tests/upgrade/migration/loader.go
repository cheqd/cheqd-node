package migration

import (
	"encoding/json"
	"os"
	"path/filepath"

	// integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/migration/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	// didutils "github.com/cheqd/cheqd-node/x/did/utils"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
)

type KeyPairBase64 struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type SignInput struct {
	VerificationMethodId string `json:"verificationMethodId"`
	PrivateKey           []byte `json:"privateKey"`
}

type DidAndMetadata struct {
	Data     didtypesv1.Did
	Metadata didtypesv1.Metadata
}

func Loader(
	path string,
	dataChunk any,
	setup migrationsetup.TestSetup,
) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	file, err := os.ReadFile(filepath.Join(cwd, GENERATED_JSON_DIR, path))
	if err != nil {
		return err
	}
	switch dataChunk := dataChunk.(type) {
	case *didtypesv1.StateValue:
		var temp_s DidAndMetadata
		var stateValue didtypesv1.StateValue
		err = json.Unmarshal(file, &temp_s)
		if err != nil {
			return err
		}

		stateValue, err = didtypesv1.NewStateValue(&temp_s.Data, &temp_s.Metadata)
		if err != nil {
			return err
		}
		*dataChunk = stateValue

	case *resourcetypesv1.Resource:
		err = json.Unmarshal(file, dataChunk)

	case *didtypes.DidDocWithMetadata:
		// err = json.Unmarshal(file, dataChunk)
		err = setup.Cdc.UnmarshalJSON(file, dataChunk)
	case *resourcetypes.ResourceWithMetadata:
		err = json.Unmarshal(file, dataChunk)
		// err = integrationhelpers.Codec.UnmarshalJSON(file, dataChunk)
	default:
		err = json.Unmarshal(file, dataChunk)
	}
	if err != nil {
		return err
	}
	return nil
}
