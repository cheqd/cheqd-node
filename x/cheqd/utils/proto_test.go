package utils_test

import (
	. "github.com/cheqd/cheqd-node/x/cheqd/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	bank_types "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var _ = Describe("Proto", func() {
	Describe("Check that DenomUnit from bank type has expected value", func() {
		Context("Denom Unit from cosmos-sdk" ,func ()  {
			It("should return expected value", func() {
				Expect(MsgTypeURL(&bank_types.DenomUnit{})).To(Equal("/cosmos.bank.v1beta1.DenomUnit"))
			})
		})
	})
})
