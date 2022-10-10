package types_test

import (
	"time"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cheqd/cheqd-node/x/cheqd/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ = Describe(`StaeValue tests`, func() {

	Context("Pack/unpack functionality", func() {
		It("should pack and unpack withour any errors", func() {

			original := &Did{
				Id: "test",
			}
			
			// Construct codec
			registry := types.NewInterfaceRegistry()
			RegisterInterfaces(registry)
			cdc := codec.NewProtoCodec(registry)

			// Marshal
			bz, err := cdc.MarshalInterface(original)
			Expect(err).To(BeNil())

			// Assert type url
			var any types.Any
			err = any.Unmarshal(bz)
			Expect(err).To(BeNil())
			Expect(any.TypeUrl).To(Equal(utils.MsgTypeURL(&Did{})))

			// Unmarshal
			var decoded StateValueData
			err = cdc.UnmarshalInterface(bz, &decoded)
			Expect(err).To(BeNil())
			Expect(&Did{}).To(BeAssignableToTypeOf(decoded))
			Expect(original).To(Equal(decoded))
		})
	})

	When("New metadata is created from context", func() {
		It("should be the same as original", func() {
			createdTime := time.Now()
			ctx := sdk.NewContext(nil, tmproto.Header{ChainID: "test_chain_id", Time: createdTime}, true, nil)
			ctx.WithTxBytes([]byte("test_tx"))
			expectedMetadata := Metadata{
				Created:     createdTime.UTC().Format(time.RFC3339),
				Updated:     "",
				Deactivated: false,
				VersionId:   utils.GetTxHash(ctx.TxBytes()),
			}

			metadata := NewMetadataFromContext(ctx)
			Expect(expectedMetadata).To(Equal(metadata))
		})
	})

	When("Metadata is updated", func() {
		It("should has ", func() {
			createdTime := time.Now()
			updatedTime := createdTime.Add(time.Hour)
		
			ctx1 := NewContext(createdTime, []byte("test1_tx"))
			ctx2 := NewContext(updatedTime, []byte("test1_tx"))
		
			expectedMetadata := Metadata{
				Created:     createdTime.UTC().Format(time.RFC3339),
				Updated:     updatedTime.UTC().Format(time.RFC3339),
				Deactivated: false,
				VersionId:   utils.GetTxHash(ctx2.TxBytes()),
			}
		
			metadata := NewMetadataFromContext(ctx1)
			metadata.Update(ctx2)
		
			Expect(expectedMetadata).To(Equal(metadata))
		})
	})

})

func NewContext(time time.Time, txBytes []byte) sdk.Context {
	ctx := sdk.NewContext(nil, tmproto.Header{ChainID: "test_chain_id", Time: time}, true, nil)
	ctx.WithTxBytes(txBytes)
	return ctx
}
