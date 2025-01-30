//go:build integration

package integration

import (
	cli "github.com/cheqd/cheqd-node/tests/integration/cli"
	helpers "github.com/cheqd/cheqd-node/tests/integration/helpers"

	sdkmath "cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Burn coins from relevant message signer", func() {
	It("should burn the coins from the given address", func() {
		// define the coins to burn, in which case 1,000,000 ncheq or 0.01 cheq
		burnCoins := sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(1_000_000))

		// get the balance of the account before burning the coins
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, didtypes.BaseMinimalDenom)

		// assert no error
		Expect(err).To(BeNil())

		// generate fixed fees, in which case 3,500,000,000 ncheq or 3.5 cheq
		fees := helpers.GenerateFees("3500000000ncheq")

		// burn the coins
		res, err := cli.BurnMsg(testdata.BASE_ACCOUNT_1, burnCoins.String(), fees)

		// assert no error
		Expect(err).To(BeNil())

		// assert the response code is 0
		Expect(res.Code).To(BeEquivalentTo(0))

		// get the balance of the account after burning the coins
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, didtypes.BaseMinimalDenom)

		// assert no error
		Expect(err).To(BeNil())

		// calculate the difference between the balance before and after burning the coins
		diff := balanceBefore.Sub(balanceAfter)

		// assert the difference is equal to the coins burnt
		total := burnCoins.Add(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(3_500_000_000)))

		// assert the difference is equal to the coins burnt
		Expect(diff).To(Equal(total))
	})

	It("shouldn't burn if sender has insufficient funds", func() {
		// define the coins to burn, in which case 1,000,000 ncheq or 0.01 cheq
		coins := sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(1_000_000))

		// get the balance of the account before burning the coins
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_3_ADDR, didtypes.BaseMinimalDenom)

		// assert no error
		Expect(err).To(BeNil())

		// generate fixed fees, in which case 3,500,000,000 ncheq or 3.5 cheq
		fees := helpers.GenerateFees("3500000000ncheq")

		// burn the coins
		res, err := cli.BurnMsg(testdata.BASE_ACCOUNT_3, coins.String(), fees)

		// assert error
		Expect(err).NotTo(BeNil())

		// assert the response code is 0
		Expect(res.Code).To(BeEquivalentTo(0))

		// get the balance of the account after burning the coins
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_3_ADDR, didtypes.BaseMinimalDenom)

		// assert no error
		Expect(err).To(BeNil())

		// assert the balance before and after burning the coins are equal
		Expect(balanceBefore).To(Equal(balanceAfter))
	})
})

var _ = Describe("Upgrade - Feemarket fees (non-taxable transactions)", func() {
	It("should successfully submit a non-taxable transaction with sufficient fees (--gas-prices)", func() {
		// query feemarket gas price for the base minimal denom
		gasPrice, err := cli.QueryFeemarketGasPrice(didtypes.BaseMinimalDenom)

		// print the gas price
		By("Gas Price: " + gasPrice.Price.String())

		// assert no error
		Expect(err).To(BeNil())

		// define the coins to send, in which case 1,000,000,000 ncheq or 1 cheq
		coins := sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(1_000_000_000))

		// compute gas price, using offset
		gasPrice.Price.Amount = gasPrice.Price.Amount.Mul(sdkmath.LegacyNewDec(didtypes.FeeOffset))

		// define feeParams
		feeParams := []string{
			"--gas", cli.Gas,
			"--gas-adjustment", cli.GasAdjustment,
			"--gas-prices", gasPrice.Price.String(),
		}

		// send the coins, balance assertions are intentionally omitted or out of scope
		res, err := cli.SendTokensTx(testdata.BASE_ACCOUNT_1, testdata.BASE_ACCOUNT_2_ADDR, coins.String(), feeParams)

		// assert no error
		Expect(err).To(BeNil())

		// assert the response code is 0
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should successfully submit a non-taxable transaction with sufficient fees (--fees)", func() {
		// query feemarket gas price for the base minimal denom
		gasPrice, err := cli.QueryFeemarketGasPrice(didtypes.BaseMinimalDenom)

		// print the gas price
		By("Gas Price: " + gasPrice.Price.String())

		// assert no error
		Expect(err).To(BeNil())

		// define the coins to send, in which case 1,000,000,000 ncheq or 1 cheq
		coins := sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(1_000_000_000))

		// define static fees, in which case gas price is multiplied by roughly 3 or greater, times the minimal base denom
		// consider multiplying in the range of [1.5, 3] times the gas price
		gasPrice.Price.Amount = gasPrice.Price.Amount.Mul(sdkmath.LegacyNewDec(3)).Mul(sdkmath.LegacyNewDec(didtypes.BaseFactor))

		// define feeParams
		feeParams := helpers.GenerateFees(gasPrice.Price.String())

		// send the coins, balance assertions are intentionally omitted or out of scope
		res, err := cli.SendTokensTx(testdata.BASE_ACCOUNT_1, testdata.BASE_ACCOUNT_2_ADDR, coins.String(), feeParams)

		// assert no error
		Expect(err).To(BeNil())

		// assert the response code is 0
		Expect(res.Code).To(BeEquivalentTo(0))
	})
})
