package migration

import (
	"encoding/json"
	"os"
	"path/filepath"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
)

type KeyPairBase64 struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type SignInput struct {
	VerificationMethodId string `json:"verificationMethodId"`
	PrivateKey           []byte `json:"privateKey"`
}

func Loader(path string, msg any) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	file, err := os.ReadFile(filepath.Join(cwd, "migration", GENERATED_JSON_DIR, path))
	if err != nil {
		return err
	}
	switch msg := msg.(type) {
	case *didtypes.MsgCreateDidDocPayload:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	case *resourcetypes.MsgCreateResourcePayload:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	case *resourcetypes.Metadata:
		err = integrationhelpers.Codec.UnmarshalJSON(file, msg)
	default:
		err = json.Unmarshal(file, msg)
	}
	if err != nil {
		return err
	}
	return nil
}
