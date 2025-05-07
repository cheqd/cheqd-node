package ante_test

import (
	"encoding/base64"

	"cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/ante"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourceutils "github.com/cheqd/cheqd-node/x/resource/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TxMsgFilters", func() {
	rounds := 1_000_000

	BeforeEach(func() {
		ante.TaxableMsgFees = ante.TaxableMsgFee{
			ante.MsgCreateDidDoc:          sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(didtypes.DefaultCreateDidTxFee))),
			ante.MsgUpdateDidDoc:          sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(didtypes.DefaultUpdateDidTxFee))),
			ante.MsgDeactivateDidDoc:      sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(didtypes.DefaultDeactivateDidTxFee))),
			ante.MsgCreateResourceDefault: sdk.NewCoins(sdk.NewCoin(resourcetypes.BaseMinimalDenom, math.NewInt(resourcetypes.DefaultCreateResourceDefaultFee))),
			ante.MsgCreateResourceImage:   sdk.NewCoins(sdk.NewCoin(resourcetypes.BaseMinimalDenom, math.NewInt(resourcetypes.DefaultCreateResourceImageFee))),
			ante.MsgCreateResourceJSON:    sdk.NewCoins(sdk.NewCoin(resourcetypes.BaseMinimalDenom, math.NewInt(resourcetypes.DefaultCreateResourceJSONFee))),
		}

		ante.BurnFactors = ante.BurnFactor{
			ante.BurnFactorDid:      math.LegacyMustNewDecFromStr("0.990000000000000000"),
			ante.BurnFactorResource: math.LegacyMustNewDecFromStr("0.990000000000000000"),
		}
	})

	Describe("GetResourceTaxableMsgFee", func() {
		It("should return the correct fee for image mimetype - 1mn rounds", func() {
			// define byte content, base64-encoded png
			content, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAABEAAAAOCAMAAAD+MweGAAADAFBMVEUAAAAAAFUAAKoAAP8AJAAAJFUAJKoAJP8ASQAASVUASaoASf8AbQAAbVUAbaoAbf8AkgAAklUAkqoAkv8AtgAAtlUAtqoAtv8A2wAA21UA26oA2/8A/wAA/1UA/6oA//8kAAAkAFUkAKokAP8kJAAkJFUkJKokJP8kSQAkSVUkSaokSf8kbQAkbVUkbaokbf8kkgAkklUkkqokkv8ktgAktlUktqoktv8k2wAk21Uk26ok2/8k/wAk/1Uk/6ok//9JAABJAFVJAKpJAP9JJABJJFVJJKpJJP9JSQBJSVVJSapJSf9JbQBJbVVJbapJbf9JkgBJklVJkqpJkv9JtgBJtlVJtqpJtv9J2wBJ21VJ26pJ2/9J/wBJ/1VJ/6pJ//9tAABtAFVtAKptAP9tJABtJFVtJKptJP9tSQBtSVVtSaptSf9tbQBtbVVtbaptbf9tkgBtklVtkqptkv9ttgBttlVttqpttv9t2wBt21Vt26pt2/9t/wBt/1Vt/6pt//+SAACSAFWSAKqSAP+SJACSJFWSJKqSJP+SSQCSSVWSSaqSSf+SbQCSbVWSbaqSbf+SkgCSklWSkqqSkv+StgCStlWStqqStv+S2wCS21WS26qS2/+S/wCS/1WS/6qS//+2AAC2AFW2AKq2AP+2JAC2JFW2JKq2JP+2SQC2SVW2Saq2Sf+2bQC2bVW2baq2bf+2kgC2klW2kqq2kv+2tgC2tlW2tqq2tv+22wC221W226q22/+2/wC2/1W2/6q2///bAADbAFXbAKrbAP/bJADbJFXbJKrbJP/bSQDbSVXbSarbSf/bbQDbbVXbbarbbf/bkgDbklXbkqrbkv/btgDbtlXbtqrbtv/b2wDb21Xb26rb2//b/wDb/1Xb/6rb////AAD/AFX/AKr/AP//JAD/JFX/JKr/JP//SQD/SVX/Sar/Sf//bQD/bVX/bar/bf//kgD/klX/kqr/kv//tgD/tlX/tqr/tv//2wD/21X/26r/2////wD//1X//6r////qm24uAAAA1ElEQVR42h1PMW4CQQwc73mlFJGCQChFIp0Rh0RBGV5AFUXKC/KPfCFdqryEgoJ8IX0KEF64q0PPnow3jT2WxzNj+gAgAGfvvDdCQIHoSnGYcGDE2nH92DoRqTYJ2bTcsKgqhIi47VdgAWNmwFSFA1UAAT2sSFcnq8a3x/zkkJrhaHT3N+hD3aH7ZuabGHX7bsSMhxwTJLr3evf1e0nBVcwmqcTZuatKoJaB7dSHjTZdM0G1HBTWefly//q2EB7/BEvk5vmzeQaJ7/xKPImpzv8/s4grhAxHl0DsqGUAAAAASUVORK5CYII=")
			Expect(err).To(BeNil())

			// perform holistic action benchmark - 100k rounds
			for range rounds {
				// detect mime type
				mimeType := resourceutils.DetectMediaType(content)
				Expect(mimeType).To(BeEquivalentTo("image/png"))

				// create indicative sign info
				signInfo := didtypes.SignInfo{
					VerificationMethodId: "",
					Signature:            []byte(""),
				}

				// create indicative signatures
				signatures := make([]*didtypes.SignInfo, 1)

				// append sign info
				signatures = append(signatures, &signInfo)

				// create resource
				resourceMsg := resourcetypes.MsgCreateResource{
					Payload: &resourcetypes.MsgCreateResourcePayload{
						Id:           "",
						CollectionId: "",
						ResourceType: "",
						Name:         "",
						Version:      "",
						Data:         content,
					},
					Signatures: signatures,
				}

				// calculate portions
				reward, burn, ok := ante.GetResourceTaxableMsgFee(sdk.Context{}, &resourceMsg)
				Expect(ok).To(BeTrue())
				Expect(reward).To(Equal(ante.GetRewardPortion(ante.TaxableMsgFees[ante.MsgCreateResourceImage], ante.GetBurnFeePortion(ante.BurnFactors[ante.BurnFactorResource], ante.TaxableMsgFees[ante.MsgCreateResourceImage]))))
				Expect(burn).To(Equal(ante.GetBurnFeePortion(ante.BurnFactors[ante.BurnFactorResource], ante.TaxableMsgFees[ante.MsgCreateResourceImage])))
				Expect(reward.Add(burn[0])[0].Amount).To(BeEquivalentTo(ante.TaxableMsgFees[ante.MsgCreateResourceImage][0].Amount))
			}
		})

		It("should return the correct fee for JSON mimetype - 1mn rounds", func() {
			// define JSON content
			content := []byte(`{"key": "value"}`)

			// perform holistic action benchmark - 1mn rounds
			for range rounds {
				// detect mime type
				mimeType := resourceutils.DetectMediaType(content)
				Expect(mimeType).To(BeEquivalentTo("application/json"))

				// create indicative sign info
				signInfo := didtypes.SignInfo{
					VerificationMethodId: "",
					Signature:            []byte(""),
				}

				// create indicative signatures
				signatures := make([]*didtypes.SignInfo, 1)

				// append sign info
				signatures = append(signatures, &signInfo)

				// create resource
				resourceMsg := resourcetypes.MsgCreateResource{
					Payload: &resourcetypes.MsgCreateResourcePayload{
						Id:           "",
						CollectionId: "",
						ResourceType: "",
						Name:         "",
						Version:      "",
						Data:         content,
					},
					Signatures: signatures,
				}

				// calculate portions
				reward, burn, ok := ante.GetResourceTaxableMsgFee(sdk.Context{}, &resourceMsg)
				Expect(ok).To(BeTrue())
				Expect(reward).To(Equal(ante.GetRewardPortion(ante.TaxableMsgFees[ante.MsgCreateResourceJSON], ante.GetBurnFeePortion(ante.BurnFactors[ante.BurnFactorResource], ante.TaxableMsgFees[ante.MsgCreateResourceJSON]))))
				Expect(burn).To(Equal(ante.GetBurnFeePortion(ante.BurnFactors[ante.BurnFactorResource], ante.TaxableMsgFees[ante.MsgCreateResourceJSON])))
				Expect(reward.Add(burn[0])[0].Amount).To(BeEquivalentTo(ante.TaxableMsgFees[ante.MsgCreateResourceJSON][0].Amount))
			}
		})

		It("should return the correct fee for default mimetype - 1mn rounds", func() {
			// define byte content, base64-encoded .txt
			content, err := base64.StdEncoding.DecodeString("VGhpcyBpcyBhIHRlc3QgdGV4dCBmaWxlLg==")
			Expect(err).To(BeNil())

			// perform holistic action benchmark - 1mn rounds
			for range rounds {
				// detect mime type
				mimeType := resourceutils.DetectMediaType(content)
				Expect(mimeType).To(BeEquivalentTo("text/plain; charset=utf-8"))

				// create indicative sign info
				signInfo := didtypes.SignInfo{
					VerificationMethodId: "",
					Signature:            []byte(""),
				}

				// create indicative signatures
				signatures := make([]*didtypes.SignInfo, 1)

				// append sign info
				signatures = append(signatures, &signInfo)

				// create resource
				resourceMsg := resourcetypes.MsgCreateResource{
					Payload: &resourcetypes.MsgCreateResourcePayload{
						Id:           "",
						CollectionId: "",
						ResourceType: "",
						Name:         "",
						Version:      "",
						Data:         content,
					},
					Signatures: signatures,
				}

				// calculate portions
				reward, burn, ok := ante.GetResourceTaxableMsgFee(sdk.Context{}, &resourceMsg)
				Expect(ok).To(BeTrue())
				Expect(reward).To(Equal(ante.GetRewardPortion(ante.TaxableMsgFees[ante.MsgCreateResourceDefault], ante.GetBurnFeePortion(ante.BurnFactors[ante.BurnFactorResource], ante.TaxableMsgFees[ante.MsgCreateResourceDefault]))))
				Expect(burn).To(Equal(ante.GetBurnFeePortion(ante.BurnFactors[ante.BurnFactorResource], ante.TaxableMsgFees[ante.MsgCreateResourceDefault])))
				Expect(reward.Add(burn[0])[0].Amount).To(BeEquivalentTo(ante.TaxableMsgFees[ante.MsgCreateResourceDefault][0].Amount))
			}
		})
	})
})
