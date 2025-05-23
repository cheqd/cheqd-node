package ante_test

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	"github.com/skip-mev/feemarket/x/feemarket/types"
	"github.com/stretchr/testify/suite"

	cheqdapp "github.com/cheqd/cheqd-node/app"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

var (
	// DefaultWindow is the default window size for the sliding window
	// used to calculate the base fee. In the base EIP-1559 implementation,
	// only the previous block is considered.
	DefaultWindow uint64 = 1

	// DefaultAlpha is not used in the base EIP-1559 implementation.
	DefaultAlpha = math.LegacyMustNewDecFromStr("0.0")

	// DefaultBeta is not used in the base EIP-1559 implementation.
	DefaultBeta = math.LegacyMustNewDecFromStr("1.0")

	// DefaultGamma is not used in the base EIP-1559 implementation.
	DefaultGamma = math.LegacyMustNewDecFromStr("0.0")

	// DefaultDelta is not used in the base EIP-1559 implementation.
	DefaultDelta = math.LegacyMustNewDecFromStr("0.0")

	// DefaultMaxBlockUtilization is the default maximum block utilization. This is the default
	// on Ethereum. This denominated in units of gas consumed in a block.
	DefaultMaxBlockUtilization uint64 = 30_000_000

	// DefaultMinBaseGasPrice is the default minimum base fee.
	DefaultMinBaseGasPrice = math.LegacyOneDec()

	// DefaultMinLearningRate is not used in the base EIP-1559 implementation.
	DefaultMinLearningRate = math.LegacyMustNewDecFromStr("0.125")

	// DefaultMaxLearningRate is not used in the base EIP-1559 implementation.
	DefaultMaxLearningRate = math.LegacyMustNewDecFromStr("0.125")

	// DefaultFeeDenom is the Cosmos SDK default bond denom.
	DefaultFeeDenom = didtypes.BaseMinimalDenom
)

// TestAccount represents an account used in the tests in x/auth/ante.
type TestAccount struct {
	acc  sdk.AccountI
	priv cryptotypes.PrivKey
}

// AnteTestSuite is a test suite to be used with ante handler tests.
type AnteTestSuite struct {
	suite.Suite

	app         *cheqdapp.TestApp
	anteHandler sdk.AnteHandler
	ctx         sdk.Context
	clientCtx   client.Context
	txBuilder   client.TxBuilder
}

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*cheqdapp.TestApp, sdk.Context, error) {
	app, err := cheqdapp.Setup(isCheckTx)
	if err != nil {
		return nil, sdk.Context{}, err
	}
	ctx := app.BaseApp.NewContext(isCheckTx)
	err = app.AccountKeeper.Params.Set(ctx, authtypes.DefaultParams())
	if err != nil {
		return nil, sdk.Context{}, err
	}

	// cheqd specific params
	didFeeParams := didtypes.DefaultGenesis().FeeParams
	err = app.DidKeeper.SetParams(ctx, *didFeeParams)
	if err != nil {
		return nil, sdk.Context{}, err
	}
	resourceFeeParams := resourcetypes.DefaultGenesis().FeeParams
	_ = app.ResourceKeeper.SetParams(ctx, *resourceFeeParams)
	err = app.FeeMarketKeeper.SetParams(ctx, types.NewParams(DefaultWindow, DefaultAlpha, DefaultBeta, DefaultGamma, DefaultDelta,
		DefaultMaxBlockUtilization, DefaultMinBaseGasPrice, DefaultMinLearningRate, DefaultMaxLearningRate, DefaultFeeDenom, true,
	))
	if err != nil {
		return nil, ctx, err
	}
	return app, ctx, nil
}

// SetupTest setups a new test, with new app, context, and anteHandler.
func (s *AnteTestSuite) SetupTest(isCheckTx bool) error {
	var err error
	s.app, s.ctx, err = createTestApp(isCheckTx)
	if err != nil {
		return err
	}
	s.ctx = s.ctx.WithBlockHeight(1)
	// We're using TestMsg encoding in some tests, so register it here.
	s.app.LegacyAmino().Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	testdata.RegisterInterfaces(s.app.InterfaceRegistry())

	s.clientCtx = client.Context{}.
		WithTxConfig(s.app.TxConfig())

	anteHandler, err := cheqdapp.NewAnteHandler(
		cheqdapp.HandlerOptions{
			AccountKeeper:   s.app.AccountKeeper,
			BankKeeper:      s.app.BankKeeper,
			FeegrantKeeper:  s.app.FeeGrantKeeper,
			DidKeeper:       s.app.DidKeeper,
			ResourceKeeper:  s.app.ResourceKeeper,
			SignModeHandler: s.app.TxConfig().SignModeHandler(),
			SigGasConsumer:  sdkante.DefaultSigVerificationGasConsumer,
			IBCKeeper:       s.app.IBCKeeper,
			FeeMarketKeeper: s.app.FeeMarketKeeper,
		},
	)
	if err != nil {
		return err
	}
	s.anteHandler = anteHandler

	return nil
}

