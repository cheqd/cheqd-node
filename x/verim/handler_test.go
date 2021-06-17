package verim

import (
	"testing"
	"log"

	"github.com/stretchr/testify/require"
	"github.com/verim-id/verim-node/x/verim/types"
)

func TestHandler_CreateNym(t *testing.T) {
	setup := Setup()

	// add new NYM
	nymMsg := TestMsgCreateNym()
	result, _ := setup.Handler(setup.Ctx, nymMsg)
	nym := types.MsgCreateNymResponse{}
	err := nym.Unmarshal(result.Data)
	if err != nil {
		log.Fatal(err)
	}

	// query NYM
	receivedNym := setup.NymKeeper.GetNym(setup.Ctx, nym.Id)

	//// check
	require.Equal(t, receivedNym.Id, nym.Id)
	require.Equal(t, nymMsg.GetCreator(), receivedNym.GetCreator())
	require.Equal(t, nymMsg.GetAlias(), receivedNym.GetAlias())
	require.Equal(t, nymMsg.GetDid(), receivedNym.GetDid())
	require.Equal(t, nymMsg.GetVerkey(), receivedNym.GetVerkey())
	require.Equal(t, nymMsg.GetRole(), receivedNym.GetRole())
}
func TestHandler_UpdateNym(t *testing.T) {
	setup := Setup()

	// add new NYM
	nymMsg := TestMsgCreateNym()
	result, _ := setup.Handler(setup.Ctx, nymMsg)
	nym := types.MsgCreateNymResponse{}
	err := nym.Unmarshal(result.Data)
	if err != nil {
		log.Fatal(err)
	}

	// update NYM
	updateNymMsg := TestMsgUpdateNym(nym.GetId())
	result, _ = setup.Handler(setup.Ctx, updateNymMsg)

	// query NYM
	receivedNym := setup.NymKeeper.GetNym(setup.Ctx, nym.Id)

	//// check
	require.Equal(t, receivedNym.Id, nym.Id)
	require.Equal(t, updateNymMsg.GetCreator(), receivedNym.GetCreator())
	require.Equal(t, updateNymMsg.GetAlias(), receivedNym.GetAlias())
	require.Equal(t, updateNymMsg.GetDid(), receivedNym.GetDid())
	require.Equal(t, updateNymMsg.GetVerkey(), receivedNym.GetVerkey())
	require.Equal(t, updateNymMsg.GetRole(), receivedNym.GetRole())
}

func TestHandler_DeleteNym(t *testing.T) {
	setup := Setup()

	// add new NYM
	nymMsg := TestMsgCreateNym()
	result, _ := setup.Handler(setup.Ctx, nymMsg)
	nym := types.MsgCreateNymResponse{}
	err := nym.Unmarshal(result.Data)
	if err != nil {
		log.Fatal(err)
	}

	// delete NYM
	updateNymMsg := TestMsgDeleteNym(nym.GetId())
	result, _ = setup.Handler(setup.Ctx, updateNymMsg)

	// query NYM
	receivedNym := setup.NymKeeper.GetNym(setup.Ctx, nym.Id)

	//// check
	require.Equal(t, receivedNym.Id, nym.Id)
	require.Equal(t, "", receivedNym.GetCreator())
	require.Equal(t, "", receivedNym.GetAlias())
	require.Equal(t, "", receivedNym.GetDid())
	require.Equal(t, "", receivedNym.GetVerkey())
	require.Equal(t, "", receivedNym.GetRole())
}
