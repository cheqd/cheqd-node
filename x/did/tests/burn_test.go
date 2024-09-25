package tests

import (
	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MsgBurn tests", func() {
	var setup testsetup.TestSetup

	BeforeEach(func() {
		setup = testsetup.Setup()
	})
	It("proper msg", func() {
		pk1 := ed25519.GenPrivKey().PubKey()
		addr1 := sdk.AccAddress(pk1.Address())
		someCoins := sdk.Coins{
			sdk.NewInt64Coin("ncheq", 1000000*1e9), // 1mn CHEQ
		}

		// mint coins to the moduleAccount
		err := setup.BankKeeper.MintCoins(setup.SdkCtx, minttypes.ModuleName, someCoins)
		Expect(err).To(BeNil())

		// FundAccount to the account address
		err = testutil.FundAccount(setup.BankKeeper, setup.SdkCtx, addr1, someCoins)
		Expect(err).To(BeNil())

		balanceBefore := setup.BankKeeper.GetAllBalances(setup.SdkCtx, addr1)
		// make a proper burn message
		burnAmount := sdk.NewCoins(sdk.NewCoin("ncheq", sdk.NewInt(100000)))
		baseMsg := types.NewMsgBurn(
			addr1.String(),
			burnAmount,
		)

		// burn the coins
		_, err = setup.MsgServer.Burn(setup.SdkCtx, baseMsg)
		Expect(err).To(BeNil())

		balanceAfter := setup.BankKeeper.GetAllBalances(setup.SdkCtx, addr1)
		differnce := balanceBefore.Sub(balanceAfter...)
		Expect(burnAmount).To(Equal(differnce))
	})
	It("empty sender", func() {
		pk1 := ed25519.GenPrivKey().PubKey()
		addr1 := sdk.AccAddress(pk1.Address())
		someCoins := sdk.Coins{
			sdk.NewInt64Coin("ncheq", 1000000*1e9), // 1mn CHEQ
		}

		// mint coins to the moduleAccount
		err := setup.BankKeeper.MintCoins(setup.SdkCtx, minttypes.ModuleName, someCoins)
		Expect(err).To(BeNil())

		// FundAccount to the account address
		err = testutil.FundAccount(setup.BankKeeper, setup.SdkCtx, addr1, someCoins)
		Expect(err).To(BeNil())

		// make a proper burn message
		baseMsg := types.NewMsgBurn(
			"",
			sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(100000))),
		)

		// burn the coins
		_, err = setup.MsgServer.Burn(setup.SdkCtx, baseMsg)
		Expect(err).NotTo(BeNil())
	})
	It("zero amount", func() {
		pk1 := ed25519.GenPrivKey().PubKey()
		addr1 := sdk.AccAddress(pk1.Address())
		someCoins := sdk.Coins{
			sdk.NewInt64Coin("ncheq", 1000000*1e9), // 1mn CHEQ
		}

		// mint coins to the moduleAccount
		err := setup.BankKeeper.MintCoins(setup.SdkCtx, minttypes.ModuleName, someCoins)
		Expect(err).To(BeNil())

		// FundAccount to the account address
		err = testutil.FundAccount(setup.BankKeeper, setup.SdkCtx, addr1, someCoins)
		Expect(err).To(BeNil())

		// make a proper burn message
		baseMsg := types.NewMsgBurn(
			addr1.String(),
			sdk.NewCoins(sdk.NewCoin("ncheq", sdk.ZeroInt())),
		)

		// burn the coins
		_, err = setup.MsgServer.Burn(setup.SdkCtx, baseMsg)
		Expect(err).NotTo(BeNil())
	})
	It("invalid denom", func() {
		pk1 := ed25519.GenPrivKey().PubKey()
		addr1 := sdk.AccAddress(pk1.Address())
		someCoins := sdk.Coins{
			sdk.NewInt64Coin("ncheq", 1000000*1e9), // 1mn CHEQ
		}

		// mint coins to the moduleAccount
		err := setup.BankKeeper.MintCoins(setup.SdkCtx, minttypes.ModuleName, someCoins)
		Expect(err).To(BeNil())

		// FundAccount to the account address
		err = testutil.FundAccount(setup.BankKeeper, setup.SdkCtx, addr1, someCoins)
		Expect(err).To(BeNil())

		// make a proper burn message
		burnAmount := sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(100000)))
		baseMsg := types.NewMsgBurn(
			addr1.String(),
			burnAmount,
		)

		// burn the coins
		_, err = setup.MsgServer.Burn(setup.SdkCtx, baseMsg)
		Expect(err).NotTo(BeNil())
	})
})
