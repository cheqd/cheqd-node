package ante_test

import (
	cheqdante "github.com/cheqd/cheqd-node/ante"
	cheqdpost "github.com/cheqd/cheqd-node/post"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fee tests", func() {
	s := new(AnteTestSuite)

	It("Test Ensure Zero Mempool Fees On Simulation CheckTx", func() {
		s.SetupTest(true) // setup
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

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

	It("Test Ensure Mempool Fees On CheckTx", func() {
		s.SetupTest(true) // setup
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

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

	It("Test Deduct Fees On Deliver Tx", func() {
		s.SetupTest(false) // setup
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

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

	It("Test TaxableTx Mempool Inclusion On CheckTx", func() {
		s.SetupTest(true) // setup
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

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

	It("Test TaxableTx Lifecycle On DeliverTx", func() {
		s.SetupTest(false) // setup
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

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
		posthandler := sdk.ChainAnteDecorators(taxDecorator)

		// get supply before tx
		supplyBeforeDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// antehandler should not error to replay in mempool or deliverTx
		_, err = antehandler(s.ctx, tx, false)
		Expect(err).To(BeNil(), "Tx errored when taxable on deliverTx")

		_, err = posthandler(s.ctx, tx, false)
		Expect(err).To(BeNil(), "Tx errored when fee payer had sufficient funds and provided sufficient fee while subtracting tax on deliverTx")

		// get fee params
		feeParams := s.app.DidKeeper.GetParams(s.ctx)

		// check balance of fee payer
		balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, didtypes.BaseMinimalDenom)
		createDidTax := feeParams.TxTypes[didtypes.DefaultKeyCreateDid]
		Expect(amount.Sub(sdk.NewInt(createDidTax.Amount.Int64()))).To(Equal(balance.Amount), "Tax was not subtracted from the fee payer")

		// get supply after tx
		supplyAfterDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// check that supply was deflated
		burnt := cheqdante.GetBurnFeePortion(feeParams.BurnFactor, sdk.NewCoins(createDidTax))
		Expect(supplyBeforeDeflation.Sub(supplyAfterDeflation...)).To(Equal(burnt), "Supply was not deflated")

		// check that reward has been sent to the fee collector
		reward := cheqdante.GetRewardPortion(sdk.NewCoins(createDidTax), burnt)
		feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)

		Expect(feeCollectorBalance.Amount).To(Equal(reward.AmountOf(didtypes.BaseMinimalDenom)), "Reward was not sent to the fee collector")
	})

	It("Test Non TaxableTx Lifecycle on DeliverTx", func() {
		s.SetupTest(false) // setup
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

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
		posthandler := sdk.ChainAnteDecorators(taxDecorator)

		// get supply before tx
		supplyBefore, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// antehandler should not error since we make sure that the fee is sufficient in DeliverTx (simulate=false only, Posthandler will check it otherwise)
		_, err = antehandler(s.ctx, tx, false)
		Expect(err).To(BeNil(), "Tx errored when non-taxable on deliverTx")

		_, err = posthandler(s.ctx, tx, false)
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

	It("Test TaxableTx Lifecycle on DeliveTx Simulation", func() {
		s.SetupTest(false) // setup
		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

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
		posthandler := sdk.ChainAnteDecorators(taxDecorator)

		// get supply before tx
		supplyBefore, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
		Expect(err).To(BeNil())

		// antehandler should not error since we make sure that the fee is sufficient in DeliverTx (simulate=false only, Posthandler will check it otherwise)
		_, err = antehandler(s.ctx, tx, true)
		Expect(err).To(BeNil(), "Tx errored when taxable on deliverTx simulation mode")

		_, err = posthandler(s.ctx, tx, true)
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
