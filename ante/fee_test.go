package ante_test

import (
	"strings"

	"cosmossdk.io/math"
	cheqdante "github.com/cheqd/cheqd-node/ante"
	cheqdpost "github.com/cheqd/cheqd-node/post"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	oraclekeeper "github.com/cheqd/cheqd-node/x/oracle/keeper"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	feeabsante "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/ante"
	"github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/types"
	feemarketante "github.com/skip-mev/feemarket/x/feemarket/ante"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

var _ = Describe("Fee tests on CheckTx", func() {
	s := new(AnteTestSuite)

	var decorators []sdk.AnteDecorator
	BeforeEach(func() {
		err := s.SetupTest(true) // setup
		Expect(err).To(BeNil(), "Error on creating test app")
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()
		decorators = []sdk.AnteDecorator{
			feemarketante.NewFeeMarketCheckDecorator( // fee market check replaces fee deduct decorator
				s.app.AccountKeeper,
				s.app.BankKeeper,
				s.app.FeeGrantKeeper,
				s.app.FeeMarketKeeper,
				ante.NewDeductFeeDecorator(
					s.app.AccountKeeper,
					s.app.BankKeeper,
					s.app.FeeGrantKeeper,
					nil,
				),
			),
		}
	})

	It("Ensure Zero Mempool Fees On Simulation", func() {
		mfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(mfd)

		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(300000000000)))
		err := testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)

		Expect(err).To(BeNil())

		// msg and signatures
		msg := testdata.NewTestMsg(addr1)
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())

		// set zero gas
		s.txBuilder.SetGasLimit(0)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// Set IsCheckTx to true
		s.ctx = s.ctx.WithIsCheckTx(true)

		_, err = antehandler(s.ctx, tx, false)
		Expect(err).NotTo(BeNil())

		// zero gas is accepted in simulation mode
		_, err = antehandler(s.ctx, tx, true)
		Expect(err).To(BeNil())
	})

	It("Ensure Mempool Fees", func() {
		mfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(mfd)

		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(300_000_000_000)))
		err := testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil())

		// msg and signatures
		msg := testdata.NewTestMsg(addr1)
		feeAmount := NewTestFeeAmount()
		gasLimit := uint64(15)
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// Set high gas price so standard test fee fails
		ncheqPrice := sdk.NewDecCoinFromDec(didtypes.BaseMinimalDenom, math.LegacyNewDec(200_000_000_000))

		params, err := s.app.FeeMarketKeeper.GetParams(s.ctx)
		Expect(err).To(BeNil())

		// Since we use the BaseGasPrice set by the feemarket
		// Set high gas price in feemarket
		params.MinBaseGasPrice = ncheqPrice.Amount
		err = s.app.FeeMarketKeeper.SetParams(s.ctx, params)
		Expect(err).To(BeNil())

		state := feemarkettypes.DefaultState()
		state.BaseGasPrice = ncheqPrice.Amount
		err = s.app.FeeMarketKeeper.SetState(s.ctx, feemarkettypes.NewState(
			state.Index,
			state.BaseGasPrice,
			state.LearningRate,
		))
		Expect(err).To(BeNil())

		params, err = s.app.FeeMarketKeeper.GetParams(s.ctx)
		Expect(err).To(BeNil())

		// Set IsCheckTx to true
		s.ctx = s.ctx.WithIsCheckTx(true)

		// antehandler errors with insufficient fees
		_, err = antehandler(s.ctx, tx, false)
		Expect(err).NotTo(BeNil(), "Decorator should have errored on too low fee for local gasPrice")

		// antehandler should not error since we do not check minGasPrice in simulation mode
		params, err = s.app.FeeMarketKeeper.GetParams(s.ctx)
		Expect(err).To(BeNil())

		// Set high gas price in feemarket
		params.Enabled = false
		err = s.app.FeeMarketKeeper.SetParams(s.ctx, params)
		Expect(err).To(BeNil())

		cacheCtx, _ := s.ctx.CacheContext()
		_, err = antehandler(cacheCtx, tx, true)
		Expect(err).To(BeNil(), "Decorator should not have errored in simulation mode")

		// Set IsCheckTx to false
		s.ctx = s.ctx.WithIsCheckTx(false)

		// antehandler should not error since we do not check minGasPrice in DeliverTx
		_, err = antehandler(s.ctx, tx, false)
		Expect(err).To(BeNil(), "MempoolFeeDecorator returned error in DeliverTx")

		// Set IsCheckTx back to true for testing sufficient mempool fee
		s.ctx = s.ctx.WithIsCheckTx(true)

		ncheqPrice = sdk.NewDecCoinFromDec(didtypes.BaseMinimalDenom, math.LegacyNewDec(0).Quo(math.LegacyNewDec(didtypes.BaseFactor)))
		lowGasPrice := []sdk.DecCoin{ncheqPrice}
		s.ctx = s.ctx.WithMinGasPrices(lowGasPrice) // 1 ncheq

		newCtx, err := antehandler(s.ctx, tx, false)
		Expect(err).To(BeNil(), "Decorator should not have errored on fee higher than local gasPrice")
		// Priority is the smallest gas price amount in any denom. Since we have only 1 gas price
		// of 10000000000ncheq, the priority here is 10*10^9.
		Expect(int64(10) * didtypes.BaseFactor).To(Equal(newCtx.Priority()))
	})

	It("TaxableTx Mempool Inclusion", func() {
		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// msg and signatures
		msg := SandboxDidDoc()
		feeAmount := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(5_000_000_000)))
		gasLimit := testdata.NewTestGasLimit()
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// set account with sufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(100_000_000_000))))
		Expect(err).To(BeNil())

		dfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(dfd)

		_, err = antehandler(s.ctx, tx, false)

		Expect(err).To(BeNil(), "Tx errored when taxable on checkTx")

		// set checkTx to false
		s.ctx = s.ctx.WithIsCheckTx(false)

		// antehandler should not error to replay in mempool or deliverTx
		_, err = antehandler(s.ctx, tx, false)

		Expect(err).To(BeNil(), "Tx errored when taxable on deliverTx")
	})
})

