package ante_test

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	cheqdante "github.com/cheqd/cheqd-node/ante"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
)

func (s *AnteTestSuite) TestDeductFeeDecorator_ZeroGas() {
	s.SetupTest(true) // setup
	s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

	mfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	coins := sdk.NewCoins(sdk.NewCoin("cheq", sdk.NewInt(300)))
	testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, coins)

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

	mfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	coins := sdk.NewCoins(sdk.NewCoin("cheq", sdk.NewInt(300)))
	testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, coins)

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
	cheqPrice := sdk.NewDecCoinFromDec("cheq", sdk.NewDec(20))
	highGasPrice := []sdk.DecCoin{cheqPrice}
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

	cheqPrice = sdk.NewDecCoinFromDec("cheq", sdk.NewDec(0).Quo(sdk.NewDec(100000)))
	lowGasPrice := []sdk.DecCoin{cheqPrice}
	s.ctx = s.ctx.WithMinGasPrices(lowGasPrice)

	newCtx, err := antehandler(s.ctx, tx, false)
	s.Require().Nil(err, "Decorator should not have errored on fee higher than local gasPrice")
	// Priority is the smallest gas price amount in any denom. Since we have only 1 gas price
	// of 10cheq, the priority here is 10.
	s.Require().Equal(int64(10), newCtx.Priority())
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
	coins := sdk.NewCoins(sdk.NewCoin("cheq", sdk.NewInt(10)))
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, coins)
	s.Require().NoError(err)

	dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, nil, nil)
	antehandler := sdk.ChainAnteDecorators(dfd)

	_, err = antehandler(s.ctx, tx, false)

	s.Require().NotNil(err, "Tx did not error when fee payer had insufficient funds")

	// Set account with sufficient funds
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin("cheq", sdk.NewInt(200))))
	s.Require().NoError(err)

	_, err = antehandler(s.ctx, tx, false)

	s.Require().Nil(err, "Tx errored after account has been set with sufficient funds")
}

