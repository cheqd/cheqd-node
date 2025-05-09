package types

import (
	fmt "fmt"

	"cosmossdk.io/x/tx/signing"
	didv2 "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	protov2 "google.golang.org/protobuf/proto"
)

func CreateGetSigners(options *signing.Options) func(msg protov2.Message) ([][]byte, error) {
	return func(msg protov2.Message) ([][]byte, error) {
		switch msg := msg.(type) {

		case *didv2.MsgCreateDidDoc:
			return [][]byte{}, nil

		case *didv2.MsgDeactivateDidDoc:
			return [][]byte{}, nil

		case *didv2.MsgUpdateDidDoc:
			return [][]byte{}, nil

		default:
			return nil, fmt.Errorf("unsupported message type: %T", msg)
		}
	}
}