var _ = Describe("Fee tests on DeliverTx", func() {
	s := new(AnteTestSuite)
	var decorators []sdk.AnteDecorator

	BeforeEach(func() {
		err := s.SetupTest(false) // setup
		Expect(err).To(BeNil(), "Error on creating test app")
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()
		decorators = []sdk.AnteDecorator{
			feemarketante.NewFeeMarketCheckDecorator( // fee market check replaces fee deduct decorator
				s.app.AccountKeeper,
				s.app.BankKeeper,
				s.app.FeeGrantKeeper,
				s.app.FeeMarketKeeper,
				ante.NewDeductFeeDecorator(
					s.app.AccountKeeper,
					s.app.BankKeeper,
					s.app.FeeGrantKeeper,
					nil,
				),
			),
		}
	})

	It("Deduct Fees", func() {
		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// msg and signatures
		msg := testdata.NewTestMsg(addr1)
		feeAmount := NewTestFeeAmount()
		gasLimit := testdata.NewTestGasLimit()
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// Set account with insufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(10_000_000_000)))
		//nocheck:errcheck
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil())

		dfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(dfd)

		_, err = antehandler(s.ctx, tx, false)

		Expect(err).NotTo(BeNil(), "Tx did not error when fee payer had insufficient funds")

		// Set account with sufficient funds
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(200_000_000_000))))
		Expect(err).To(BeNil())

		_, err = antehandler(s.ctx, tx, false)

		Expect(err).To(BeNil(), "Tx errored after account has been set with sufficient funds")
	})
	It("TaxableTx Lifecycle - DID: MsgCreateDidDoc", func() {
		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// msg and signatures
		msg := SandboxDidDoc()
		feeAmount := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(50_000_000_000)))

		gasLimit := testdata.NewTestGasLimit()
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)
		s.txBuilder.SetFeePayer(addr1)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// set account with sufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		amount := math.NewInt(100_000_000_000)
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
		Expect(err).To(BeNil())

		dfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(dfd)
		taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper, s.app.FeeMarketKeeper, s.app.OracleKeeper, s.app.FeeabsKeeper)
		posthandler := sdk.ChainPostDecorators(taxDecorator)

		// get supply before tx
		supplyBeforeDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// antehandler should not error to replay in mempool or deliverTx
		_, err = antehandler(s.ctx, tx, false)
		Expect(err).To(BeNil(), "Tx errored when taxable on deliverTx")

		_, err = posthandler(s.ctx, tx, false, true)
		Expect(err).To(BeNil(), "Tx errored when fee payer had sufficient funds and provided sufficient fee while subtracting tax on deliverTx")

		// get fee params
		feeParams, err := s.app.DidKeeper.GetParams(s.ctx)
		Expect(err).To(BeNil())

		// check balance of fee payer
		balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, didtypes.BaseMinimalDenom)
		Expect(amount.Sub(math.NewInt(feeParams.CreateDid[0].MinAmount.Int64()))).To(Equal(balance.Amount), "Tax was not subtracted from the fee payer")

		// get supply after tx
		supplyAfterDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// check that supply was deflated
		burnt := cheqdante.GetBurnFeePortion(feeParams.BurnFactor, sdk.NewCoins(sdk.NewCoin(feeParams.CreateDid[0].Denom, *feeParams.CreateDid[0].MinAmount)))
		Expect(supplyBeforeDeflation.Sub(supplyAfterDeflation...)).To(Equal(burnt), "Supply was not deflated")

		// check that reward has been sent to the fee collector
		reward := cheqdante.GetRewardPortion(sdk.NewCoins(sdk.NewCoin(feeParams.CreateDid[0].Denom, *feeParams.CreateDid[0].MinAmount)), burnt)

		// calculate oracle share (0.5%)
		oracleShareRate := math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 3) // 0.005 = 0.5%
		rewardAmt := reward.AmountOf(didtypes.BaseMinimalDenom)
		rewardDec := math.LegacyNewDecFromInt(rewardAmt)
		oracleShare := rewardDec.Mul(oracleShareRate).TruncateInt()

		// expected fee collector reward = reward - oracle share
		expectedFeeCollectorReward := rewardAmt.Sub(oracleShare)

		// check fee collector balance
		feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)
		Expect(feeCollectorBalance.Amount).To(Equal(expectedFeeCollectorReward), "Fee Collector did not receive correct reward after oracle share deduction")

		// check oracle module balance
		oracleModule := s.app.AccountKeeper.GetModuleAddress("oracle")
		oracleBalance := s.app.BankKeeper.GetBalance(s.ctx, oracleModule, didtypes.BaseMinimalDenom)
		Expect(oracleBalance.Amount).To(Equal(oracleShare), "Oracle module did not receive the correct reward share")
	})

	It("TaxableTx Lifecycle - DLR: MsgCreateResource JSON", func() {
		s.app.OracleKeeper.SetAverage(s.ctx, oraclekeeper.KeyEMA(didtypes.BaseDenom), math.LegacyMustNewDecFromStr("0.016"))

		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// msg and signatures
		msg := SandboxResource()

		feeAmount := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(3500000000000)))
		gasLimit := uint64(2_000_000)
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)
		s.txBuilder.SetFeePayer(addr1)
		s.ctx = s.ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoinFromDec(didtypes.BaseMinimalDenom, math.LegacyMustNewDecFromStr("5000"))))

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// set account with sufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		amount := math.NewInt(100_000_000_000)
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
		Expect(err).To(BeNil())

		dfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(dfd)

		taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper, s.app.FeeMarketKeeper, s.app.OracleKeeper, s.app.FeeabsKeeper)
		posthandler := sdk.ChainPostDecorators(taxDecorator)

		// get supply before tx
		supplyBeforeDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// antehandler should not error
		_, err = antehandler(s.ctx, tx, false)
		Expect(err).To(BeNil(), "Tx errored when taxable on deliverTx")

		_, err = posthandler(s.ctx, tx, false, true)
		Expect(err).To(BeNil(), "Tx errored when fee payer had sufficient funds and provided sufficient fee while subtracting tax on deliverTx")

		// get fee params
		feeParams, err := s.app.ResourceKeeper.GetParams(s.ctx)
		Expect(err).To(BeNil())

		// check balance of fee payer
		balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, didtypes.BaseMinimalDenom)
		Expect(amount.Sub(math.NewInt(feeParams.Json[0].MinAmount.Int64()))).To(Equal(balance.Amount), "Tax was not subtracted from the fee payer")

		// get supply after tx
		supplyAfterDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// check that supply was deflated
		burnt := cheqdante.GetBurnFeePortion(feeParams.BurnFactor, sdk.NewCoins(sdk.NewCoin(feeParams.Json[0].Denom, *feeParams.Json[0].MinAmount)))
		Expect(supplyBeforeDeflation.Sub(supplyAfterDeflation...)).To(Equal(burnt), "Supply was not deflated")

		// reward and oracle share logic
		reward := cheqdante.GetRewardPortion(sdk.NewCoins(sdk.NewCoin(feeParams.Json[0].Denom, *feeParams.Json[0].MinAmount)), burnt)
		rewardAmt := reward.AmountOf(didtypes.BaseMinimalDenom) // type: math.Int
		rewardCoin := sdk.NewCoin(didtypes.BaseMinimalDenom, rewardAmt)
		rewardDecCoin := sdk.NewDecCoinFromCoin(rewardCoin)

		oracleShareRate := math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 3) // 0.005 = 0.5%
		oracleShare := rewardDecCoin.Amount.Mul(oracleShareRate).TruncateInt()

		expectedFeeCollectorReward := rewardAmt.Sub(oracleShare)

		feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)

		Expect(feeCollectorBalance.Amount).To(Equal(expectedFeeCollectorReward), "Reward sent to fee collector is incorrect after subtracting oracle share")

		// Optional: verify oracle module received its share
		oracleModule := s.app.AccountKeeper.GetModuleAddress("oracle")
		oracleBalance := s.app.BankKeeper.GetBalance(s.ctx, oracleModule, didtypes.BaseMinimalDenom)

		Expect(oracleBalance.Amount).To(Equal(oracleShare), "Oracle did not receive correct reward share")
	})

	It("Non TaxableTx Lifecycle", func() {
		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// msg and signatures
		msg := testdata.NewTestMsg(addr1)
		feeAmount := NewTestFeeAmount()
		gasLimit := testdata.NewTestGasLimit()
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)
		s.txBuilder.SetFeePayer(addr1)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// set account with sufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		amount := math.NewInt(300_000_000_000)
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
		Expect(err).To(BeNil())

		dfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(dfd)

		taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper, s.app.FeeMarketKeeper, s.app.OracleKeeper, s.app.FeeabsKeeper)
		posthandler := sdk.ChainPostDecorators(taxDecorator)

		// get supply before tx
		supplyBefore, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// antehandler should not error since we make sure that the fee is sufficient in DeliverTx (simulate=false only, Posthandler will check it otherwise)
		newCtx, err := antehandler(s.ctx, tx, false)
		Expect(err).To(BeNil(), "Tx errored when non-taxable on deliverTx")
		_, _, proposer := testdata.KeyTestPubAddr()
		s.ctx = newCtx
		a := s.ctx.BlockHeader()
		a.ProposerAddress = proposer
		newCtx = s.ctx.WithBlockHeader(a)
		s.ctx = newCtx
		_, err = posthandler(s.ctx, tx, false, true)
		Expect(err).To(BeNil(), "Tx errored when non-taxable on deliverTx from posthandler")

		// check balance of fee payer
		balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, didtypes.BaseMinimalDenom)
		Expect(amount.Sub(math.NewInt(feeAmount.AmountOf(didtypes.BaseMinimalDenom).Int64()))).To(Equal(balance.Amount), "Fee amount subtracted was not equal to fee amount required for non-taxable tx")

		// get supply after tx
		supplyAfter, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// check that supply was not deflated
		Expect(supplyBefore).To(Equal(supplyAfter), "Supply was deflated")

		// check that reward has been sent to the fee collector
		feeCollector := s.app.AccountKeeper.GetModuleAddress(feemarkettypes.FeeCollectorName)
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)

		Expect((feeCollectorBalance.Amount).GT(math.NewInt(0)))
	})

	It("Non TaxableTx Lifecycle - Ensure minimum gas prices", func() {
		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// msg and signatures
		msg := testdata.NewTestMsg(addr1)
		feeAmount := sdk.NewCoins(sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 1*didtypes.BaseFactor)) // 1 CHEQ
		gasLimit := testdata.NewTestGasLimit()
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)
		s.txBuilder.SetFeePayer(addr1)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// set account with sufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		amount := math.NewInt(300_000_000_000) // 300 CHEQ
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
		Expect(err).To(BeNil())

		dfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(dfd)

		// get feemarket params
		feemarketParams, err := s.app.FeeMarketKeeper.GetParams(s.ctx)
		Expect(err).To(BeNil())

		// enforced enablement
		feemarketParams.Enabled = true

		// enforce burn
		feemarketParams.DistributeFees = false

		// enforce maximum block utilisation
		feemarketParams.MaxBlockUtilization = feemarkettypes.DefaultMaxBlockUtilization

		// set minimum gas prices to realistic value
		feemarketParams.MinBaseGasPrice = math.LegacyMustNewDecFromStr("0.5")

		err = s.app.FeeMarketKeeper.SetParams(s.ctx, feemarketParams)
		Expect(err).To(BeNil())

		// get feemarket state
		feemarketState, err := s.app.FeeMarketKeeper.GetState(s.ctx)
		Expect(err).To(BeNil())

		// set base gas price to realistic value
		feemarketState.BaseGasPrice = math.LegacyMustNewDecFromStr("0.5")

		// set feemarket state
		err = s.app.FeeMarketKeeper.SetState(s.ctx, feemarketState)
		Expect(err).To(BeNil())

		taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper, s.app.FeeMarketKeeper, s.app.OracleKeeper, s.app.FeeabsKeeper)
		posthandler := sdk.ChainPostDecorators(taxDecorator)

		// get supply before tx
		supplyBefore, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// antehandler should not error since we make sure that the fee is sufficient in DeliverTx (simulate=false only, Posthandler will check it otherwise)
		newCtx, err := antehandler(s.ctx, tx, false)
		Expect(err).To(BeNil(), "Tx errored when non-taxable on deliverTx")
		_, _, proposer := testdata.KeyTestPubAddr()
		s.ctx = newCtx
		a := s.ctx.BlockHeader()
		a.ProposerAddress = proposer
		newCtx = s.ctx.WithBlockHeader(a)
		s.ctx = newCtx
		_, err = posthandler(s.ctx, tx, false, true)
		Expect(err).To(BeNil(), "Tx errored when non-taxable on deliverTx from posthandler")

		// check balance of fee payer
		balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, didtypes.BaseMinimalDenom)
		Expect(amount.Sub(math.NewInt(feeAmount.AmountOf(didtypes.BaseMinimalDenom).Int64()))).To(Equal(balance.Amount), "Fee amount subtracted was not equal to fee amount required for non-taxable tx")

		// get supply after tx
		supplyAfter, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// check that supply was not deflated
		Expect(supplyBefore).To(Equal(supplyAfter), "Supply was deflated")

		// check that reward has been sent to the fee collector
		feeCollector := s.app.AccountKeeper.GetModuleAddress(feemarkettypes.FeeCollectorName)
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)

		Expect((feeCollectorBalance.Amount).GT(math.NewInt(0)))
	})

	It("TaxableTx Lifecycle on Simulation", func() {
		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		msg := SandboxDidDoc()
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())

		// set zero gas
		s.txBuilder.SetGasLimit(0)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// set account with sufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		amount := math.NewInt(100_000_000_000)
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
		Expect(err).To(BeNil())

		dfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(dfd)

		taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper, s.app.FeeMarketKeeper, s.app.OracleKeeper, s.app.FeeabsKeeper)
		posthandler := sdk.ChainPostDecorators(taxDecorator)

		// get supply before tx
		supplyBefore, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// antehandler should not error since we make sure that the fee is sufficient in DeliverTx (simulate=false only, Posthandler will check it otherwise)
		_, err = antehandler(s.ctx, tx, true)
		Expect(err).To(BeNil(), "Tx errored when taxable on deliverTx simulation mode")

		_, err = posthandler(s.ctx, tx, true, true)
		Expect(err).To(BeNil(), "Tx errored when fee payer had sufficient funds and provided sufficient fee while skipping tax simulation mode")

		// check balance of fee payer
		balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, didtypes.BaseMinimalDenom)
		Expect(amount).To(Equal(balance.Amount), "Tax was subtracted from fee payer when taxable tx was simulated")

		// get supply after tx
		supplyAfter, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// check that supply was not deflated
		Expect(supplyBefore).To(Equal(supplyAfter), "Supply was deflated when taxable tx simulation mode")

		// check that no fee has been sent to the fee collector
		feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)

		Expect(feeCollectorBalance.Amount).To(Equal(math.NewInt(0)), "Reward was sent to the fee collector when taxable tx simulation mode")
	})
})

