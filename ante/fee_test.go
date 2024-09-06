package ante_test

import (
	"strings"

	cheqdante "github.com/cheqd/cheqd-node/ante"
	cheqdpost "github.com/cheqd/cheqd-node/post"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/osmosis-labs/fee-abstraction/v7/x/feeabs/ante"
	"github.com/osmosis-labs/fee-abstraction/v7/x/feeabs/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fee tests on CheckTx", func() {
	s := new(AnteTestSuite)

	BeforeEach(func() {
		err := s.SetupTest(true) // setup
		Expect(err).To(BeNil(), "Error on creating test app")
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()
	})

	It("Ensure Zero Mempool Fees On Simulation", func() {
		mfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, nil)
		antehandler := sdk.ChainAnteDecorators(mfd)

		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(300000000000)))
		err := testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, coins)
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
		mfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, nil)
		antehandler := sdk.ChainAnteDecorators(mfd)

		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(300_000_000_000)))
		err := testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, coins)
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
		ncheqPrice := sdk.NewDecCoinFromDec(didtypes.BaseMinimalDenom, sdk.NewDec(20_000_000_000))
		highGasPrice := []sdk.DecCoin{ncheqPrice}
		s.ctx = s.ctx.WithMinGasPrices(highGasPrice)

		// Set IsCheckTx to true
		s.ctx = s.ctx.WithIsCheckTx(true)

		// antehandler errors with insufficient fees
		_, err = antehandler(s.ctx, tx, false)
		Expect(err).NotTo(BeNil(), "Decorator should have errored on too low fee for local gasPrice")

		// antehandler should not error since we do not check minGasPrice in simulation mode
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

		ncheqPrice = sdk.NewDecCoinFromDec(didtypes.BaseMinimalDenom, sdk.NewDec(0).Quo(sdk.NewDec(didtypes.BaseFactor)))
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
		feeAmount := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(5_000_000_000)))
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
		err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(100_000_000_000))))
		Expect(err).To(BeNil())

		dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, nil, nil)
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

	BeforeEach(func() {
		err := s.SetupTest(false) // setup
		Expect(err).To(BeNil(), "Error on creating test app")
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()
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
		coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(10_000_000_000)))
		//nocheck:errcheck
		err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, coins)
		Expect(err).To(BeNil())

		dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, nil, nil)
		antehandler := sdk.ChainAnteDecorators(dfd)

		_, err = antehandler(s.ctx, tx, false)

		Expect(err).NotTo(BeNil(), "Tx did not error when fee payer had insufficient funds")

		// Set account with sufficient funds
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(200_000_000_000))))
		Expect(err).To(BeNil())

		_, err = antehandler(s.ctx, tx, false)

		Expect(err).To(BeNil(), "Tx errored after account has been set with sufficient funds")
	})

	It("TaxableTx Lifecycle", func() {
		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// msg and signatures
		msg := SandboxDidDoc()
		feeAmount := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(5_000_000_000)))
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
		amount := sdk.NewInt(100_000_000_000)
		err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
		Expect(err).To(BeNil())

		dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, nil, nil)
		antehandler := sdk.ChainAnteDecorators(dfd)

		taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper)
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
		feeParams := s.app.DidKeeper.GetParams(s.ctx)

		// check balance of fee payer
		balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, didtypes.BaseMinimalDenom)
		Expect(amount.Sub(sdk.NewInt(feeParams.CreateDid.Amount.Int64()))).To(Equal(balance.Amount), "Tax was not subtracted from the fee payer")

		// get supply after tx
		supplyAfterDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// check that supply was deflated
		burnt := cheqdante.GetBurnFeePortion(feeParams.BurnFactor, sdk.NewCoins(feeParams.CreateDid))
		Expect(supplyBeforeDeflation.Sub(supplyAfterDeflation...)).To(Equal(burnt), "Supply was not deflated")

		// check that reward has been sent to the fee collector
		reward := cheqdante.GetRewardPortion(sdk.NewCoins(feeParams.CreateDid), burnt)
		feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)

		Expect(feeCollectorBalance.Amount).To(Equal(reward.AmountOf(didtypes.BaseMinimalDenom)), "Reward was not sent to the fee collector")
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
		amount := sdk.NewInt(300_000_000_000)
		err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
		Expect(err).To(BeNil())

		dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, nil, nil)
		antehandler := sdk.ChainAnteDecorators(dfd)

		taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper)
		posthandler := sdk.ChainPostDecorators(taxDecorator)

		// get supply before tx
		supplyBefore, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// antehandler should not error since we make sure that the fee is sufficient in DeliverTx (simulate=false only, Posthandler will check it otherwise)
		_, err = antehandler(s.ctx, tx, false)
		Expect(err).To(BeNil(), "Tx errored when non-taxable on deliverTx")

		_, err = posthandler(s.ctx, tx, false, true)
		Expect(err).To(BeNil(), "Tx errored when non-taxable on deliverTx from posthandler")

		// check balance of fee payer
		balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, didtypes.BaseMinimalDenom)
		Expect(amount.Sub(sdk.NewInt(feeAmount.AmountOf(didtypes.BaseMinimalDenom).Int64()))).To(Equal(balance.Amount), "Fee amount subtracted was not equal to fee amount required for non-taxable tx")

		// get supply after tx
		supplyAfter, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// check that supply was not deflated
		Expect(supplyBefore).To(Equal(supplyAfter), "Supply was deflated")

		// check that reward has been sent to the fee collector
		feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)

		Expect(feeCollectorBalance.Amount).To(Equal(feeAmount.AmountOf(didtypes.BaseMinimalDenom)), "Fee was not sent to the fee collector")
	})

	It("TaxableTx Lifecycle on Simulation", func() {
		// keys and addresses
		priv1, _, addr1 := testdata.KeyTestPubAddr()

		// msg and signatures
		msg := testdata.NewTestMsg(addr1)
		Expect(s.txBuilder.SetMsgs(msg)).To(BeNil())

		// set zero gas
		s.txBuilder.SetGasLimit(0)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
		tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
		Expect(err).To(BeNil())

		// set account with sufficient funds
		acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
		s.app.AccountKeeper.SetAccount(s.ctx, acc)
		amount := sdk.NewInt(100_000_000_000)
		err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
		Expect(err).To(BeNil())

		dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, nil, nil)
		antehandler := sdk.ChainAnteDecorators(dfd)

		taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper)
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

		Expect(feeCollectorBalance.Amount).To(Equal(sdk.NewInt(0)), "Reward was sent to the fee collector when taxable tx simulation mode")
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
			MinSwapAmount:           0,
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
			minGasPrice := sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.NewInt64Coin("ncheq", 100))...)
			suite.txBuilder.SetGasLimit(gasLimit)
			suite.txBuilder.SetFeeAmount(feeAmount)
			suite.ctx = suite.ctx.WithMinGasPrices(minGasPrice)

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			mempoolDecorator := ante.NewFeeAbstrationMempoolFeeDecorator(suite.app.FeeabsKeeper)
			anteHandler := sdk.ChainAnteDecorators(mempoolDecorator)

			_, err := anteHandler(suite.ctx, tx, false)

			// Expect error due to insufficient fee
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), errors.ErrInsufficientFee.Error())).To(BeTrue())
		})

		It("should fail with insufficient native fee", func() {
			feeAmount := sdk.NewCoins(sdk.NewInt64Coin("ncheq", 100))
			minGasPrice := sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.NewInt64Coin("ncheq", 1000))...)
			suite.txBuilder.SetGasLimit(gasLimit)
			suite.txBuilder.SetFeeAmount(feeAmount)
			suite.ctx = suite.ctx.WithMinGasPrices(minGasPrice)

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			mempoolDecorator := ante.NewFeeAbstrationMempoolFeeDecorator(suite.app.FeeabsKeeper)
			anteHandler := sdk.ChainAnteDecorators(mempoolDecorator)

			_, err := anteHandler(suite.ctx, tx, false)

			// Expect error due to insufficient fee
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), errors.ErrInsufficientFee.Error())).To(BeTrue())
		})

		It("should pass with sufficient native fee", func() {
			feeAmount := sdk.NewCoins(sdk.NewInt64Coin("ncheq", 1000*int64(gasLimit)))
			minGasPrice := sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.NewInt64Coin("ncheq", 1000))...)
			suite.txBuilder.SetGasLimit(gasLimit)
			suite.txBuilder.SetFeeAmount(feeAmount)
			suite.ctx = suite.ctx.WithMinGasPrices(minGasPrice)

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			mempoolDecorator := ante.NewFeeAbstrationMempoolFeeDecorator(suite.app.FeeabsKeeper)
			anteHandler := sdk.ChainAnteDecorators(mempoolDecorator)

			_, err := anteHandler(suite.ctx, tx, false)

			// No error is expected
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with unknown ibc fee denom", func() {
			feeAmount := sdk.NewCoins(sdk.NewInt64Coin("ibcfee", 1000*int64(gasLimit)))
			minGasPrice := sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.NewInt64Coin("ncheq", 1000))...)
			suite.txBuilder.SetGasLimit(gasLimit)
			suite.txBuilder.SetFeeAmount(feeAmount)
			suite.ctx = suite.ctx.WithMinGasPrices(minGasPrice)

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			mempoolDecorator := ante.NewFeeAbstrationMempoolFeeDecorator(suite.app.FeeabsKeeper)
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
			suite.app.FeeabsKeeper.SetTwapRate(suite.ctx, "ibcfee", sdk.NewDec(1))

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			mempoolDecorator := ante.NewFeeAbstrationMempoolFeeDecorator(suite.app.FeeabsKeeper)
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
		minGasPrice = sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.NewInt64Coin("ncheq", 1000))...)
		feeAmount = sdk.NewCoins(sdk.NewInt64Coin("ncheq", 1000*int64(gasLimit)))
		ibcFeeAmount = sdk.NewCoins(sdk.NewInt64Coin("ibcfee", 1000*int64(gasLimit)))

		mockHostZoneConfig = types.HostChainFeeAbsConfig{
			IbcDenom:                "ibcfee",
			OsmosisPoolTokenDenomIn: "osmosis",
			PoolId:                  1,
			Status:                  types.HostChainFeeAbsStatus_UPDATED,
			MinSwapAmount:           0,
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

		params := suite.app.StakingKeeper.GetParams(suite.ctx)
		params.BondDenom = "ncheq"
		suite.app.StakingKeeper.SetParams(suite.ctx, params)

		// this line will create the module account
		_ = suite.app.AccountKeeper.GetModuleAccount(suite.ctx, types.ModuleName)

	})

	Describe("Handling Deduct Fee Decorator", func() {
		When("native fee is insufficient", func() {
			It("should fail due to insufficient native fee", func() {
				suite.app.FeeabsKeeper.SetTwapRate(suite.ctx, "ibcfee", sdk.NewDec(1))
				_, _, addr := testdata.KeyTestPubAddr()
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr)
				err := acc.SetAccountNumber(1)
				Expect(err).ToNot(HaveOccurred())

				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

				suite.txBuilder.SetFeePayer(acc.GetAddress())
				// Construct and run the ante handler
				tx := suite.txBuilder.GetTx()
				deductFeeDecorator := ante.NewFeeAbstractionDeductFeeDecorate(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.FeeabsKeeper, suite.app.FeeGrantKeeper)
				anteHandler := sdk.ChainAnteDecorators(deductFeeDecorator)

				_, err = anteHandler(suite.ctx, tx, false)

				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), errors.ErrInsufficientFunds.Error())).To(BeTrue())
			})
		})
	})

	When("native fee is sufficient", func() {
		It("should pass with sufficient native fee", func() {
			suite.app.FeeabsKeeper.SetTwapRate(suite.ctx, "ibcfee", sdk.NewDec(1))

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			deductFeeDecorator := ante.NewFeeAbstractionDeductFeeDecorate(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.FeeabsKeeper, suite.app.FeeGrantKeeper)
			anteHandler := sdk.ChainAnteDecorators(deductFeeDecorator)

			_, err := anteHandler(suite.ctx, tx, false)

			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("ibc fee is insufficient", func() {
		It("should fail due to insufficient ibc fee", func() {
			err := suite.app.FeeabsKeeper.SetHostZoneConfig(suite.ctx, mockHostZoneConfig)
			Expect(err).ToNot(HaveOccurred())
			suite.app.FeeabsKeeper.SetTwapRate(suite.ctx, "ibcfee", sdk.NewDec(1))

			suite.txBuilder.SetFeeAmount(ibcFeeAmount)

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			deductFeeDecorator := ante.NewFeeAbstractionDeductFeeDecorate(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.FeeabsKeeper, suite.app.FeeGrantKeeper)
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
			suite.app.FeeabsKeeper.SetTwapRate(suite.ctx, "ibcfee", sdk.NewDec(1))
			suite.txBuilder.SetFeeAmount(ibcFeeAmount)

			feeabsAddr := suite.app.FeeabsKeeper.GetFeeAbsModuleAddress()

			err = suite.mintCoins(feeabsAddr, sdk.NewCoins(feeAmount...))
			Expect(err).ToNot(HaveOccurred())

			err = suite.mintCoins(testAcc.acc.GetAddress(), sdk.NewCoins(ibcFeeAmount...))
			Expect(err).ToNot(HaveOccurred())

			// Construct and run the ante handler
			tx := suite.txBuilder.GetTx()
			deductFeeDecorator := ante.NewFeeAbstractionDeductFeeDecorate(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.FeeabsKeeper, suite.app.FeeGrantKeeper)
			anteHandler := sdk.ChainAnteDecorators(deductFeeDecorator)

			_, err = anteHandler(suite.ctx, tx, false)

			Expect(err).ToNot(HaveOccurred())
		})
	})
})

func (s *AnteTestSuite) mintCoins(addr sdk.AccAddress, someCoins sdk.Coins) error {
	err := s.app.BankKeeper.MintCoins(s.ctx, minttypes.ModuleName, someCoins)
	if err != nil {
		return err
	}

	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, addr, someCoins)
	if err != nil {
		return err
	}

	return nil
}
