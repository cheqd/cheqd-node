//go:build upgrade_integration

package integration

import (
	cli "github.com/cheqd/cheqd-node/tests/upgrade/integration/v3/cli"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Burn coins from relevant message signer", func() {
	It("should wait for node catching up", func() {
		By("pinging the node status until catching up is flagged as false")
		err := cli.WaitForCaughtUp(cli.Validator0, cli.CliBinaryName, cli.VotingPeriod*6)
		Expect(err).To(BeNil())
	})

	It("should burn the coins from the given address (here container/validator)", func() {
		coins := sdk.NewCoins(sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdk.NewInt(1000)})
		res, err := cli.BurnMsg(cli.Validator0, coins.String())
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})
})