var _ = Describe("Fee abstraction", func() {
	var gasLimit uint64
	var mockHostZoneConfig types.HostChainFeeAbsConfig

	// mockHostZoneConfig is used to mock the host zone config, with ibcfee as the ibc fee denom to be used as alternative fee
	BeforeEach(func() {
		gasLimit = 200000
		mockHostZoneConfig = types.HostChainFeeAbsConfig{
			IbcDenom:                "ibcfee",
			OsmosisPoolTokenDenomIn: "osmosis",
			PoolId:                  1,
			Status:                  types.HostChainFeeAbsStatus_UPDATED,
		}
	})

	// Define test cases inside a context
	Context("Testing MempoolDecorator with different fee amounts", func() {
		var suite *AnteTestSuite

		BeforeEach(func() {
			// Set up the test suite for each test case
			suite = new(AnteTestSuite)
			err := suite.SetupTest(true) // setup
			Expect(err).To(BeNil(), "Error on creating test app")
			suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
		})

		It("should fail with empty fee", func() {
			feeAmount := sdk.Coins{}
			minGasPrice := sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 100))...)
			suite.txBuilder.SetGasLimit(gasLimit)
			suite.txBuilder.SetFeeAmount(feeAmount)
			suite.ctx = suite.ctx.WithMinGasPrices(minGasPrice)

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			mempoolDecorator := feeabsante.NewFeeAbstrationMempoolFeeDecorator(suite.app.FeeabsKeeper)
			anteHandler := sdk.ChainAnteDecorators(mempoolDecorator)

			_, err := anteHandler(suite.ctx, tx, false)

			// Expect error due to insufficient fee
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), errors.ErrInsufficientFee.Error())).To(BeTrue())
		})

		It("should fail with insufficient native fee", func() {
			feeAmount := sdk.NewCoins(sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 100))
			minGasPrice := sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 1000))...)
			suite.txBuilder.SetGasLimit(gasLimit)
			suite.txBuilder.SetFeeAmount(feeAmount)
			suite.ctx = suite.ctx.WithMinGasPrices(minGasPrice)

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			mempoolDecorator := feeabsante.NewFeeAbstrationMempoolFeeDecorator(suite.app.FeeabsKeeper)
			anteHandler := sdk.ChainAnteDecorators(mempoolDecorator)

			_, err := anteHandler(suite.ctx, tx, false)

			// Expect error due to insufficient fee
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), errors.ErrInsufficientFee.Error())).To(BeTrue())
		})

		It("should pass with sufficient native fee", func() {
			feeAmount := sdk.NewCoins(sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 1000*int64(gasLimit)))
			minGasPrice := sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 1000))...)
			suite.txBuilder.SetGasLimit(gasLimit)
			suite.txBuilder.SetFeeAmount(feeAmount)
			suite.ctx = suite.ctx.WithMinGasPrices(minGasPrice)

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			mempoolDecorator := feeabsante.NewFeeAbstrationMempoolFeeDecorator(suite.app.FeeabsKeeper)
			anteHandler := sdk.ChainAnteDecorators(mempoolDecorator)

			_, err := anteHandler(suite.ctx, tx, false)

			// No error is expected
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with unknown ibc fee denom", func() {
			feeAmount := sdk.NewCoins(sdk.NewInt64Coin("ibcfee", 1000*int64(gasLimit)))
			minGasPrice := sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 1000))...)
			suite.txBuilder.SetGasLimit(gasLimit)
			suite.txBuilder.SetFeeAmount(feeAmount)
			suite.ctx = suite.ctx.WithMinGasPrices(minGasPrice)

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			mempoolDecorator := feeabsante.NewFeeAbstrationMempoolFeeDecorator(suite.app.FeeabsKeeper)
			anteHandler := sdk.ChainAnteDecorators(mempoolDecorator)

			_, err := anteHandler(suite.ctx, tx, false)

			// Expect error due to unknown ibc fee denom
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), errors.ErrInvalidCoins.Error())).To(BeTrue())
		})

		It("should pass with sufficient ibc fee", func() {
			feeAmount := sdk.NewCoins(sdk.NewInt64Coin("ibcfee", 1000*int64(gasLimit)))
			minGasPrice := sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.NewInt64Coin("stake", 1000))...)
			suite.txBuilder.SetGasLimit(gasLimit)
			suite.txBuilder.SetFeeAmount(feeAmount)
			suite.ctx = suite.ctx.WithMinGasPrices(minGasPrice)

			// Configure the HostZoneConfig
			err := suite.app.FeeabsKeeper.SetHostZoneConfig(suite.ctx, mockHostZoneConfig)
			Expect(err).ToNot(HaveOccurred())
			suite.app.FeeabsKeeper.SetTwapRate(suite.ctx, "ibcfee", math.LegacyNewDec(1))

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			mempoolDecorator := feeabsante.NewFeeAbstrationMempoolFeeDecorator(suite.app.FeeabsKeeper)
			anteHandler := sdk.ChainAnteDecorators(mempoolDecorator)

			_, err = anteHandler(suite.ctx, tx, false)

			// No error is expected
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

