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
)

func (s *AnteTestSuite) TestDeductFeeDecorator_ZeroGas() {
	s.SetupTest(true) // setup
	s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

	mfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(300000000000)))
	err := testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, coins)
	s.Require().NoError(err)

	// msg and signatures
	msg := testdata.NewTestMsg(addr1)
	s.Require().NoError(s.txBuilder.SetMsgs(msg))

	// set zero gas
	s.txBuilder.SetGasLimit(0)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
	s.Require().NoError(err)

	// Set IsCheckTx to true
	s.ctx = s.ctx.WithIsCheckTx(true)

	_, err = antehandler(s.ctx, tx, false)
	s.Require().Error(err)

	// zero gas is accepted in simulation mode
	_, err = antehandler(s.ctx, tx, true)
	s.Require().NoError(err)
}

func (s *AnteTestSuite) TestEnsureMempoolFees() {
	s.SetupTest(true) // setup
	s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

	mfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(300_000_000_000)))
	err := testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, coins)
	s.Require().NoError(err)

	// msg and signatures
	msg := testdata.NewTestMsg(addr1)
	feeAmount := NewTestFeeAmount()
	gasLimit := uint64(15)
	s.Require().NoError(s.txBuilder.SetMsgs(msg))
	s.txBuilder.SetFeeAmount(feeAmount)
	s.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
	s.Require().NoError(err)

	// Set high gas price so standard test fee fails
	ncheqPrice := sdk.NewDecCoinFromDec(didtypes.BaseMinimalDenom, sdk.NewDec(20_000_000_000))
	highGasPrice := []sdk.DecCoin{ncheqPrice}
	s.ctx = s.ctx.WithMinGasPrices(highGasPrice)

	// Set IsCheckTx to true
	s.ctx = s.ctx.WithIsCheckTx(true)

	// antehandler errors with insufficient fees
	_, err = antehandler(s.ctx, tx, false)
	s.Require().NotNil(err, "Decorator should have errored on too low fee for local gasPrice")

	// antehandler should not error since we do not check minGasPrice in simulation mode
	cacheCtx, _ := s.ctx.CacheContext()
	_, err = antehandler(cacheCtx, tx, true)
	s.Require().Nil(err, "Decorator should not have errored in simulation mode")

	// Set IsCheckTx to false
	s.ctx = s.ctx.WithIsCheckTx(false)

	// antehandler should not error since we do not check minGasPrice in DeliverTx
	_, err = antehandler(s.ctx, tx, false)
	s.Require().Nil(err, "MempoolFeeDecorator returned error in DeliverTx")

	// Set IsCheckTx back to true for testing sufficient mempool fee
	s.ctx = s.ctx.WithIsCheckTx(true)

	ncheqPrice = sdk.NewDecCoinFromDec(didtypes.BaseMinimalDenom, sdk.NewDec(0).Quo(sdk.NewDec(didtypes.BaseFactor)))
	lowGasPrice := []sdk.DecCoin{ncheqPrice}
	s.ctx = s.ctx.WithMinGasPrices(lowGasPrice) // 1 ncheq

	newCtx, err := antehandler(s.ctx, tx, false)
	s.Require().Nil(err, "Decorator should not have errored on fee higher than local gasPrice")
	// Priority is the smallest gas price amount in any denom. Since we have only 1 gas price
	// of 10000000000ncheq, the priority here is 10*10^9.
	s.Require().Equal(int64(10)*didtypes.BaseFactor, newCtx.Priority())
}

func (s *AnteTestSuite) TestDeductFees() {
	s.SetupTest(false) // setup
	s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	msg := testdata.NewTestMsg(addr1)
	feeAmount := NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	s.Require().NoError(s.txBuilder.SetMsgs(msg))
	s.txBuilder.SetFeeAmount(feeAmount)
	s.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
	s.Require().NoError(err)

	// Set account with insufficient funds
	acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	coins := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(10_000_000_000)))
	//nocheck:errcheck
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, coins)
	s.Require().NoError(err)

	dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, nil, s.app.DidKeeper, s.app.ResourceKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(dfd)

	_, err = antehandler(s.ctx, tx, false)

	s.Require().NotNil(err, "Tx did not error when fee payer had insufficient funds")

	// Set account with sufficient funds
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(200_000_000_000))))
	s.Require().NoError(err)

	_, err = antehandler(s.ctx, tx, false)

	s.Require().Nil(err, "Tx errored after account has been set with sufficient funds")
}

