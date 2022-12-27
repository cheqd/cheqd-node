package setup

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
)

type KeyPairBase64 struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type SignInput struct {
	VerificationMethodID string `json:"verificationMethodId"`
	PrivateKey           []byte `json:"privateKey"`
}

type DidAndMetadata struct {
	Data     didtypesv1.Did
	Metadata didtypesv1.Metadata
}

type ILoader interface {
	LoadFile(path string, dataChunk any, setup TestSetup) error
	GetListOfFiles(pathToDir, prefix string) ([]string, error)
}

type Loader struct{}

func (l Loader) GetListOfFiles(pathToDir, prefix string) ([]string, error) {
	filesToLoad := []string{}
	err := filepath.Walk(pathToDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasPrefix(info.Name(), prefix) {
			filesToLoad = append(filesToLoad, path)
		}
		return nil
	})
	return filesToLoad, err
}

func (l Loader) LoadFile(
	path string,
	dataChunk any,
	setup TestSetup,
) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	switch dataChunk := dataChunk.(type) {
	case *didtypesv1.StateValue:
		var tempS DidAndMetadata
		var stateValue didtypesv1.StateValue
		err = json.Unmarshal(file, &tempS)
		if err != nil {
			return err
		}

		stateValue, err = didtypesv1.NewStateValue(&tempS.Data, &tempS.Metadata)
		if err != nil {
			return err
		}
		*dataChunk = stateValue

	case *resourcetypesv1.Resource:
		err = json.Unmarshal(file, dataChunk)

	case *didtypes.DidDocWithMetadata:
		err = setup.Cdc.UnmarshalJSON(file, dataChunk)
	case *resourcetypes.ResourceWithMetadata:
		err = json.Unmarshal(file, dataChunk)
	default:
		err = json.Unmarshal(file, dataChunk)
	}
	if err != nil {
		return err
	}
	return nil
}