var _ = Describe("DeductFeeDecorator", func() {
	var gasLimit uint64
	var minGasPrice sdk.DecCoins
	var feeAmount sdk.Coins
	var ibcFeeAmount sdk.Coins
	var mockHostZoneConfig types.HostChainFeeAbsConfig
	var suite *AnteTestSuite
	var testAcc TestAccount

	// Setup the common test data
	BeforeEach(func() {
		gasLimit = 200000
		minGasPrice = sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 1000))...)
		feeAmount = sdk.NewCoins(sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 1000*int64(gasLimit)))
		ibcFeeAmount = sdk.NewCoins(sdk.NewInt64Coin("ibcfee", 1000*int64(gasLimit)))

		mockHostZoneConfig = types.HostChainFeeAbsConfig{
			IbcDenom:                "ibcfee",
			OsmosisPoolTokenDenomIn: "osmosis",
			PoolId:                  1,
			Status:                  types.HostChainFeeAbsStatus_UPDATED,
		}

		// Initialize suite for each test case
		suite = new(AnteTestSuite)
		err := suite.SetupTest(true) // setup
		Expect(err).To(BeNil(), "Error on creating test app")
		suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
		acc, err := suite.CreateTestAccounts(1)
		testAcc = acc[0]
		Expect(err).To(BeNil())
		Expect(len(acc)).To(Equal(1))
		suite.txBuilder.SetGasLimit(gasLimit)
		suite.txBuilder.SetFeeAmount(feeAmount)
		suite.txBuilder.SetFeePayer(acc[0].acc.GetAddress())
		suite.ctx = suite.ctx.WithMinGasPrices(minGasPrice)
		// minFee, _ := minGasPrice.TruncateDecimal()

		params, err := suite.app.StakingKeeper.GetParams(suite.ctx)
		Expect(err).To(BeNil(), "Error getting staking params")
		params.BondDenom = didtypes.BaseMinimalDenom
		err = suite.app.StakingKeeper.SetParams(suite.ctx, params)
		Expect(err).To(BeNil(), "Error setting the params")

		// this line will create the module account
		_ = suite.app.AccountKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	})

	When("native fee is sufficient", func() {
		It("should pass with sufficient native fee", func() {
			suite.app.FeeabsKeeper.SetTwapRate(suite.ctx, "ibcfee", math.LegacyNewDec(1))

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			deductFeeDecorator := feeabsante.NewFeeAbstractionDeductFeeDecorate(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.FeeabsKeeper, suite.app.FeeGrantKeeper)
			anteHandler := sdk.ChainAnteDecorators(deductFeeDecorator)

			_, err := anteHandler(suite.ctx, tx, false)

			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("ibc fee is insufficient", func() {
		It("should fail due to insufficient ibc fee", func() {
			err := suite.app.FeeabsKeeper.SetHostZoneConfig(suite.ctx, mockHostZoneConfig)
			Expect(err).ToNot(HaveOccurred())
			suite.app.FeeabsKeeper.SetTwapRate(suite.ctx, "ibcfee", math.LegacyNewDec(1))

			suite.txBuilder.SetFeeAmount(ibcFeeAmount)

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			deductFeeDecorator := feeabsante.NewFeeAbstractionDeductFeeDecorate(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.FeeabsKeeper, suite.app.FeeGrantKeeper)
			anteHandler := sdk.ChainAnteDecorators(deductFeeDecorator)

			_, err = anteHandler(suite.ctx, tx, false)

			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), errors.ErrInsufficientFunds.Error())).To(BeTrue())
		})
	})

	When("ibc fee is sufficient", func() {
		It("should pass with sufficient ibc fee", func() {
			err := suite.app.FeeabsKeeper.SetHostZoneConfig(suite.ctx, mockHostZoneConfig)
			Expect(err).ToNot(HaveOccurred())
			suite.app.FeeabsKeeper.SetTwapRate(suite.ctx, "ibcfee", math.LegacyNewDec(1))
			suite.txBuilder.SetFeeAmount(ibcFeeAmount)

			feeabsAddr := suite.app.FeeabsKeeper.GetFeeAbsModuleAddress()

			// err = suite.mintCoins(feeabsAddr, sdk.NewCoins(feeAmount...))
			err = testutil.FundAccount(suite.ctx, suite.app.BankKeeper, feeabsAddr, feeAmount)
			Expect(err).ToNot(HaveOccurred())

			err = testutil.FundAccount(suite.ctx, suite.app.BankKeeper, testAcc.acc.GetAddress(), ibcFeeAmount)

			// err = suite.mintCoins(testAcc.acc.GetAddress(), sdk.NewCoins(ibcFeeAmount...))
			Expect(err).ToNot(HaveOccurred())

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			deductFeeDecorator := feeabsante.NewFeeAbstractionDeductFeeDecorate(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.FeeabsKeeper, suite.app.FeeGrantKeeper)
			anteHandler := sdk.ChainAnteDecorators(deductFeeDecorator)

			_, err = anteHandler(suite.ctx, tx, false)

			Expect(err).ToNot(HaveOccurred())
		})
	})
})

