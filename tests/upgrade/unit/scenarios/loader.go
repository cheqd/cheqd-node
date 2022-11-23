package scenarios

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	// integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/unit/setup"
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

type ILoader interface {
	LoadFile(path string, dataChunk any, setup migrationsetup.TestSetup) error
	GetLsitOfFiles(path_to_dir, prefix string) ([]string, error)
}

type Loader struct{}

func (l Loader) GetLsitOfFiles(path_to_dir, prefix string) ([]string, error) {

	files_to_load := []string{}
	err := filepath.Walk(path_to_dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasPrefix(info.Name(), prefix) {
			files_to_load = append(files_to_load, path)
		}
		return nil
	})
	return files_to_load, err
}

func (l Loader) LoadFile(
	path string,
	dataChunk any,
	setup migrationsetup.TestSetup,
) error {
	file, err := os.ReadFile(path)
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
