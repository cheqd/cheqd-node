package migration

import (
	"encoding/json"
	"os"

	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
)

type KeyPairBase64 struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type SignInput struct {
	VerificationMethodId string `json:"verificationMethodId"`
	PrivateKey           string `json:"privateKey"`
}

func Loader[T cheqdtypes.MsgCreateDidPayload | resourcetypes.MsgCreateResourcePayload | resourcetypes.ResourceHeader | KeyPairBase64 | SignInput](path string, msg *T) error {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&msg)
	if err != nil {
		return err
	}
	return nil
}