var _ = Describe("Test Deduct Coins", func() {
	// Create a new AnteTestSuite instance
	s := new(AnteTestSuite)

	BeforeEach(func() {
		err := s.SetupTest(false)                             // Initialize the test environment
		Expect(err).To(BeNil(), "Error on creating test app") // Ensure no error occurred
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()     // Create a new transaction builder
	})

	It("valid coins", func() {
		_, _, addr1 := testdata.KeyTestPubAddr()                       // Generate a test address
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1) // Create a new account
		s.app.AccountKeeper.SetAccount(s.ctx, acc)                     // Set the account in the account keeper

		// Define the initial coins for the account
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(10_000_000_000)))

		// Fund the account with the defined coins
		err := testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil(), "Error on funding account")

		// Check whether DeductCoins work for valid coins
		err = cheqdpost.DeductCoins(s.app.BankKeeper, s.ctx, coins, false)

		// === Validate the results ===
		Expect(err).To(BeNil(), "Error on deducting coins")
	})

	It("valid zero coin", func() {
		_, _, addr1 := testdata.KeyTestPubAddr() // Generate a test address

		// Create a test message and set transaction parameters
		msg := testdata.NewTestMsg(addr1)
		feeAmount := NewTestFeeAmount()              // Define fee amount
		gasLimit := testdata.NewTestGasLimit()       // Define gas limit
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil()) // Set the message in the transaction builder
		s.txBuilder.SetFeeAmount(feeAmount)          // Set the fee amount
		s.txBuilder.SetGasLimit(gasLimit)            // Set the gas limit

		// Create a new account and set it in the account keeper
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)

		// Define the initial coins and fee coins for the account
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(10_000_000_000)))
		feeCoin := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(0))) // Zero fee coin

		// Fund the account with the initial coins
		err := testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil(), "Error on funding account")

		// === Execute the DeductCoins function with zero fee ===
		err = cheqdpost.DeductCoins(s.app.BankKeeper, s.ctx, feeCoin, false)

		Expect(err).To(BeNil(), "Error on deducting zero coin")
	})
})

var _ = Describe("Test Deduct coins and distribute fees", func() {
	// Create a new AnteTestSuite instance
	s := new(AnteTestSuite)

	BeforeEach(func() {
		err := s.SetupTest(false)                             // Initialize the test environment
		Expect(err).To(BeNil(), "Error on creating test app") // Ensure no error occurred
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()     // Create a new transaction builder
	})

	It("valid coins and distribute to feeCollector", func() {
		_, _, addr1 := testdata.KeyTestPubAddr() // Generate a test address

		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1) // Create a new account
		s.app.AccountKeeper.SetAccount(s.ctx, acc)                     // Set the account in the account keeper

		// Define the initial coins for the account
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(10_000_000_000)))

		// Fund the account with the defined coins
		err := testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil(), "Error on funding account")

		// Define the fee to be deducted and the flag for distributing fees
		distributeFees := true
		deductFee := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(5000)))

		// Send the coins to the fee collector module
		// Initially feemarket module has zero funds so we send it module
		// So while performing the deductCoins we don't get error
		err = s.app.BankKeeper.SendCoinsFromAccountToModule(s.ctx, addr1, feemarkettypes.FeeCollectorName, deductFee)
		Expect(err).To(BeNil(), "Error on sending coins to fee collector")

		// Deduct the coins and distribute the fees
		err = cheqdpost.DeductCoins(s.app.BankKeeper, s.ctx, deductFee, distributeFees)
		Expect(err).To(BeNil(), "Error on deducting coins")

		// Check that the fee has been sent to the fee collector
		feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)                   // Get fee collector address
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom) // Get balance

		// Ensure the deducted fee has been correctly transferred to the fee collector
		Expect(feeCollectorBalance.Amount).To(Equal(deductFee.AmountOf(didtypes.BaseMinimalDenom)))
	})

	It("valid zero coin and distribute to feeCollector", func() {
		_, _, addr1 := testdata.KeyTestPubAddr() // Generate a test address

		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1) // Create a new account
		s.app.AccountKeeper.SetAccount(s.ctx, acc)                     // Set the account in the account keeper

		// Define the initial coins for the account
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(10_000_000_000)))

		// Fund the account with the defined coins
		err := testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil(), "Error on funding account")

		// Define the zero fee to be deducted and the flag for distributing fees
		distributeFees := true
		deductFee := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(0))) // Zero fee

		// Send the zero coins to the fee collector module
		err = s.app.BankKeeper.SendCoinsFromAccountToModule(s.ctx, addr1, feemarkettypes.FeeCollectorName, deductFee)
		Expect(err).To(BeNil(), "Error on sending zero coin to fee collector")

		// Deduct the zero coins and distribute the fees
		err = cheqdpost.DeductCoins(s.app.BankKeeper, s.ctx, deductFee, distributeFees)
		Expect(err).To(BeNil(), "Error on deducting zero coin")

		// === Validate the results ===
		// Check that the fee has been sent to the fee collector
		feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)                   // Get fee collector address
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom) // Get balance

		// Ensure the zero fee has been correctly handled
		Expect(feeCollectorBalance.Amount).To(Equal(deductFee.AmountOf(didtypes.BaseMinimalDenom)))
	})
})

