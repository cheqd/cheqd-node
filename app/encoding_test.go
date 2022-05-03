package app

import (
	"testing"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/stretchr/testify/require"
)

func Test_MsgCreateDidPayload_UnmarshalJSON(t *testing.T) {
	cdc := MakeEncodingConfig()

	createPayloadJson := `{"id": "did:cheqd:alice"}`

	var createPayload types.MsgCreateDidPayload
	err := cdc.Codec.UnmarshalJSON([]byte(createPayloadJson), &createPayload)
	require.NoError(t, err)

	require.Equal(t, "did:cheqd:alice", createPayload.Id, "json unmarshal doesn't work")
}
