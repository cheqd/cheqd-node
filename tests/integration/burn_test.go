//go:build upgrade_integration

package integration

import (
	"fmt"

	cli "github.com/cheqd/cheqd-node/tests/integration/cli"
	helpers "github.com/cheqd/cheqd-node/tests/integration/helpers"

	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Burn coins from relevant message signer", func() {
	It("should burn the coins from the given address", func() {
		coins := sdk.NewCoin("ncheq", sdk.NewInt(10000000))

		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, "ncheq")
		Expect(err).To(BeNil())
		fmt.Println("balance before>>>>>>>>>.", balanceBefore)
		fees := helpers.GenerateFees("50000000000ncheq")
		res, err := cli.BurnMsg(testdata.BASE_ACCOUNT_1, coins.String(), fees)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, "ncheq")
		Expect(err).To(BeNil())
		fmt.Println("balance after>>>>>>>>>>>", balanceAfter)
		diff := balanceBefore.Sub(balanceAfter)

		// tnx fee + coins burned
		feededucted := coins.Add(sdk.NewCoin("ncheq", sdk.NewInt(50000000000)))
		Expect(diff).To(Equal(feededucted))
	})

	It("shouldn't burn as their are insufficient funds in the sender", func() {
		coins := sdk.NewCoin("ncheq", sdk.NewInt(10000000))

		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_3_ADDR, "ncheq")
		Expect(err).To(BeNil())
		fmt.Println("balance before>>>>>>>>>.", balanceBefore)
		fees := helpers.GenerateFees("50000000000ncheq")
		res, err := cli.BurnMsg(testdata.BASE_ACCOUNT_3, coins.String(), fees)
		Expect(err).NotTo(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_3_ADDR, "ncheq")
		Expect(err).To(BeNil())
		fmt.Println("balance after>>>>>>>>>>>", balanceAfter)
	})
})