// CreateTestAccounts creates `numAccs` accounts, and return all relevant
// information about them including their private keys.
func (s *AnteTestSuite) CreateTestAccounts(numAccs int) ([]TestAccount, error) {
	var accounts []TestAccount

	for i := 0; i < numAccs; i++ {
		priv, _, addr := testdata.KeyTestPubAddr()
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr)
		err := acc.SetAccountNumber(uint64(i + 1000))
		if err != nil {
			return nil, err
		}
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		someCoins := sdk.Coins{
			sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 1000000*1e9), // 1mn CHEQ
		}
		err = s.app.BankKeeper.MintCoins(s.ctx, minttypes.ModuleName, someCoins)
		if err != nil {
			return nil, err
		}

		err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, addr, someCoins)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, TestAccount{acc, priv})
	}

	return accounts, nil
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (s *AnteTestSuite) CreateTestTx(privs []cryptotypes.PrivKey, accNums []uint64, accSeqs []uint64, chainID string) (xauthsigning.Tx, error) {
	// First round: we gather all the signer infos. We use the "set empty
	// signature" hack to do that.
	sigsV2 := make([]signing.SignatureV2, 0, len(privs))
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  signing.SignMode(s.clientCtx.TxConfig.SignModeHandler().DefaultMode()),
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := s.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	// Second round: all signer infos are set, so each signer can sign.
	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		sigV2, err := tx.SignWithPrivKey(s.ctx,
			signing.SignMode(s.clientCtx.TxConfig.SignModeHandler().DefaultMode()), signerData,
			s.txBuilder, priv, s.clientCtx.TxConfig, accSeqs[i])
		if err != nil {
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err = s.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}
	return s.txBuilder.GetTx(), nil
}

// SetDidFeeParams is a helper function to set did fee params.
func (s *AnteTestSuite) SetDidFeeParams(feeParams didtypes.FeeParams) error {
	return s.app.DidKeeper.SetParams(s.ctx, feeParams)
}

// SetResourceFeeParams is a helper function to set resource fee params.
func (s *AnteTestSuite) SetResourceFeeParams(feeParams resourcetypes.FeeParams) error {
	return s.app.ResourceKeeper.SetParams(s.ctx, feeParams)
}

func (s *AnteTestSuite) SetFeeMarketFeeDenom() error {
	err := s.app.FeeMarketKeeper.SetParams(s.ctx, types.Params{FeeDenom: didtypes.BaseMinimalDenom})
	if err != nil {
		return err
	}
	return nil
}

// TestCase represents a test case used in test tables.
type TestCase struct {
	desc     string
	simulate bool
	expPass  bool
	expErr   error
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (s *AnteTestSuite) RunTestCase(privs []cryptotypes.PrivKey, msgs []sdk.Msg, feeAmount sdk.Coins, gasLimit uint64, accNums, accSeqs []uint64, chainID string, tc TestCase) {
	s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
		s.Require().NoError(s.txBuilder.SetMsgs(msgs...))
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)

		// Theoretically speaking, ante handler unit tests should only test
		// ante handlers, but here we sometimes also test the tx creation
		// process.
		tx, txErr := s.CreateTestTx(privs, accNums, accSeqs, chainID)
		newCtx, anteErr := s.anteHandler(s.ctx, tx, tc.simulate)

		if tc.expPass {
			s.Require().NoError(txErr)
			s.Require().NoError(anteErr)
			s.Require().NotNil(newCtx)

			s.ctx = newCtx
		} else {
			switch {
			case txErr != nil:
				s.Require().Error(txErr)
				s.Require().True(errors.Is(txErr, tc.expErr))

			case anteErr != nil:
				s.Require().Error(anteErr)
				s.Require().True(errors.Is(anteErr, tc.expErr))

			default:
				s.Fail("expected one of txErr,anteErr to be an error")
			}
		}
	})
}
