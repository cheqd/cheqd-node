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

var _ = Describe("Upgrade - Feemarket fees (non-taxable transactions) negative", func() {
	It("should fail to submit a non-taxable transaction with insufficient fees (--gas-prices)", func() {
		// query feemarket gas price for the base minimal denom
		gasPrice, err := cli.QueryFeemarketGasPrice(didtypes.BaseMinimalDenom)

		// print the gas price
		println("Gas Price: " + gasPrice.String())

		// assert no error
		Expect(err).To(BeNil())

		// define the coins to send, in which case 1,000,000,000 ncheq or 1 cheq
		coins := sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(1_000_000_000))

		// compute gas price, using offset
		gasPrice.Price.Amount = gasPrice.Price.Amount.Mul(sdkmath.LegacyNewDec(didtypes.FeeOffset))

		// invalidate the gas price, in which case 100 times less than the required
		insufficientGasPrice := gasPrice.Price.Amount.Mul(sdkmath.LegacyMustNewDecFromStr("0.01"))

		// define feeParams
		feeParams := []string{
			"--gas", cli.Gas,
			"--gas-adjustment", cli.GasAdjustment,
			"--gas-prices", insufficientGasPrice.String(),
		}

		// send the coins, balance assertions are intentionally omitted or out of scope
		res, err := cli.SendTokensTx(testdata.BASE_ACCOUNT_1, testdata.BASE_ACCOUNT_2_ADDR, coins.String(), feeParams)

		// assert error
		Expect(err).ToNot(BeNil())

		// assert the response code is 13
		Expect(res.Code).To(BeEquivalentTo(13))
	})

	It("should fail to submit a non-taxable transaction with insufficient fees (--fees)", func() {
		// query feemarket gas price for the base minimal denom
		gasPrice, err := cli.QueryFeemarketGasPrice(didtypes.BaseMinimalDenom)

		// print the gas price
		println("Gas Price: " + gasPrice.String())

		// assert no error
		Expect(err).To(BeNil())

		// define the coins to send, in which case 1,000,000,000 ncheq or 1 cheq
		coins := sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(1_000_000_000))

		// define static fees, in which case gas price is multiplied by roughly 3 or greater, times the minimal base denom
		// consider multiplying in the range of [1.5, 3] times the gas price
		gasPrice.Price.Amount = gasPrice.Price.Amount.Mul(sdkmath.LegacyNewDec(3)).Mul(sdkmath.LegacyNewDec(didtypes.BaseFactor))

		// invalidate the static fees, in which case 100 times less than the required
		insufficientGasPrice := gasPrice.Price.Amount.Mul(sdkmath.LegacyMustNewDecFromStr("0.01"))

		// define feeParams
		feeParams := helpers.GenerateFees(insufficientGasPrice.String())

		// send the coins, balance assertions are intentionally omitted or out of scope
		res, err := cli.SendTokensTx(testdata.BASE_ACCOUNT_1, testdata.BASE_ACCOUNT_2_ADDR, coins.String(), feeParams)

		// assert error
		Expect(err).ToNot(BeNil())

		// assert the response code is 13
		Expect(res.Code).To(BeEquivalentTo(13))
	})
})