func (s *AnteTestSuite) TestTaxableTxMempoolExcluded() {
	s.SetupTest(true) // setup
	s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	msg := SandboxDidDoc()
	feeAmount := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(1_000_000_000)))
	gasLimit := testdata.NewTestGasLimit()
	s.Require().NoError(s.txBuilder.SetMsgs(msg))
	s.txBuilder.SetFeeAmount(feeAmount)
	s.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
	s.Require().NoError(err)

	// set account with sufficient funds
	acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(10_000_000_000))))
	s.Require().NoError(err)

	dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, nil, s.app.DidKeeper, s.app.ResourceKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(dfd)

	_, err = antehandler(s.ctx, tx, false)

	s.Require().NotNil(err, "Tx did not error when fee payer had sufficient funds and provided lower fee than required")

	// set checkTx to false
	s.ctx = s.ctx.WithIsCheckTx(false)

	// antehandler should error since we make sure that the fee is sufficient in DeliverTx (simulate=false only, Posthandler will check it otherwise)
	_, err = antehandler(s.ctx, tx, false)

	s.Require().NotNil(err, "Tx did not error when fee payer had sufficient funds and provided lower fee than required when checkTx is false")
}

func (s *AnteTestSuite) TestTaxableTxMempoolIncluded() {
	s.SetupTest(true) // setup
	s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	msg := SandboxDidDoc()
	feeAmount := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(5_000_000_000)))
	gasLimit := testdata.NewTestGasLimit()
	s.Require().NoError(s.txBuilder.SetMsgs(msg))
	s.txBuilder.SetFeeAmount(feeAmount)
	s.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
	s.Require().NoError(err)

	// set account with sufficient funds
	acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(100_000_000_000))))
	s.Require().NoError(err)

	dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, nil, s.app.DidKeeper, s.app.ResourceKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(dfd)

	_, err = antehandler(s.ctx, tx, false)

	s.Require().Nil(err, "Tx errored when fee payer had sufficient funds and provided sufficient fee")

	// set checkTx to false
	s.ctx = s.ctx.WithIsCheckTx(false)

	// antehandler should not error since we make sure that the fee is sufficient in DeliverTx (simulate=false only, Posthandler will check it otherwise)
	_, err = antehandler(s.ctx, tx, false)

	s.Require().Nil(err, "Tx errored when fee payer had sufficient funds and provided sufficient fee when checkTx is false")
}

func (s *AnteTestSuite) TestTaxableTxOverallLifecycleNonSimulated() {
	s.SetupTest(false) // setup
	s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	msg := SandboxDidDoc()
	feeAmount := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(5_000_000_000)))
	gasLimit := testdata.NewTestGasLimit()
	s.Require().NoError(s.txBuilder.SetMsgs(msg))
	s.txBuilder.SetFeeAmount(feeAmount)
	s.txBuilder.SetGasLimit(gasLimit)
	s.txBuilder.SetFeePayer(addr1)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
	s.Require().NoError(err)

	// set account with sufficient funds
	acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	amount := sdk.NewInt(100_000_000_000)
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
	s.Require().NoError(err)

	dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, nil, s.app.DidKeeper, s.app.ResourceKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(dfd)

	taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper)
	posthandler := sdk.ChainAnteDecorators(taxDecorator)

	// get supply before tx
	supplyBeforeDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
	s.Require().NoError(err)

	// antehandler should not error since we make sure that the fee is sufficient in DeliverTx (simulate=false only, Posthandler will check it otherwise)
	_, err = antehandler(s.ctx, tx, false)
	s.Require().Nil(err, "Tx errored when fee payer had sufficient funds and provided sufficient fee")

	_, err = posthandler(s.ctx, tx, false)
	s.Require().Nil(err, "Tx errored when fee payer had sufficient funds and provided sufficient fee while subtracting tax")

	// get fee params
	feeParams := s.app.DidKeeper.GetParams(s.ctx)

	// check balance of fee payer
	balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, didtypes.BaseMinimalDenom)
	createDidTax := feeParams.TxTypes[didtypes.DefaultKeyCreateDid]
	s.Require().Equal(amount.Sub(sdk.NewInt(createDidTax.Amount.Int64())), balance.Amount, "Tax was not subtracted from the fee payer")

	// get supply after tx
	supplyAfterDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
	s.Require().NoError(err)

	// check that supply was deflated
	burnt := cheqdante.GetBurnFeePortion(s.ctx, feeParams.BurnFactor, sdk.NewCoins(createDidTax))
	s.Require().Equal(supplyBeforeDeflation.Sub(supplyAfterDeflation...), burnt, "Supply was not deflated")

	// check that reward has been sent to the fee collector
	reward := cheqdante.GetRewardPortion(sdk.NewCoins(createDidTax), burnt)
	feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)

	s.Require().Equal(feeCollectorBalance.Amount, reward.AmountOf(didtypes.BaseMinimalDenom), "Reward was not sent to the fee collector")
}