var _ = Describe("Test Send tip", func() {
	// Create a new AnteTestSuite instance
	s := new(AnteTestSuite)

	BeforeEach(func() {
		err := s.SetupTest(false)                             // Initialize the test environment
		Expect(err).To(BeNil(), "Error on creating test app") // Ensure no error occurred during setup
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()     // Create a new transaction builder
	})

	It("valid coins", func() {
		_, _, addr1 := testdata.KeyTestPubAddr() // Generate a test address

		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1) // Create a new account
		s.app.AccountKeeper.SetAccount(s.ctx, acc)                     // Set the account in the account keeper

		// Define the initial coins for the account
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(10_000_000_000)))

		// Fund the account with the defined coins
		err := testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil(), "Error on funding account")

		// Define the fee to be deducted and the tip amount
		deductFee := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(5000)))
		tip := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(1000)))

		// Generate another test address for the proposer (recipient of the tip)
		_, _, proposer := testdata.KeyTestPubAddr()

		// Send the fee coins to the fee collector module
		err = s.app.BankKeeper.SendCoinsFromAccountToModule(s.ctx, addr1, feemarkettypes.FeeCollectorName, deductFee)
		Expect(err).To(BeNil(), "Error on sending coins to fee collector")

		// Send the tip to the block proposer
		err = cheqdpost.SendTip(s.app.BankKeeper, s.ctx, proposer, tip)
		Expect(err).To(BeNil(), "Error on sending tip to proposer")

		// Check the balance of the block proposer to confirm receipt of the tip
		balance := s.app.BankKeeper.GetBalance(s.ctx, proposer, didtypes.BaseMinimalDenom)

		// Ensure the tip amount has been credited to the block proposer's account
		Expect(balance.Amount).To(Equal(tip.AmountOf(didtypes.BaseMinimalDenom)))
	})

	It("valid zero coin and distribute to feeCollector", func() {
		_, _, addr1 := testdata.KeyTestPubAddr() // Generate a test address

		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1) // Create a new account
		s.app.AccountKeeper.SetAccount(s.ctx, acc)                     // Set the account in the account keeper

		// Define the initial coins for the account
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(10_000_000_000)))

		// Fund the account with the defined coins
		err := testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil(), "Error on funding account")

		// Define the zero fee to be deducted and the flag for distributing fees
		distributeFees := true
		deductFee := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(0))) // Zero fee

		// Send the zero coins to the fee collector module
		err = s.app.BankKeeper.SendCoinsFromAccountToModule(s.ctx, addr1, feemarkettypes.FeeCollectorName, deductFee)
		Expect(err).To(BeNil(), "Error on sending zero coin to fee collector")

		// Deduct the zero coins and distribute the fees
		err = cheqdpost.DeductCoins(s.app.BankKeeper, s.ctx, deductFee, distributeFees)
		Expect(err).To(BeNil(), "Error on deducting zero coin")

		// Check that the zero fee has been sent to the fee collector
		feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)                   // Get fee collector address
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom) // Get balance

		// Ensure the zero fee has been correctly handled
		Expect(feeCollectorBalance.Amount).To(Equal(deductFee.AmountOf(didtypes.BaseMinimalDenom)))
	})
})

var _ = Describe("Test PostHandle", func() {
	// Create a new AnteTestSuite instance and define the decorators array
	s := new(AnteTestSuite)
	var decorators []sdk.AnteDecorator

	BeforeEach(func() {
		// Initialize the test environment with a complete setup
		err := s.SetupTest(true)
		Expect(err).To(BeNil(), "Error on creating test app")

		// Create a new transaction builder
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

		// Define the AnteDecorators (replacing the default fee deduct decorator with FeeMarketCheckDecorator)
		decorators = []sdk.AnteDecorator{
			feemarketante.NewFeeMarketCheckDecorator(
				s.app.AccountKeeper,
				s.app.BankKeeper,
				s.app.FeeGrantKeeper,
				s.app.FeeMarketKeeper,
				ante.NewDeductFeeDecorator(
					s.app.AccountKeeper,
					s.app.BankKeeper,
					s.app.FeeGrantKeeper,
					nil,
				),
			),
		}
	})

	It("signer has no funds", func() {
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// Prepare a test message and signatures
		msg := testdata.NewTestMsg(addr1)
		feeAmount := NewTestFeeAmount()
		gasLimit := testdata.NewTestGasLimit()

		// Set message, fee, and gas limit in the tx builder
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)

		// Create the transaction
		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// Set account with zero funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(0)))
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil())

		// Create and execute the AnteHandler
		dfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(dfd)

		// Expect an error due to insufficient funds
		_, err = antehandler(s.ctx, tx, false)
		Expect(err).NotTo(BeNil(), "signer has no funds")
	})

	It("signer has no funds --simulate", func() {
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// Prepare a test message and signatures
		msg := testdata.NewTestMsg(addr1)
		feeAmount := NewTestFeeAmount()
		gasLimit := testdata.NewTestGasLimit()

		// Set message, fee, and gas limit in the tx builder
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)

		// Create the transaction
		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// Set account with zero funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(0)))
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil())

		// Create and execute the AnteHandler (in simulation mode)
		dfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(dfd)

		// Expect no error in simulation
		_, err = antehandler(s.ctx, tx, true)
		Expect(err).To(BeNil(), "errored while in simulation when signer has no funds")
	})

	It("0 gas given should fail", func() {
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// Prepare a test message and signatures
		msg := testdata.NewTestMsg(addr1)
		feeAmount := NewTestFeeAmount()

		// Set message, fee, and gas limit (0 gas) in the tx builder
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(0)

		// Create the transaction
		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// Fund the account with sufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(1_00_00_00_00_00_00_000)))
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil())

		// Create and execute the AnteHandler
		dfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(dfd)

		// Expect an error due to 0 gas provided
		_, err = antehandler(s.ctx, tx, false)
		Expect(err).NotTo(BeNil(), "must provide a positive gas")
	})

	It("0 gas given should pass in simulation", func() {
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// Prepare a test message and signatures
		msg := testdata.NewTestMsg(addr1)
		feeAmount := NewTestFeeAmount()

		// Set message, fee, and gas limit (0 gas) in the tx builder
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(0)

		// Create the transaction
		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// Fund the account with sufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(1_00_00_00_00_00_00_000)))
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil())

		// Create and execute the AnteHandler (in simulation mode)
		dfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(dfd)

		// Expect success in simulation with 0 gas
		simulate := true
		_, err = antehandler(s.ctx, tx, simulate)
		Expect(err).To(BeNil(), "should pass in simulation")
	})

	It("signer has enough funds, should pass with tip", func() {
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// Prepare a test message and signatures
		msg := testdata.NewTestMsg(addr1)
		feeAmount := NewTestFeeAmount()
		gasLimit := testdata.NewTestGasLimit()

		// Set message, fee, and gas limit in the tx builder
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)

		// Create the transaction
		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// Fund the account with sufficient funds

		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(1000000000000000)))
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, coins)
		Expect(err).To(BeNil())

		// Create and execute the AnteHandler
		dfd := cheqdante.NewOverAllDecorator(decorators...)
		antehandler := sdk.ChainAnteDecorators(dfd)
		simulate := false
		newCtx, err := antehandler(s.ctx, tx, simulate)
		Expect(err).To(BeNil())

		// Setup block header with proposer
		_, _, proposer := testdata.KeyTestPubAddr()
		s.ctx = newCtx
		blockHeader := s.ctx.BlockHeader()
		blockHeader.ProposerAddress = proposer
		s.ctx = s.ctx.WithBlockHeader(blockHeader)

		// Create and execute the PostHandler with the tax decorator
		taxDecorator := cheqdpost.NewTaxDecorator(
			s.app.AccountKeeper,
			s.app.BankKeeper,
			s.app.FeeGrantKeeper,
			s.app.DidKeeper,
			s.app.ResourceKeeper,
			s.app.FeeMarketKeeper,
			s.app.OracleKeeper,
			s.app.FeeabsKeeper,
		)
		posthandler := sdk.ChainPostDecorators(taxDecorator)

		_, err = posthandler(s.ctx, tx, simulate, true)
		Expect(err).To(BeNil())

		// Validate that the proposer received the tip
		proposerBalance := s.app.BankKeeper.GetAllBalances(s.ctx, proposer)
		Expect(proposerBalance.AmountOf(didtypes.BaseMinimalDenom)).NotTo(BeNil())
	})
})

