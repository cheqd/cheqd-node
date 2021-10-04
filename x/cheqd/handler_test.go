package cheqd

import (
	"log"
	"testing"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	ptypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateDid(t *testing.T) {
	setup := Setup()

	// add new Did
	didMsg := setup.CreateDid()
	data, _ := ptypes.NewAnyWithValue(didMsg)
	result, _ := setup.Handler(setup.Ctx, setup.WrapRequest(data))
	did := types.MsgCreateDidResponse{}
	err := did.Unmarshal(result.Data)

	if err != nil {
		log.Fatal(err)
	}

	// query Did
	receivedDid := setup.Keeper.GetDid(setup.Ctx, did.Id)

	//// check
	require.Equal(t, did.Id, receivedDid.Id)
	require.Equal(t, didMsg.GetController(), receivedDid.GetController())
	require.Equal(t, didMsg.GetVerificationMethod(), receivedDid.GetVerificationMethod())
	require.Equal(t, didMsg.GetAuthentication(), receivedDid.GetAuthentication())
	require.Equal(t, didMsg.GetAssertionMethod(), receivedDid.GetAssertionMethod())
	require.Equal(t, didMsg.GetCapabilityInvocation(), receivedDid.GetCapabilityInvocation())
	require.Equal(t, didMsg.GetCapabilityDelegation(), receivedDid.GetCapabilityDelegation())
	require.Equal(t, didMsg.GetKeyAgreement(), receivedDid.GetKeyAgreement())
	require.Equal(t, didMsg.GetAlsoKnownAs(), receivedDid.GetAlsoKnownAs())
	require.Equal(t, didMsg.GetService(), receivedDid.GetService())
}
