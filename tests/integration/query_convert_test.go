//go:build integration

package integration

import (
	"fmt"
	"time"

	cli "github.com/cheqd/cheqd-node/tests/integration/cli"
	oraclekeeper "github.com/cheqd/cheqd-node/x/oracle/keeper"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration - Convert USDC to CHEQ using MA", func() {
	var resp *oracletypes.ConvertUSDCtoCHEQResponse
	var err error
	var parsedAmount sdk.Coin

	BeforeEach(func() {
		oracleParams, err := cli.QueryOracleParams()
		Expect(err).To(BeNil())

		historicStampPeriod := oracleParams.Params.HistoricStampPeriod
		averagingWindow := oraclekeeper.AveragingWindow

		targetHeight := int64(historicStampPeriod) * int64(averagingWindow)
		fmt.Printf("Waiting until block height ≥ %d to trigger ComputeAllAverages...\n", targetHeight)

		for {
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			if currentHeight >= targetHeight {
				break
			}
			time.Sleep(2 * time.Second)
		}
	})

	It("should convert using SMA", func() {
		By("executing SMA query")
		resp, err = cli.QueryConvertUSDCtoCHEQ("1000000000000000000usd", "sma", "", nil)
		Expect(err).To(BeNil())
		parsedAmount, err = sdk.ParseCoinNormalized(resp.Amount)
		Expect(err).To(BeNil())
		Expect(parsedAmount.Denom).To(Equal("ncheq"))
		Expect(parsedAmount.Amount.Int64()).To(BeNumerically(">", 0))
	})

	It("should convert using EMA", func() {
		By("executing EMA query")
		resp, err = cli.QueryConvertUSDCtoCHEQ("2000000000000000000usd", "ema", "", nil)
		Expect(err).To(BeNil())
		parsedAmount, err = sdk.ParseCoinNormalized(resp.Amount)
		Expect(err).To(BeNil())
		Expect(parsedAmount.Denom).To(Equal("ncheq"))
		Expect(parsedAmount.Amount.Int64()).To(BeNumerically(">", 0))
	})

	It("should convert using WMA RECENT", func() {
		By("executing WMA RECENT query")
		resp, err = cli.QueryConvertUSDCtoCHEQ("3000000000000000000usd", "wma", "RECENT", nil)
		Expect(err).To(BeNil())
		parsedAmount, err = sdk.ParseCoinNormalized(resp.Amount)
		Expect(err).To(BeNil())
		Expect(parsedAmount.Denom).To(Equal("ncheq"))
		Expect(parsedAmount.Amount.Int64()).To(BeNumerically(">", 0))
	})

	It("should convert using WMA OLDEST", func() {
		By("executing WMA OLDEST query")
		resp, err = cli.QueryConvertUSDCtoCHEQ("4000000000000000000usd", "wma", "OLDEST", nil)
		Expect(err).To(BeNil())
		parsedAmount, err = sdk.ParseCoinNormalized(resp.Amount)
		Expect(err).To(BeNil())
		Expect(parsedAmount.Denom).To(Equal("ncheq"))
		Expect(parsedAmount.Amount.Int64()).To(BeNumerically(">", 0))
	})

	It("should convert using WMA BALANCED", func() {
		By("executing WMA BALANCED query")
		resp, err = cli.QueryConvertUSDCtoCHEQ("5000000000000000000usd", "wma", "BALANCED", nil)
		Expect(err).To(BeNil())
		parsedAmount, err = sdk.ParseCoinNormalized(resp.Amount)
		Expect(err).To(BeNil())
		Expect(parsedAmount.Denom).To(Equal("ncheq"))
		Expect(parsedAmount.Amount.Int64()).To(BeNumerically(">", 0))
	})

	It("should convert using WMA CUSTOM with matching weights", func() {
		By("executing WMA CUSTOM query with valid weights")
		weights := []int64{1, 2, 3}
		resp, err = cli.QueryConvertUSDCtoCHEQ("6000000000000000000usd", "wma", "CUSTOM", weights)
		Expect(err).To(BeNil())
		parsedAmount, err = sdk.ParseCoinNormalized(resp.Amount)
		Expect(err).To(BeNil())
		Expect(parsedAmount.Denom).To(Equal("ncheq"))
		Expect(parsedAmount.Amount.Int64()).To(BeNumerically(">", 0))
	})

	It("should fail WMA CUSTOM with mismatched weights", func() {
		By("executing WMA CUSTOM query with too few weights")
		weights := []int64{1, 2}
		_, err = cli.QueryConvertUSDCtoCHEQ("7000000000000000000usd", "wma", "CUSTOM", weights)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("weights length"))
	})

	It("should fail on unsupported MA type", func() {
		By("executing query with invalid ma_type")
		_, err = cli.QueryConvertUSDCtoCHEQ("8000000000000000000usd", "badma", "", nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid MA type"))
	})

	It("should fail on unsupported denom", func() {
		By("executing query with invalid denom")
		_, err = cli.QueryConvertUSDCtoCHEQ("1000000000000000000atom", "sma", "", nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("expected denom to be 'usd'"))
	})
})
