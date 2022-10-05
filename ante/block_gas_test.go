package ante_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	// "github.com/cosmos/cosmos-sdk/baseapp"
	cheqdsimapp "github.com/cheqd/cheqd-node/simapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"

	// banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cheqdtests "github.com/cheqd/cheqd-node/x/cheqd/tests"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	resourcetests "github.com/cheqd/cheqd-node/x/resource/tests"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

const (
	RangeSlippage                   = 12_000
	ComputationalGasMsgCreateDid    = 134110
	ComputationalGasMsgCreateDid80x = 2110398 // 80x the MsgCreateDid gas usage
)

// blockMaxGas = uint64(cheqdsimapp.DefaultConsensusParams.Block.MaxGas) // 200000 keeping this for reference on TODOs
var keyPair = cheqdtests.GenerateKeyPair()

func TestBaseApp_BlockGas(t *testing.T) {
	testcases := []struct {
		name                    string
		gasToConsume            uint64    // gas to consume in the msg execution
		msg                     []sdk.Msg // msgs to execute
		feeAmount               sdk.Coins // fee amount
		panicTx                 bool      // panic explicitly in tx execution
		expErr                  bool
		createDidBeforeResource *cheqdtypes.MsgCreateDid
	}{
		{
			"less than block gas meter - single MsgCreateDid",
			ComputationalGasMsgCreateDid,
			[]sdk.Msg{NewTestDidMsg_CreateDid_Valid(nil)},
			NewTestFeeAmountMinimalDenomEFixedFee_CreateDid(),
			false,
			false,
			nil,
		},
		{
			"more than block gas meter - 80x MsgCreateDid",
			2120438,
			[]sdk.Msg{
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
				NewTestDidMsg_CreateDid_Valid(nil),
			},
			NewTestFeeAmountMinimalDenomEFixedFee_CreateDid().MulInt(sdk.NewInt(80)),
			false,
			true,
			nil,
		},
		{
			"less than block gas meter - single MsgCreateResource",
			70000,
			[]sdk.Msg{NewTestResourceMsg_Json_Valid(keyPair, resourcetests.SchemaData)},
			NewTestFeeAmountMinimalDenomEFixedFee_CreateResourceJson(),
			false,
			true, // TODO: Explicitly set to true, because the KVStore needs state to be set. This should be false.
			NewTestDidMsg_CreateDid_Valid(keyPair),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var app *cheqdsimapp.SimApp
			encCfg := cheqdsimapp.MakeTestEncodingConfig()
			app = cheqdsimapp.NewSimApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, map[int64]bool{}, "", 0, encCfg, simapp.EmptyAppOptions{})
			genState := cheqdsimapp.GenesisStateWithSingleValidator(t, app)
			stateBytes, err := tmjson.MarshalIndent(genState, "", " ")

			require.NoError(t, err)

			app.InitChain(abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			})

			ctx := app.NewContext(false, tmproto.Header{})

			// tx fee
			feeCoin := tc.feeAmount[0]

			// test account and fund
			priv1, _, addr1 := testdata.KeyTestPubAddr()
			expectedSeq := uint64(0)

			// check for did, if resource creation is being tested
			if tc.createDidBeforeResource != nil {
				// set did in state
				didDoc := tc.createDidBeforeResource.Payload.ToDid()
				metadata := cheqdtypes.NewMetadataFromContext(ctx)
				stateValue, err := cheqdtypes.NewStateValue(&didDoc, &metadata)
				require.NoError(t, err)
				err = app.CheqdKeeper.SetDid(&ctx, &stateValue)
				require.NoError(t, err)
			}

			err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, tc.feeAmount)
			require.NoError(t, err)
			err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr1, tc.feeAmount)
			require.NoError(t, err)
			require.Equal(t, feeCoin.Amount, app.BankKeeper.GetBalance(ctx, addr1, feeCoin.Denom).Amount)
			seq, _ := app.AccountKeeper.GetSequence(ctx, addr1)
			require.Equal(t, expectedSeq, seq)

			// msg and signatures
			txBuilder := encCfg.TxConfig.NewTxBuilder()
			require.NoError(t, txBuilder.SetMsgs(tc.msg...))
			txBuilder.SetFeeAmount(tc.feeAmount)
			txBuilder.SetFeePayer(addr1)
			txBuilder.SetGasLimit(txtypes.MaxGasWanted) // tx validation checks that gasLimit can't be bigger than this

			privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{7}, []uint64{seq}
			_, txBytes, err := createTestTx(encCfg.TxConfig, txBuilder, privs, accNums, accSeqs, ctx.ChainID())
			require.NoError(t, err)

			app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1}})
			rsp := app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})

			// check result
			ctx = app.GetContextForDeliverTx(txBytes)

			if tc.expErr {
				if tc.panicTx {
					require.Equal(t, sdkerrors.ErrPanic.ABCICode(), rsp.Code)
				} else {
					require.Equal(t, sdkerrors.ErrOutOfGas.ABCICode(), rsp.Code)
				}
			} else {
				require.Equal(t, uint32(0), rsp.Code)
			}
			// check block gas is always consumed

			if tc.gasToConsume > txtypes.MaxGasWanted {
				// capped by gasLimit
				tc.gasToConsume = txtypes.MaxGasWanted
			}
			// CONTRACT: gasToConsume is +/- 12k gas units from the actual gas consumed (required for larger computations)
			require.GreaterOrEqual(t, tc.gasToConsume+12_000, ctx.BlockGasMeter().GasConsumed())
			require.LessOrEqual(t, tc.gasToConsume-12_000, ctx.BlockGasMeter().GasConsumed())
			// tx fee is always deducted
			require.Equal(t, int64(0), app.BankKeeper.GetBalance(ctx, addr1, feeCoin.Denom).Amount.Int64())
			// sender's sequence is always increased
			seq, err = app.AccountKeeper.GetSequence(ctx, addr1)
			require.NoError(t, err)
			expectedSeq++
			require.Equal(t, expectedSeq, seq)
		})
	}
}

func createTestTx(txConfig client.TxConfig, txBuilder client.TxBuilder, privs []cryptotypes.PrivKey, accNums []uint64, accSeqs []uint64, chainID string) (xauthsigning.Tx, []byte, error) {
	// First round: we gather all the signer infos. We use the "set empty
	// signature" hack to do that.
	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  txConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, nil, err
	}

	// Second round: all signer infos are set, so each signer can sign.
	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			Address:       sdk.AccAddress(priv.PubKey().Bytes()).String(),
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		sigV2, err := tx.SignWithPrivKey(
			txConfig.SignModeHandler().DefaultMode(), signerData,
			txBuilder, priv, txConfig, accSeqs[i])
		if err != nil {
			return nil, nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err = txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, nil, err
	}

	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, nil, err
	}

	return txBuilder.GetTx(), txBytes, nil
}
