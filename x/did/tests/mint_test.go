package tests

import (
	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/did/types"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MsgBurn tests", func() {
	var setup testsetup.TestSetup

	BeforeEach(func() {
		setup = testsetup.Setup()
	})

	It("Valid message format", func() {
		pk1 := ed25519.GenPrivKey().PubKey()
		addr1 := sdk.AccAddress(pk1.Address())
		mintAmount := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(100000)))
		governanceAddress := setup.AccountKeeper.GetModuleAccount(setup.SdkCtx, govtypes.ModuleName).GetAddress().String()

		baseMsg := types.NewMsgMint(
			governanceAddress,
			addr1.String(),
			mintAmount,
		)
		_, err := setup.MsgServer.Mint(setup.SdkCtx, baseMsg)
		Expect(err).To(BeNil())
	})

	It("Not the expected authority address", func() {
		pk1 := ed25519.GenPrivKey().PubKey()
		addr1 := sdk.AccAddress(pk1.Address())
		pk2 := ed25519.GenPrivKey().PubKey()
		add2 := sdk.AccAddress(pk2.Address())
		mintAmount := sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(100000)))

		baseMsg := types.NewMsgMint(
			add2.String(),
			addr1.String(),
			mintAmount,
		)
		_, err := setup.MsgServer.Mint(setup.SdkCtx, baseMsg)
		Expect(err).NotTo(BeNil())
	})
})