func (s *AnteTestSuite) TestTaxableTxOverallLifecycleSimulated() {
	s.SetupTest(false) // setup
	s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	msg := SandboxDidDoc()
	feeAmount := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(5_000_000_000)))
	gasLimit := testdata.NewTestGasLimit()
	s.Require().NoError(s.txBuilder.SetMsgs(msg))
	s.txBuilder.SetFeeAmount(feeAmount)
	s.txBuilder.SetGasLimit(gasLimit)
	s.txBuilder.SetFeePayer(addr1)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
	s.Require().NoError(err)

	// set account with sufficient funds
	acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	amount := sdk.NewInt(100_000_000_000)
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
	s.Require().NoError(err)

	dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, nil, s.app.DidKeeper, s.app.ResourceKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(dfd)

	taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper)
	posthandler := sdk.ChainAnteDecorators(taxDecorator)

	// get supply before tx
	supplyBeforeDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
	s.Require().NoError(err)

	// antehandler should not error since we make sure that the fee is sufficient in DeliverTx (simulate=false only, Posthandler will check it otherwise)
	_, err = antehandler(s.ctx, tx, true)
	s.Require().Nil(err, "Tx errored when fee payer had sufficient funds and provided sufficient fee")

	_, err = posthandler(s.ctx, tx, true)
	s.Require().Nil(err, "Tx errored when fee payer had sufficient funds and provided sufficient fee while subtracting tax")

	// get fee params
	feeParams := s.app.DidKeeper.GetParams(s.ctx)

	// check balance of fee payer
	balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, didtypes.BaseMinimalDenom)
	createDidTax := feeParams.TxTypes[didtypes.DefaultKeyCreateDid]
	s.Require().Equal(amount.Sub(sdk.NewInt(createDidTax.Amount.Int64())), balance.Amount, "Tax was not subtracted from the fee payer")

	// get supply after tx
	supplyAfterDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
	s.Require().NoError(err)

	// check that supply was deflated
	burnt := cheqdante.GetBurnFeePortion(s.ctx, feeParams.BurnFactor, sdk.NewCoins(createDidTax))
	s.Require().Equal(supplyBeforeDeflation.Sub(supplyAfterDeflation...), burnt, "Supply was not deflated")

	// check that reward has been sent to the fee collector
	reward := cheqdante.GetRewardPortion(sdk.NewCoins(createDidTax), burnt)
	feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)

	s.Require().Equal(feeCollectorBalance.Amount, reward.AmountOf(didtypes.BaseMinimalDenom), "Reward was not sent to the fee collector")
}

func (s *AnteTestSuite) TestNonTaxableTxOverallLifecycleSimulated() {
	s.SetupTest(false) // setup
	s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	msg := testdata.NewTestMsg(addr1)
	feeAmount := NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	s.Require().NoError(s.txBuilder.SetMsgs(msg))
	s.txBuilder.SetFeeAmount(feeAmount)
	s.txBuilder.SetGasLimit(gasLimit)
	s.txBuilder.SetFeePayer(addr1)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
	s.Require().NoError(err)

	// set account with sufficient funds
	acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	amount := sdk.NewInt(300_000_000_000)
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, amount)))
	s.Require().NoError(err)

	dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, nil, s.app.DidKeeper, s.app.ResourceKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(dfd)

	taxDecorator := cheqdpost.NewTaxDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, s.app.DidKeeper, s.app.ResourceKeeper)
	posthandler := sdk.ChainAnteDecorators(taxDecorator)

	// get supply before tx
	supplyBeforeDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
	s.Require().NoError(err)

	// antehandler should not error since we make sure that the fee is sufficient in DeliverTx (simulate=false only, Posthandler will check it otherwise)
	_, err = antehandler(s.ctx, tx, true)
	s.Require().Nil(err, "Tx errored when fee payer had sufficient funds and provided sufficient fee")

	_, err = posthandler(s.ctx, tx, true)
	s.Require().Nil(err, "Tx errored when fee payer had sufficient funds and provided sufficient fee while subtracting tax")

	// check balance of fee payer
	balance := s.app.BankKeeper.GetBalance(s.ctx, addr1, didtypes.BaseMinimalDenom)
	s.Require().Equal(amount.Sub(sdk.NewInt(feeAmount.AmountOf(didtypes.BaseMinimalDenom).Int64())), balance.Amount, "Tax was not subtracted from the fee payer")

	// get supply after tx
	supplyAfterDeflation, _, err := s.app.BankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
	s.Require().NoError(err)

	// check that supply was not deflated
	s.Require().Equal(supplyBeforeDeflation, supplyAfterDeflation, "Supply was deflated")

	// check that reward has been sent to the fee collector
	feeCollector := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	feeCollectorBalance := s.app.BankKeeper.GetBalance(s.ctx, feeCollector, didtypes.BaseMinimalDenom)

	s.Require().Equal(feeCollectorBalance.Amount, feeAmount.AmountOf(didtypes.BaseMinimalDenom), "Fee was not sent to the fee collector")
}