var _ = Describe("Fee abstraction along with fee market", func() {
	s := new(AnteTestSuite)
	gasLimit := 200000
	ibcFeeAmount := sdk.NewCoins(sdk.NewInt64Coin("ibcfee", 1000*int64(gasLimit)))
	feeAmount := sdk.NewCoins(sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 1000*int64(gasLimit)))

	mockHostZoneConfig := types.HostChainFeeAbsConfig{
		IbcDenom:                "ibcfee",
		OsmosisPoolTokenDenomIn: "osmosis",
		PoolId:                  1,
		Status:                  types.HostChainFeeAbsStatus_UPDATED,
	}

	var feeabsModAcc sdk.ModuleAccountI

	var decorators []sdk.AnteDecorator
	BeforeEach(func() {
		err := s.SetupTest(false)
		Expect(err).To(BeNil(), "Error on creating test app")
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

		decorators = []sdk.AnteDecorator{
			feeabsante.NewFeeAbstrationMempoolFeeDecorator(s.app.FeeabsKeeper),
			feeabsante.NewFeeAbstractionDeductFeeDecorate(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeabsKeeper, s.app.FeeGrantKeeper),
			cheqdante.NewOverAllDecorator(
				feemarketante.NewFeeMarketCheckDecorator(
					// fee market check replaces fee deduct decorator
					s.app.AccountKeeper,
					s.app.BankKeeper,
					s.app.FeeGrantKeeper,
					s.app.FeeMarketKeeper,
					ante.NewDeductFeeDecorator(
						s.app.AccountKeeper,
						s.app.BankKeeper,
						s.app.FeeGrantKeeper,
						nil,
					),
				),
			),
		}

		feeabsModAcc = s.app.FeeabsKeeper.GetFeeAbsModuleAccount(s.ctx)
		s.app.AccountKeeper.SetModuleAccount(s.ctx, feeabsModAcc)

		params, err := s.app.StakingKeeper.GetParams(s.ctx)
		Expect(err).To(BeNil(), "Error getting staking params")
		params.BondDenom = didtypes.BaseMinimalDenom
		err = s.app.StakingKeeper.SetParams(s.ctx, params)
		Expect(err).To(BeNil(), "Error setting the params")
	})

	It("Ensure native tx fee txns are working", func() {
		anteHandler := sdk.ChainAnteDecorators(decorators...)

		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// Prepare a test message and signatures
		msg := testdata.NewTestMsg(addr1)
		gasLimit := testdata.NewTestGasLimit()

		// Set message, fee, and gas limit in the tx builder
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)
		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, feeAmount)
		Expect(err).To(BeNil())

		_, err = anteHandler(s.ctx, tx, true)
		Expect(err).To(BeNil())
	})

	It("Ensure to convert the IBC Denom to native fee", func() {
		err := s.app.FeeabsKeeper.SetHostZoneConfig(s.ctx, mockHostZoneConfig)
		Expect(err).ToNot(HaveOccurred())
		s.app.FeeabsKeeper.SetTwapRate(s.ctx, "ibcfee", math.LegacyNewDec(1))
		minGasPrice := sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 100))...)
		s.ctx = s.ctx.WithMinGasPrices(minGasPrice)

		anteHandler := sdk.ChainAnteDecorators(decorators...)

		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// Prepare a test message and signatures
		msg := testdata.NewTestMsg(addr1)
		gasLimit := testdata.NewTestGasLimit()

		// Set message, fee, and gas limit in the tx builder
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(ibcFeeAmount)
		s.txBuilder.SetGasLimit(gasLimit)
		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, ibcFeeAmount)
		Expect(err).To(BeNil())
		err = testutil.FundModuleAccount(s.ctx, s.app.BankKeeper, types.ModuleName, feeAmount)
		Expect(err).To(BeNil())

		_, err = anteHandler(s.ctx, tx, true)
		Expect(err).To(BeNil())
	})

	It("Ensure to convert the IBC Denom to native fee for taxable txn", func() {
		s.app.OracleKeeper.SetAverage(s.ctx, oraclekeeper.KeyEMA(didtypes.BaseDenom), math.LegacyMustNewDecFromStr("0.016"))
		err := s.app.FeeabsKeeper.SetHostZoneConfig(s.ctx, mockHostZoneConfig)
		Expect(err).ToNot(HaveOccurred())
		ibcDenom := "ibcfee"
		s.app.FeeabsKeeper.SetTwapRate(s.ctx, ibcDenom, math.LegacyNewDec(1))

		anteHandler := sdk.ChainAnteDecorators(decorators...)

		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// msg and signatures
		msg := SandboxDidDoc()
		feeAmount := sdk.NewCoins(sdk.NewCoin(ibcDenom, math.NewInt(75_000_000_000)))
		gasLimit := testdata.NewTestGasLimit()
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)
		s.txBuilder.SetFeePayer(addr1)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// set account with sufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		amount := math.NewInt(90_000_000_000)
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, sdk.NewCoins(sdk.NewCoin(ibcDenom, amount)))
		Expect(err).To(BeNil())
		err = testutil.FundModuleAccount(s.ctx, s.app.BankKeeper, types.ModuleName, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
		Expect(err).To(BeNil())

		taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper, s.app.FeeMarketKeeper, s.app.OracleKeeper, s.app.FeeabsKeeper)
		posthandler := sdk.ChainPostDecorators(taxDecorator)

		// get supply before tx
		supplyBeforeDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		_, err = anteHandler(s.ctx, tx, true)
		Expect(err).To(BeNil())

		_, err = posthandler(s.ctx, tx, false, true)
		Expect(err).To(BeNil(), "Tx errored when fee payer had sufficient funds and provided sufficient fee while subtracting tax on deliverTx")

		// get fee params
		feeParams, err := s.app.DidKeeper.GetParams(s.ctx)
		Expect(err).To(BeNil())

		// check balance of fee payer
		balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, ibcDenom)
		cheqPrice, _ := s.app.OracleKeeper.GetEMA(s.ctx, didtypes.BaseDenom)
		hostconfig, _ := s.app.FeeabsKeeper.GetHostZoneConfig(s.ctx, ibcDenom)
		nativeFee, err := s.app.FeeabsKeeper.CalculateNativeFromIBCCoins(s.ctx, feeAmount, hostconfig)
		Expect(err).To(BeNil())
		twapRate, err := s.app.FeeabsKeeper.GetTwapRate(s.ctx, ibcDenom)
		ibcFee := math.LegacyNewDecFromInt(nativeFee.AmountOf(didtypes.BaseMinimalDenom)).Quo(twapRate).TruncateInt()
		Expect(err).To(BeNil())
		getFee, err := cheqdante.GetFeeForMsg(feeAmount, feeParams.CreateDid, cheqPrice, nativeFee)
		Expect(err).To(BeNil())

		Expect(amount.Sub(ibcFee).Equal(balance.Amount)).To(BeTrue(), "Tax was not subtracted from the fee payer")

		// get supply after tx
		supplyAfterDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// check that supply was deflated
		burnt := cheqdante.GetBurnFeePortion(feeParams.BurnFactor, getFee)
		convertedburnt, err := cheqdpost.ConvertToCheq(burnt, cheqPrice)
		Expect(err).To(BeNil())

		Expect(supplyBeforeDeflation.Sub(supplyAfterDeflation...)).To(Equal(convertedburnt), "Supply was not deflated")

		// check that reward has been sent to the fee collector
		reward := cheqdante.GetRewardPortion(getFee, burnt)

		convertedreward, err := cheqdpost.ConvertToCheq(reward, cheqPrice)
		Expect(err).To(BeNil())

		// calculate oracle share (0.5%)
		oracleShareRate := math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 3) // 0.005 = 0.5%
		rewardAmt := convertedreward.AmountOf(didtypes.BaseMinimalDenom)
		rewardDec := math.LegacyNewDecFromInt(rewardAmt)
		oracleShare := rewardDec.Mul(oracleShareRate).TruncateInt()

		// expected fee collector reward = reward - oracle share
		expectedFeeCollectorReward := rewardAmt.Sub(oracleShare)

		// check fee collector balance
		feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)
		Expect(feeCollectorBalance.Amount).To(Equal(expectedFeeCollectorReward), "Fee Collector did not receive correct reward after oracle share deduction")

		// check oracle module balance
		oracleModule := s.app.AccountKeeper.GetModuleAddress("oracle")
		oracleBalance := s.app.BankKeeper.GetBalance(s.ctx, oracleModule, didtypes.BaseMinimalDenom)
		Expect(oracleBalance.Amount).To(Equal(oracleShare), "Oracle module did not receive the correct reward share")
	})

	It("Ensure to convert the IBC Denom to native fee for non taxable txn", func() {
		err := s.app.FeeabsKeeper.SetHostZoneConfig(s.ctx, mockHostZoneConfig)
		Expect(err).ToNot(HaveOccurred())
		ibcDenom := "ibcfee"
		s.app.FeeabsKeeper.SetTwapRate(s.ctx, ibcDenom, math.LegacyNewDec(1))

		anteHandler := sdk.ChainAnteDecorators(decorators...)

		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// msg and signatures
		msg := testdata.NewTestMsg(addr1)
		feeAmount := sdk.NewCoins(sdk.NewCoin(ibcDenom, math.NewInt(50_000_000_000)))
		gasLimit := testdata.NewTestGasLimit()
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)
		s.txBuilder.SetFeePayer(addr1)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// set account with sufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		amount := math.NewInt(50_000_000_000)
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, sdk.NewCoins(sdk.NewCoin(ibcDenom, amount)))
		Expect(err).To(BeNil())

		err = testutil.FundModuleAccount(s.ctx, s.app.BankKeeper, types.ModuleName, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
		Expect(err).To(BeNil())

		taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper, s.app.FeeMarketKeeper, s.app.OracleKeeper, s.app.FeeabsKeeper)
		posthandler := sdk.ChainPostDecorators(taxDecorator)

		// get supply before tx
		supplyBefore, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		newCtx, err := anteHandler(s.ctx, tx, true)
		Expect(err).To(BeNil())

		_, _, proposer := testdata.KeyTestPubAddr()
		s.ctx = newCtx
		a := s.ctx.BlockHeader()
		a.ProposerAddress = proposer
		newCtx = s.ctx.WithBlockHeader(a)
		s.ctx = newCtx

		_, err = posthandler(s.ctx, tx, false, true)
		Expect(err).To(BeNil(), "Tx errored when fee payer had sufficient funds and provided sufficient fee while subtracting tax on deliverTx")

		// check balance of fee payer
		balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, ibcDenom)
		Expect(amount.Sub(math.NewInt(feeAmount.AmountOf(ibcDenom).Int64())).Equal(balance.Amount)).To(BeTrue(), "Fee amount subtracted was not equal to fee amount required for non-taxable tx")

		// get supply after tx
		supplyAfter, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// check that supply was not deflated
		Expect(supplyBefore).To(Equal(supplyAfter), "Supply was deflated")

		// check that reward has been sent to the fee collector
		feeCollector := s.app.AccountKeeper.GetModuleAddress(feemarkettypes.FeeCollectorName)
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)

		Expect((feeCollectorBalance.Amount).GT(math.NewInt(0)))
	})

	It("Ensure taxable txn working fine after integrating the fee-abs", func() {
		anteHandler := sdk.ChainAnteDecorators(decorators...)

		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// msg and signatures
		msg := SandboxDidDoc()
		feeAmount := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(50_000_000_000)))
		gasLimit := testdata.NewTestGasLimit()
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())
		s.txBuilder.SetFeeAmount(feeAmount)
		s.txBuilder.SetGasLimit(gasLimit)
		s.txBuilder.SetFeePayer(addr1)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// set account with sufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		amount := math.NewInt(50_000_000_000)
		err = testutil.FundAccount(s.ctx, s.app.BankKeeper, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
		Expect(err).To(BeNil())

		taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper, s.app.FeeMarketKeeper, s.app.OracleKeeper, s.app.FeeabsKeeper)
		posthandler := sdk.ChainPostDecorators(taxDecorator)

		// get supply before tx
		supplyBeforeDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		_, err = anteHandler(s.ctx, tx, true)
		Expect(err).To(BeNil())

		_, err = posthandler(s.ctx, tx, false, true)
		Expect(err).To(BeNil(), "Tx errored when fee payer had sufficient funds and provided sufficient fee while subtracting tax on deliverTx")

		// get fee params
		feeParams, err := s.app.DidKeeper.GetParams(s.ctx)
		Expect(err).To(BeNil())

		// check balance of fee payer
		balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, didtypes.BaseMinimalDenom)
		Expect(amount.Sub(*feeParams.CreateDid[0].MinAmount).Equal(balance.Amount)).To(BeTrue(), "Tax was not subtracted from the fee payer")

		// get supply after tx
		supplyAfterDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// check that supply was deflated
		burnt := cheqdante.GetBurnFeePortion(feeParams.BurnFactor, sdk.NewCoins(sdk.NewCoin(feeParams.CreateDid[0].Denom, *feeParams.CreateDid[0].MinAmount)))
		Expect(supplyBeforeDeflation.Sub(supplyAfterDeflation...)).To(Equal(burnt), "Supply was not deflated")

		// check that reward has been sent to the fee collector
		reward := cheqdante.GetRewardPortion(sdk.NewCoins(sdk.NewCoin(feeParams.CreateDid[0].Denom, *feeParams.CreateDid[0].MinAmount)), burnt)

		// calculate oracle share (0.5%)
		oracleShareRate := math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 3) // 0.005 = 0.5%
		rewardAmt := reward.AmountOf(didtypes.BaseMinimalDenom)
		rewardDec := math.LegacyNewDecFromInt(rewardAmt)
		oracleShare := rewardDec.Mul(oracleShareRate).TruncateInt()

		// expected fee collector reward = reward - oracle share
		expectedFeeCollectorReward := rewardAmt.Sub(oracleShare)

		// check fee collector balance
		feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)
		Expect(feeCollectorBalance.Amount).To(Equal(expectedFeeCollectorReward), "Fee Collector did not receive correct reward after oracle share deduction")

		// check oracle module balance
		oracleModule := s.app.AccountKeeper.GetModuleAddress("oracle")
		oracleBalance := s.app.BankKeeper.GetBalance(s.ctx, oracleModule, didtypes.BaseMinimalDenom)
		Expect(oracleBalance.Amount).To(Equal(oracleShare), "Oracle module did not receive the correct reward share")
	})
})
