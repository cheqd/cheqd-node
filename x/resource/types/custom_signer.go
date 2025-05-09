package types

import (
	fmt "fmt"

	"cosmossdk.io/x/tx/signing"
	resourcev2 "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	protov2 "google.golang.org/protobuf/proto"
)

func CreateGetSigners(options *signing.Options) func(msg protov2.Message) ([][]byte, error) {
	return func(msg protov2.Message) ([][]byte, error) {
		switch msg := msg.(type) {

		case *resourcev2.MsgCreateResource:
			return [][]byte{}, nil

		default:
			return nil, fmt.Errorf("unsupported message type: %T", msg)
		}
	}
}
