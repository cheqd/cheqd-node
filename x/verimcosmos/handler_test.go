package verimcosmos

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/verim-id/verim-cosmos/x/verimcosmos/types"
	"testing"
)

func TestHandler_CreateNym(t *testing.T) {
	setup := Setup()

	// add new model
	nymMsg := TestMsgCreateNym()
	result, er := setup.Handler(setup.Ctx, nymMsg)
	//require.Equal(t, sdk.CodeOK, result.)
	//var id = types.MsgCreateNymResponse.Unmarshal(*[]byte(result.Data))
	id := proto.Unmarshal(result.Data, &types.MsgCreateNymResponse)
	//var receivedModelInfo types.ModelInfo
	//_ = setup.Cdc(result, &receivedModelInfo)
	print(id)
	print(er)

	// query model
	//receivedNym := queryGetNym(setup, result.)
	//receivedNym := setup.NymKeeper.GetNym(setup.Ctx, id)
	//
	//// check
	//require.Equal(t, receivedNym., nymMsg)
}

//func queryGetNym(setup TestSetup, id unit64) types.Nym {
//
//
//
//	result, _ := setup.Querier(
//		setup.Ctx,
//		[]string{keeper.QueryModel, fmt.Sprintf("%v", vid), TestQueryGetNym(id)},
//		abci.RequestQuery{},
//	)
//
//	var receivedNym types.Nym
//	_ = setup.Cdc.UnmarshalJSON(result, &receivedNym)
//
//	return receivedNym
//
//}