func (s *AnteTestSuite) TestCheckDeductFeeWithCustomFixedFee_SingleMsgTaxableTx() {
	s.SetupTest(true) // setup
	s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	msg := NewTestDidMsg()
	feeAmount := NewTestFeeAmountMinimalDenomLTFixedFee()
	gasLimit := testdata.NewTestGasLimit()
	s.Require().NoError(s.txBuilder.SetMsgs(msg))
	s.txBuilder.SetFeeAmount(feeAmount)
	s.txBuilder.SetGasLimit(gasLimit)
	s.txBuilder.SetFeePayer(addr1)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
	s.Require().NoError(err)

	// Set account with insufficient funds
	acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	coins := sdk.NewCoins(sdk.NewCoin("ncheq", sdk.NewInt(200 * 1e9)))
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, coins)
	s.Require().NoError(err)

	dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(dfd)

	_, err = antehandler(s.ctx, tx, false)

	s.Require().NotNil(err, "Tx did not error when fee payer had insufficient funds")

	// Set account with sufficient funds
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin("ncheq", sdk.NewInt(1e13))))
	s.Require().NoError(err)

	// Set more than sufficient fee provided 10,000 CHEQ > 300 CHEQ
	s.txBuilder.SetFeeAmount(feeAmount.MulInt(sdk.NewInt(10)))
	tx = s.txBuilder.GetTx()

	// Set new event manager to capture fee deduction workflow events only
	em := sdk.NewEventManager()
	s.ctx = s.ctx.WithEventManager(em)

	_, err = antehandler(s.ctx, tx, false)

	s.Require().Nil(err, "Tx errored after account has been set with sufficient funds")

	// 1. Check type & order of events emitted
	//     a. `Transfer` event - Fee deduction from fee payer to `cheqd` module account
	//     b. `BurnFee` event - Fee burn from `cheqd` module account
	//     c. `Transfer` event - Fee distribution from `cheqd` module account to fee collector
	// 2. Check event attributes & sender/recipient addresses

	// Prepare relevant module account addresses
	cheqdModuleAccAddr := s.app.AccountKeeper.GetModuleAddress(cheqdtypes.ModuleName)
	feeCollectorModuleAccAddr := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)

	// 1. Check type & order of events emitted
	events := em.Events()
	s.Require().Len(events, 11)
	s.Require().Equal(events[0].Type, "coin_spent") 		// 1.1 Fee deduction from fee payer	
	s.Require().Equal(events[1].Type, "coin_received")		// 1.2 Fee received by `cheqd` module account
	s.Require().Equal(events[2].Type, "transfer") 			// 1.3 Fee transfer encapsulating fee deduction and fee received
	s.Require().Equal(events[3].Type, "message") 			// 1.4 Tx message specifying the sender
	s.Require().Equal(events[4].Type, "tx") 				// 1.5 Tx specifying the fee amount, fee payer
	s.Require().Equal(events[5].Type, "coin_spent") 		// 1.6 Fee burn initiated as coin spent from `cheqd` module account (burn portion)
	s.Require().Equal(events[6].Type, "burn")				// 1.7 Fee burn specifying burner as `cheqd` module account (burn portion)
	s.Require().Equal(events[7].Type, "coin_spent")			// 1.8 Fee deduction from `cheqd` module account (rewards portion)
	s.Require().Equal(events[8].Type, "coin_received")		// 1.9 Fee received by `feeCollector` module account (rewards portion)
	s.Require().Equal(events[9].Type, "transfer")			// 1.10 Fee transfer encapsulating fee deduction and fee received (rewards portion)
	s.Require().Equal(events[10].Type, "message")			// 1.11 Tx message specifying the sender

	// Prepare relevant event attributes
	senderKey							:=		[]byte(sdk.AttributeKeySender)
	amountKey							:=		[]byte(sdk.AttributeKeyAmount)
	feeKey								:=		[]byte(sdk.AttributeKeyFee)
	feePayerKey							:=		[]byte(sdk.AttributeKeyFeePayer)
	spenderKey							:=		[]byte("spender")
	receiverKey							:=		[]byte("receiver")
	recipientKey						:=	 	[]byte("recipient")
	burnerKey							:=		[]byte("burner")

	addr1Value							:=		[]byte(addr1.String())
	cheqdModuleAccAddrValue				:=		[]byte(cheqdModuleAccAddr.String())
	feeCollectorModuleAccAddrValue 		:=		[]byte(feeCollectorModuleAccAddr.String())

	bigFeeValue							:=		[]byte("1500000000000ncheq")
	splitFeeValue						:=		[]byte("750000000000ncheq")

	// 2. Check event attributes & sender/recipient addresses
	s.Require().Equal(events[0].Attributes[0].Key, spenderKey) 									// 2.1.a Fee deduction from fee payer (sender attr)
	s.Require().Equal(events[0].Attributes[0].Value, addr1Value)								// 2.1.b Value matches fee payer address
	s.Require().Equal(events[0].Attributes[1].Key, amountKey)									// 2.1.c Fee deduction from fee payer (amount attr)
	s.Require().Equal(events[0].Attributes[1].Value, bigFeeValue)								// 2.1.d Value matches fee amount
	s.Require().Equal(events[1].Attributes[0].Key, receiverKey)									// 2.2.a Fee received by `cheqd` module account (receiver attr)
	s.Require().Equal(events[1].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.2.b Value matches `cheqd` module account address
	s.Require().Equal(events[1].Attributes[1].Key, amountKey)									// 2.2.c Fee received by `cheqd` module account (amount attr)
	s.Require().Equal(events[1].Attributes[1].Value, bigFeeValue)								// 2.2.d Value matches fee amount
	s.Require().Equal(events[2].Attributes[0].Key, recipientKey)								// 2.3.a Fee transfer encapsulating fee deduction and fee received (recipient attr)
	s.Require().Equal(events[2].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.3.b Value matches `cheqd` module account address
	s.Require().Equal(events[2].Attributes[1].Key, senderKey)									// 2.3.c Fee transfer encapsulating fee deduction and fee received (sender attr)
	s.Require().Equal(events[2].Attributes[1].Value, addr1Value)								// 2.3.d Value matches fee payer address
	s.Require().Equal(events[2].Attributes[2].Key, amountKey)									// 2.3.e Fee transfer encapsulating fee deduction and fee received (amount attr)
	s.Require().Equal(events[2].Attributes[2].Value, bigFeeValue)								// 2.3.f Value matches fee amount
	s.Require().Equal(events[3].Attributes[0].Key, senderKey)									// 2.4.a Tx message specifying the sender (sender attr)
	s.Require().Equal(events[3].Attributes[0].Value, addr1Value)								// 2.4.b Value matches fee payer address
	s.Require().Equal(events[4].Attributes[0].Key, feeKey)										// 2.5.a Tx specifying the fee amount, fee payer (fee attr)
	s.Require().Equal(events[4].Attributes[0].Value, bigFeeValue)								// 2.5.b Value matches fee amount
	s.Require().Equal(events[4].Attributes[1].Key, feePayerKey)									// 2.5.c Tx specifying the fee amount, fee payer (fee_payer attr)
	s.Require().Equal(events[4].Attributes[1].Value, addr1Value)								// 2.5.d Value matches fee payer address
	s.Require().Equal(events[5].Attributes[0].Key, spenderKey)									// 2.6.a Fee burn initiated as coin spent from `cheqd` module account (burn portion) (spender attr)
	s.Require().Equal(events[5].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.6.b Value matches `cheqd` module account address
	s.Require().Equal(events[5].Attributes[1].Key, amountKey)									// 2.6.c Fee burn initiated as coin spent from `cheqd` module account (burn portion) (amount attr)
	s.Require().Equal(events[5].Attributes[1].Value, splitFeeValue)								// 2.6.d Value matches fee amount
	s.Require().Equal(events[6].Attributes[0].Key, burnerKey)									// 2.7.a Fee burn specifying burner as `cheqd` module account (burn portion) (burner attr)
	s.Require().Equal(events[6].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.7.b Value matches `cheqd` module account address
	s.Require().Equal(events[6].Attributes[1].Key, amountKey)									// 2.7.c Fee burn specifying burner as `cheqd` module account (burn portion) (amount attr)
	s.Require().Equal(events[6].Attributes[1].Value, splitFeeValue)								// 2.7.d Value matches fee amount
	s.Require().Equal(events[7].Attributes[0].Key, spenderKey)									// 2.8.a Fee deduction from `cheqd` module account (rewards portion) (spender attr)
	s.Require().Equal(events[7].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.8.b Value matches `cheqd` module account address
	s.Require().Equal(events[7].Attributes[1].Key, amountKey)									// 2.8.c Fee deduction from `cheqd` module account (rewards portion) (amount attr)
	s.Require().Equal(events[7].Attributes[1].Value, splitFeeValue)								// 2.8.d Value matches fee amount
	s.Require().Equal(events[8].Attributes[0].Key, receiverKey)									// 2.9.a Fee received by fee collector (rewards portion) (receiver attr)
	s.Require().Equal(events[8].Attributes[0].Value, feeCollectorModuleAccAddrValue)			// 2.9.b Value matches fee collector address
	s.Require().Equal(events[8].Attributes[1].Key, amountKey)									// 2.9.c Fee received by fee collector (rewards portion) (amount attr)
	s.Require().Equal(events[8].Attributes[1].Value, splitFeeValue)								// 2.9.d Value matches fee amount
	s.Require().Equal(events[9].Attributes[0].Key, recipientKey)								// 2.10.a Fee transfer encapsulating fee deduction and fee received (rewards portion) (recipient attr)
	s.Require().Equal(events[9].Attributes[0].Value, feeCollectorModuleAccAddrValue)			// 2.10.b Value matches fee collector address
	s.Require().Equal(events[9].Attributes[1].Key, senderKey)									// 2.10.c Fee transfer encapsulating fee deduction and fee received (rewards portion) (sender attr)
	s.Require().Equal(events[9].Attributes[1].Value, cheqdModuleAccAddrValue)					// 2.10.d Value matches `cheqd` module account address
	s.Require().Equal(events[9].Attributes[2].Key, amountKey)									// 2.10.e Fee transfer encapsulating fee deduction and fee received (rewards portion) (amount attr)
	s.Require().Equal(events[9].Attributes[2].Value, splitFeeValue)								// 2.10.f Value matches fee amount
	s.Require().Equal(events[10].Attributes[0].Key, senderKey)									// 2.11.a Tx message specifying the sender (sender attr)
	s.Require().Equal(events[10].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.11.b Value matches `cheqd` module account address
}

func (s *AnteTestSuite) TestCheckDeductFeeWithCustomFixedFee_MultipleMsgTaxableTx() {
	s.SetupTest(true) // setup
	s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	cheqdMsg := NewTestDidMsg()
	resourceMsg := NewTestResourceMsg()
	feeAmount := NewTestFeeAmountMinimalDenomLTFixedFee()
	gasLimit := testdata.NewTestGasLimit()
	s.Require().NoError(s.txBuilder.SetMsgs([]sdk.Msg{cheqdMsg, resourceMsg, cheqdMsg, resourceMsg, cheqdMsg}...)) // 3x cheqdMsg, 2x resourceMsg = 5 taxable msgs
	s.txBuilder.SetFeeAmount(feeAmount)
	s.txBuilder.SetGasLimit(gasLimit)
	s.txBuilder.SetFeePayer(addr1)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
	s.Require().NoError(err)

	// Set account with insufficient funds
	acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	coins := NewTestFeeAmountMinimalDenomEFixedFee() // 300 CHEQ aka single Msg fee
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, coins)
	s.Require().NoError(err)

	dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(dfd)

	_, err = antehandler(s.ctx, tx, false)

	s.Require().NotNil(err, "Tx did not error when fee payer had insufficient funds")

	// Set account with sufficient funds
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin("ncheq", sdk.NewInt(1e13))))
	s.Require().NoError(err)

	// Set exactly expected fee 1500 CHEQ
	feeAmount = NewTestFeeAmountMinimalDenomGTFixedFee()
	s.txBuilder.SetFeeAmount(feeAmount)
	tx = s.txBuilder.GetTx()

	// Set new event manager to capture fee deduction workflow events only
	em := sdk.NewEventManager()
	s.ctx = s.ctx.WithEventManager(em)

	_, err = antehandler(s.ctx, tx, false)

	s.Require().Nil(err, "Tx errored after account has been set with sufficient funds")

	// 1. Check type & order of events emitted
	//     a. `Transfer` event - Fee deduction from fee payer to `cheqd` module account
	//     b. `BurnFee` event - Fee burn from `cheqd` module account
	//     c. `Transfer` event - Fee distribution from `cheqd` module account to fee collector
	// 2. Check event attributes & sender/recipient addresses

	// Prepare relevant module account addresses
	cheqdModuleAccAddr := s.app.AccountKeeper.GetModuleAddress(cheqdtypes.ModuleName)
	feeCollectorModuleAccAddr := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)

	// 1. Check type & order of events emitted
	events := em.Events()
	s.Require().Len(events, 11)
	s.Require().Equal(events[0].Type, "coin_spent") 		// 1.1 Fee deduction from fee payer	
	s.Require().Equal(events[1].Type, "coin_received")		// 1.2 Fee received by `cheqd` module account
	s.Require().Equal(events[2].Type, "transfer") 			// 1.3 Fee transfer encapsulating fee deduction and fee received
	s.Require().Equal(events[3].Type, "message") 			// 1.4 Tx message specifying the sender
	s.Require().Equal(events[4].Type, "tx") 				// 1.5 Tx specifying the fee amount, fee payer
	s.Require().Equal(events[5].Type, "coin_spent") 		// 1.6 Fee burn initiated as coin spent from `cheqd` module account (burn portion)
	s.Require().Equal(events[6].Type, "burn")				// 1.7 Fee burn specifying burner as `cheqd` module account (burn portion)
	s.Require().Equal(events[7].Type, "coin_spent")			// 1.8 Fee deduction from `cheqd` module account (rewards portion)
	s.Require().Equal(events[8].Type, "coin_received")		// 1.9 Fee received by `feeCollector` module account (rewards portion)
	s.Require().Equal(events[9].Type, "transfer")			// 1.10 Fee transfer encapsulating fee deduction and fee received (rewards portion)
	s.Require().Equal(events[10].Type, "message")			// 1.11 Tx message specifying the sender

	// Prepare relevant event attributes
	senderKey							:=		[]byte(sdk.AttributeKeySender)
	amountKey							:=		[]byte(sdk.AttributeKeyAmount)
	feeKey								:=		[]byte(sdk.AttributeKeyFee)
	feePayerKey							:=		[]byte(sdk.AttributeKeyFeePayer)
	spenderKey							:=		[]byte("spender")
	receiverKey							:=		[]byte("receiver")
	recipientKey						:=	 	[]byte("recipient")
	burnerKey							:=		[]byte("burner")

	addr1Value							:=		[]byte(addr1.String())
	cheqdModuleAccAddrValue				:=		[]byte(cheqdModuleAccAddr.String())
	feeCollectorModuleAccAddrValue 		:=		[]byte(feeCollectorModuleAccAddr.String())

	bigFeeValue							:=		[]byte("1500000000000ncheq")
	splitFeeValue						:=		[]byte("750000000000ncheq")

	// 2. Check event attributes & sender/recipient addresses
	s.Require().Equal(events[0].Attributes[0].Key, spenderKey) 									// 2.1.a Fee deduction from fee payer (sender attr)
	s.Require().Equal(events[0].Attributes[0].Value, addr1Value)								// 2.1.b Value matches fee payer address
	s.Require().Equal(events[0].Attributes[1].Key, amountKey)									// 2.1.c Fee deduction from fee payer (amount attr)
	s.Require().Equal(events[0].Attributes[1].Value, bigFeeValue)								// 2.1.d Value matches fee amount
	s.Require().Equal(events[1].Attributes[0].Key, receiverKey)									// 2.2.a Fee received by `cheqd` module account (receiver attr)
	s.Require().Equal(events[1].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.2.b Value matches `cheqd` module account address
	s.Require().Equal(events[1].Attributes[1].Key, amountKey)									// 2.2.c Fee received by `cheqd` module account (amount attr)
	s.Require().Equal(events[1].Attributes[1].Value, bigFeeValue)								// 2.2.d Value matches fee amount
	s.Require().Equal(events[2].Attributes[0].Key, recipientKey)								// 2.3.a Fee transfer encapsulating fee deduction and fee received (recipient attr)
	s.Require().Equal(events[2].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.3.b Value matches `cheqd` module account address
	s.Require().Equal(events[2].Attributes[1].Key, senderKey)									// 2.3.c Fee transfer encapsulating fee deduction and fee received (sender attr)
	s.Require().Equal(events[2].Attributes[1].Value, addr1Value)								// 2.3.d Value matches fee payer address
	s.Require().Equal(events[2].Attributes[2].Key, amountKey)									// 2.3.e Fee transfer encapsulating fee deduction and fee received (amount attr)
	s.Require().Equal(events[2].Attributes[2].Value, bigFeeValue)								// 2.3.f Value matches fee amount
	s.Require().Equal(events[3].Attributes[0].Key, senderKey)									// 2.4.a Tx message specifying the sender (sender attr)
	s.Require().Equal(events[3].Attributes[0].Value, addr1Value)								// 2.4.b Value matches fee payer address
	s.Require().Equal(events[4].Attributes[0].Key, feeKey)										// 2.5.a Tx specifying the fee amount, fee payer (fee attr)
	s.Require().Equal(events[4].Attributes[0].Value, bigFeeValue)								// 2.5.b Value matches fee amount
	s.Require().Equal(events[4].Attributes[1].Key, feePayerKey)									// 2.5.c Tx specifying the fee amount, fee payer (fee_payer attr)
	s.Require().Equal(events[4].Attributes[1].Value, addr1Value)								// 2.5.d Value matches fee payer address
	s.Require().Equal(events[5].Attributes[0].Key, spenderKey)									// 2.6.a Fee burn initiated as coin spent from `cheqd` module account (burn portion) (spender attr)
	s.Require().Equal(events[5].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.6.b Value matches `cheqd` module account address
	s.Require().Equal(events[5].Attributes[1].Key, amountKey)									// 2.6.c Fee burn initiated as coin spent from `cheqd` module account (burn portion) (amount attr)
	s.Require().Equal(events[5].Attributes[1].Value, splitFeeValue)								// 2.6.d Value matches fee amount
	s.Require().Equal(events[6].Attributes[0].Key, burnerKey)									// 2.7.a Fee burn specifying burner as `cheqd` module account (burn portion) (burner attr)
	s.Require().Equal(events[6].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.7.b Value matches `cheqd` module account address
	s.Require().Equal(events[6].Attributes[1].Key, amountKey)									// 2.7.c Fee burn specifying burner as `cheqd` module account (burn portion) (amount attr)
	s.Require().Equal(events[6].Attributes[1].Value, splitFeeValue)								// 2.7.d Value matches fee amount
	s.Require().Equal(events[7].Attributes[0].Key, spenderKey)									// 2.8.a Fee deduction from `cheqd` module account (rewards portion) (spender attr)
	s.Require().Equal(events[7].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.8.b Value matches `cheqd` module account address
	s.Require().Equal(events[7].Attributes[1].Key, amountKey)									// 2.8.c Fee deduction from `cheqd` module account (rewards portion) (amount attr)
	s.Require().Equal(events[7].Attributes[1].Value, splitFeeValue)								// 2.8.d Value matches fee amount
	s.Require().Equal(events[8].Attributes[0].Key, receiverKey)									// 2.9.a Fee received by fee collector (rewards portion) (receiver attr)
	s.Require().Equal(events[8].Attributes[0].Value, feeCollectorModuleAccAddrValue)			// 2.9.b Value matches fee collector address
	s.Require().Equal(events[8].Attributes[1].Key, amountKey)									// 2.9.c Fee received by fee collector (rewards portion) (amount attr)
	s.Require().Equal(events[8].Attributes[1].Value, splitFeeValue)								// 2.9.d Value matches fee amount
	s.Require().Equal(events[9].Attributes[0].Key, recipientKey)								// 2.10.a Fee transfer encapsulating fee deduction and fee received (rewards portion) (recipient attr)
	s.Require().Equal(events[9].Attributes[0].Value, feeCollectorModuleAccAddrValue)			// 2.10.b Value matches fee collector address
	s.Require().Equal(events[9].Attributes[1].Key, senderKey)									// 2.10.c Fee transfer encapsulating fee deduction and fee received (rewards portion) (sender attr)
	s.Require().Equal(events[9].Attributes[1].Value, cheqdModuleAccAddrValue)					// 2.10.d Value matches `cheqd` module account address
	s.Require().Equal(events[9].Attributes[2].Key, amountKey)									// 2.10.e Fee transfer encapsulating fee deduction and fee received (rewards portion) (amount attr)
	s.Require().Equal(events[9].Attributes[2].Value, splitFeeValue)								// 2.10.f Value matches fee amount
	s.Require().Equal(events[10].Attributes[0].Key, senderKey)									// 2.11.a Tx message specifying the sender (sender attr)
	s.Require().Equal(events[10].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.11.b Value matches `cheqd` module account address
}

func (s *AnteTestSuite) TestCheckDeductFeeWithCustomFixedFee_MixedMsgTaxableNonTaxableTx() {
	s.SetupTest(true) // setup
	s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	cheqdMsg := NewTestDidMsg()
	resourceMsg := NewTestResourceMsg()
	testMsg := testdata.NewTestMsg()
	feeAmount := NewTestFeeAmountMinimalDenomLTFixedFee()
	gasLimit := testdata.NewTestGasLimit()
	s.Require().NoError(s.txBuilder.SetMsgs([]sdk.Msg{cheqdMsg, resourceMsg, cheqdMsg, resourceMsg, cheqdMsg, testMsg, testMsg, testMsg }...)) // 3x cheqdMsg + 2x resourceMsg + 3x testMsg = 8 msgs, 5 taxable, 3 non-taxable
	s.txBuilder.SetFeeAmount(feeAmount)
	s.txBuilder.SetGasLimit(gasLimit)
	s.txBuilder.SetFeePayer(addr1)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
	s.Require().NoError(err)

	// Set account with insufficient funds
	acc := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, addr1)
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	coins := NewTestFeeAmountMinimalDenomEFixedFee() // 300 CHEQ aka single Msg fee
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, coins)
	s.Require().NoError(err)

	dfd := cheqdante.NewDeductFeeDecorator(s.app.AccountKeeper, s.app.BankKeeper, s.app.FeeGrantKeeper, nil)
	antehandler := sdk.ChainAnteDecorators(dfd)

	_, err = antehandler(s.ctx, tx, false)

	s.Require().NotNil(err, "Tx did not error when fee payer had insufficient funds")

	// Set account with sufficient funds
	s.app.AccountKeeper.SetAccount(s.ctx, acc)
	err = testutil.FundAccount(s.app.BankKeeper, s.ctx, addr1, sdk.NewCoins(sdk.NewCoin("ncheq", sdk.NewInt(1e13))))
	s.Require().NoError(err)

	// Set exactly expected fee 1500 CHEQ, for 8 msgs, 5 taxable, 3 non-taxable
	// Non-taxable msgs are 3x testMsg, each with FF - BM*FF = 300 - 0.5*300 = 150 CHEQ to cover the minGasPrice set by the validator
	feeAmount = NewTestFeeAmountMinimalDenomGTFixedFee()
	s.txBuilder.SetFeeAmount(feeAmount)
	tx = s.txBuilder.GetTx()

	// Set new event manager to capture fee deduction workflow events only
	em := sdk.NewEventManager()
	s.ctx = s.ctx.WithEventManager(em)

	_, err = antehandler(s.ctx, tx, false)

	s.Require().Nil(err, "Tx errored after account has been set with sufficient funds")

	// 1. Check type & order of events emitted
	//     a. `Transfer` event - Fee deduction from fee payer to `cheqd` module account
	//     b. `BurnFee` event - Fee burn from `cheqd` module account
	//     c. `Transfer` event - Fee distribution from `cheqd` module account to fee collector
	// 2. Check event attributes & sender/recipient addresses

	// Prepare relevant module account addresses
	cheqdModuleAccAddr := s.app.AccountKeeper.GetModuleAddress(cheqdtypes.ModuleName)
	feeCollectorModuleAccAddr := s.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)

	// 1. Check type & order of events emitted
	events := em.Events()
	s.Require().Len(events, 11)
	s.Require().Equal(events[0].Type, "coin_spent") 		// 1.1 Fee deduction from fee payer	
	s.Require().Equal(events[1].Type, "coin_received")		// 1.2 Fee received by `cheqd` module account
	s.Require().Equal(events[2].Type, "transfer") 			// 1.3 Fee transfer encapsulating fee deduction and fee received
	s.Require().Equal(events[3].Type, "message") 			// 1.4 Tx message specifying the sender
	s.Require().Equal(events[4].Type, "tx") 				// 1.5 Tx specifying the fee amount, fee payer
	s.Require().Equal(events[5].Type, "coin_spent") 		// 1.6 Fee burn initiated as coin spent from `cheqd` module account (burn portion)
	s.Require().Equal(events[6].Type, "burn")				// 1.7 Fee burn specifying burner as `cheqd` module account (burn portion)
	s.Require().Equal(events[7].Type, "coin_spent")			// 1.8 Fee deduction from `cheqd` module account (rewards portion)
	s.Require().Equal(events[8].Type, "coin_received")		// 1.9 Fee received by `feeCollector` module account (rewards portion)
	s.Require().Equal(events[9].Type, "transfer")			// 1.10 Fee transfer encapsulating fee deduction and fee received (rewards portion)
	s.Require().Equal(events[10].Type, "message")			// 1.11 Tx message specifying the sender

	// Prepare relevant event attributes
	senderKey							:=		[]byte(sdk.AttributeKeySender)
	amountKey							:=		[]byte(sdk.AttributeKeyAmount)
	feeKey								:=		[]byte(sdk.AttributeKeyFee)
	feePayerKey							:=		[]byte(sdk.AttributeKeyFeePayer)
	spenderKey							:=		[]byte("spender")
	receiverKey							:=		[]byte("receiver")
	recipientKey						:=	 	[]byte("recipient")
	burnerKey							:=		[]byte("burner")

	addr1Value							:=		[]byte(addr1.String())
	cheqdModuleAccAddrValue				:=		[]byte(cheqdModuleAccAddr.String())
	feeCollectorModuleAccAddrValue 		:=		[]byte(feeCollectorModuleAccAddr.String())

	bigFeeValue							:=		[]byte("1500000000000ncheq")
	splitFeeValue						:=		[]byte("750000000000ncheq")

	// 2. Check event attributes & sender/recipient addresses
	s.Require().Equal(events[0].Attributes[0].Key, spenderKey) 									// 2.1.a Fee deduction from fee payer (sender attr)
	s.Require().Equal(events[0].Attributes[0].Value, addr1Value)								// 2.1.b Value matches fee payer address
	s.Require().Equal(events[0].Attributes[1].Key, amountKey)									// 2.1.c Fee deduction from fee payer (amount attr)
	s.Require().Equal(events[0].Attributes[1].Value, bigFeeValue)								// 2.1.d Value matches fee amount
	s.Require().Equal(events[1].Attributes[0].Key, receiverKey)									// 2.2.a Fee received by `cheqd` module account (receiver attr)
	s.Require().Equal(events[1].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.2.b Value matches `cheqd` module account address
	s.Require().Equal(events[1].Attributes[1].Key, amountKey)									// 2.2.c Fee received by `cheqd` module account (amount attr)
	s.Require().Equal(events[1].Attributes[1].Value, bigFeeValue)								// 2.2.d Value matches fee amount
	s.Require().Equal(events[2].Attributes[0].Key, recipientKey)								// 2.3.a Fee transfer encapsulating fee deduction and fee received (recipient attr)
	s.Require().Equal(events[2].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.3.b Value matches `cheqd` module account address
	s.Require().Equal(events[2].Attributes[1].Key, senderKey)									// 2.3.c Fee transfer encapsulating fee deduction and fee received (sender attr)
	s.Require().Equal(events[2].Attributes[1].Value, addr1Value)								// 2.3.d Value matches fee payer address
	s.Require().Equal(events[2].Attributes[2].Key, amountKey)									// 2.3.e Fee transfer encapsulating fee deduction and fee received (amount attr)
	s.Require().Equal(events[2].Attributes[2].Value, bigFeeValue)								// 2.3.f Value matches fee amount
	s.Require().Equal(events[3].Attributes[0].Key, senderKey)									// 2.4.a Tx message specifying the sender (sender attr)
	s.Require().Equal(events[3].Attributes[0].Value, addr1Value)								// 2.4.b Value matches fee payer address
	s.Require().Equal(events[4].Attributes[0].Key, feeKey)										// 2.5.a Tx specifying the fee amount, fee payer (fee attr)
	s.Require().Equal(events[4].Attributes[0].Value, bigFeeValue)								// 2.5.b Value matches fee amount
	s.Require().Equal(events[4].Attributes[1].Key, feePayerKey)									// 2.5.c Tx specifying the fee amount, fee payer (fee_payer attr)
	s.Require().Equal(events[4].Attributes[1].Value, addr1Value)								// 2.5.d Value matches fee payer address
	s.Require().Equal(events[5].Attributes[0].Key, spenderKey)									// 2.6.a Fee burn initiated as coin spent from `cheqd` module account (burn portion) (spender attr)
	s.Require().Equal(events[5].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.6.b Value matches `cheqd` module account address
	s.Require().Equal(events[5].Attributes[1].Key, amountKey)									// 2.6.c Fee burn initiated as coin spent from `cheqd` module account (burn portion) (amount attr)
	s.Require().Equal(events[5].Attributes[1].Value, splitFeeValue)								// 2.6.d Value matches fee amount
	s.Require().Equal(events[6].Attributes[0].Key, burnerKey)									// 2.7.a Fee burn specifying burner as `cheqd` module account (burn portion) (burner attr)
	s.Require().Equal(events[6].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.7.b Value matches `cheqd` module account address
	s.Require().Equal(events[6].Attributes[1].Key, amountKey)									// 2.7.c Fee burn specifying burner as `cheqd` module account (burn portion) (amount attr)
	s.Require().Equal(events[6].Attributes[1].Value, splitFeeValue)								// 2.7.d Value matches fee amount
	s.Require().Equal(events[7].Attributes[0].Key, spenderKey)									// 2.8.a Fee deduction from `cheqd` module account (rewards portion) (spender attr)
	s.Require().Equal(events[7].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.8.b Value matches `cheqd` module account address
	s.Require().Equal(events[7].Attributes[1].Key, amountKey)									// 2.8.c Fee deduction from `cheqd` module account (rewards portion) (amount attr)
	s.Require().Equal(events[7].Attributes[1].Value, splitFeeValue)								// 2.8.d Value matches fee amount
	s.Require().Equal(events[8].Attributes[0].Key, receiverKey)									// 2.9.a Fee received by fee collector (rewards portion) (receiver attr)
	s.Require().Equal(events[8].Attributes[0].Value, feeCollectorModuleAccAddrValue)			// 2.9.b Value matches fee collector address
	s.Require().Equal(events[8].Attributes[1].Key, amountKey)									// 2.9.c Fee received by fee collector (rewards portion) (amount attr)
	s.Require().Equal(events[8].Attributes[1].Value, splitFeeValue)								// 2.9.d Value matches fee amount
	s.Require().Equal(events[9].Attributes[0].Key, recipientKey)								// 2.10.a Fee transfer encapsulating fee deduction and fee received (rewards portion) (recipient attr)
	s.Require().Equal(events[9].Attributes[0].Value, feeCollectorModuleAccAddrValue)			// 2.10.b Value matches fee collector address
	s.Require().Equal(events[9].Attributes[1].Key, senderKey)									// 2.10.c Fee transfer encapsulating fee deduction and fee received (rewards portion) (sender attr)
	s.Require().Equal(events[9].Attributes[1].Value, cheqdModuleAccAddrValue)					// 2.10.d Value matches `cheqd` module account address
	s.Require().Equal(events[9].Attributes[2].Key, amountKey)									// 2.10.e Fee transfer encapsulating fee deduction and fee received (rewards portion) (amount attr)
	s.Require().Equal(events[9].Attributes[2].Value, splitFeeValue)								// 2.10.f Value matches fee amount
	s.Require().Equal(events[10].Attributes[0].Key, senderKey)									// 2.11.a Tx message specifying the sender (sender attr)
	s.Require().Equal(events[10].Attributes[0].Value, cheqdModuleAccAddrValue)					// 2.11.b Value matches `cheqd` module account address
}