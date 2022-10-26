package cli

import (
	"encoding/base64"

	cheqdcli "github.com/cheqd/cheqd-node/x/cheqd/client/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateResource(collectionId string, resourceId string, resourceName string, resourceType string, resourceFile string, signInputs []cheqdcli.SignInput, container string) (sdk.TxResponse, error) {
	args := []string{
		"--collection-id", collectionId,
		"--resource-id", resourceId,
		"--resource-name", resourceName,
		"--resource-type", resourceType,
		"--resource-file", resourceFile,
	}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodId)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivKey))
	}

	return Tx(container, CLI_BINARY_NAME, "resource", "create-resource", OperatorAccounts[container], args...)
}
