package keeper_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/x/oracle/abci"
	"github.com/cheqd/cheqd-node/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenerateSalt generates a random salt, size length/2,  as a HEX encoded string.
func GenerateSalt(length int) (string, error) {
	if length == 0 {
		return "", fmt.Errorf("failed to generate salt: zero length")
	}

	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (s *IntegrationTestSuite) TestMsgServer_AggregateExchangeRatePrevote() {
	ctx := s.ctx

	exchangeRatesStr := "123.2:CHEQ"
	salt, err := GenerateSalt(32)
	s.Require().NoError(err)
	hash := types.GetAggregateVoteHash(salt, exchangeRatesStr, valAddr)

	invalidHash := &types.MsgAggregateExchangeRatePrevote{
		Hash:      "invalid_hash",
		Feeder:    addr.String(),
		Validator: valAddr.String(),
	}
	invalidFeeder := &types.MsgAggregateExchangeRatePrevote{
		Hash:      hash.String(),
		Feeder:    "invalid_feeder",
		Validator: valAddr.String(),
	}
	invalidValidator := &types.MsgAggregateExchangeRatePrevote{
		Hash:      hash.String(),
		Feeder:    addr.String(),
		Validator: "invalid_val",
	}
	validMsg := &types.MsgAggregateExchangeRatePrevote{
		Hash:      hash.String(),
		Feeder:    addr.String(),
		Validator: valAddr.String(),
	}

	_, err = s.msgServer.AggregateExchangeRatePrevote(ctx, invalidHash)
	s.Require().Error(err)
	_, err = s.msgServer.AggregateExchangeRatePrevote(ctx, invalidFeeder)
	s.Require().Error(err)
	_, err = s.msgServer.AggregateExchangeRatePrevote(ctx, invalidValidator)
	s.Require().Error(err)
	_, err = s.msgServer.AggregateExchangeRatePrevote(ctx, validMsg)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestMsgServer_AggregateExchangeRateVote() {
	ctx := s.ctx

	ratesStr := "CHEQ:123.2"
	ratesStrInvalidCoin := "CHEQ:123.2,badcoin:234.5"
	salt, err := GenerateSalt(32)
	s.Require().NoError(err)
	hash := types.GetAggregateVoteHash(salt, ratesStr, valAddr)
	hashInvalidRate := types.GetAggregateVoteHash(salt, ratesStrInvalidCoin, valAddr)

	prevoteMsg := &types.MsgAggregateExchangeRatePrevote{
		Hash:      hash.String(),
		Feeder:    addr.String(),
		Validator: valAddr.String(),
	}
	voteMsg := &types.MsgAggregateExchangeRateVote{
		Feeder:        addr.String(),
		Validator:     valAddr.String(),
		Salt:          salt,
		ExchangeRates: ratesStr,
	}
	voteMsgInvalidRate := &types.MsgAggregateExchangeRateVote{
		Feeder:        addr.String(),
		Validator:     valAddr.String(),
		Salt:          salt,
		ExchangeRates: ratesStrInvalidCoin,
	}

	// Flattened acceptList symbols to make checks easier
	acceptList := s.app.OracleKeeper.GetParams(ctx).AcceptList
	var acceptListFlat []string
	for _, v := range acceptList {
		acceptListFlat = append(acceptListFlat, v.SymbolDenom)
	}

	// No existing prevote
	_, err = s.msgServer.AggregateExchangeRateVote(ctx, voteMsg)
	s.Require().EqualError(err, errors.Wrap(types.ErrNoAggregatePrevote, valAddr.String()).Error())
	_, err = s.msgServer.AggregateExchangeRatePrevote(ctx, prevoteMsg)
	s.Require().NoError(err)
	// Reveal period mismatch
	_, err = s.msgServer.AggregateExchangeRateVote(ctx, voteMsg)
	s.Require().EqualError(err, types.ErrRevealPeriodMissMatch.Error())

	// Valid
	s.app.OracleKeeper.SetAggregateExchangeRatePrevote(
		ctx,
		valAddr,
		types.NewAggregateExchangeRatePrevote(
			hash, valAddr, 8,
		))
	_, err = s.msgServer.AggregateExchangeRateVote(ctx, voteMsg)
	s.Require().NoError(err)
	vote, err := s.app.OracleKeeper.GetAggregateExchangeRateVote(ctx, valAddr)
	s.Require().Nil(err)
	for _, v := range vote.ExchangeRates {
		s.Require().Contains(acceptListFlat, v.Denom)
	}

	// Valid, but with an exchange rate which isn't in AcceptList
	s.app.OracleKeeper.SetAggregateExchangeRatePrevote(
		ctx,
		valAddr,
		types.NewAggregateExchangeRatePrevote(
			hashInvalidRate, valAddr, 8,
		))
	_, err = s.msgServer.AggregateExchangeRateVote(ctx, voteMsgInvalidRate)
	s.Require().NoError(err)
	vote, err = s.app.OracleKeeper.GetAggregateExchangeRateVote(ctx, valAddr)
	s.Require().NoError(err)
	for _, v := range vote.ExchangeRates {
		s.Require().Contains(acceptListFlat, v.Denom)
	}
}

func (s *IntegrationTestSuite) TestMsgServer_DelegateFeedConsent() {
	app, ctx := s.app, s.ctx

	feederAddr := sdk.AccAddress([]byte("addr________________"))
	feederAcc := app.AccountKeeper.NewAccountWithAddress(ctx, feederAddr)
	app.AccountKeeper.SetAccount(ctx, feederAcc)

	_, err := s.msgServer.DelegateFeedConsent(ctx, &types.MsgDelegateFeedConsent{
		Operator: valAddr.String(),
		Delegate: feederAddr.String(),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestMsgServer_UpdateGovParams() {
	govAccAddr := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()
	testCases := []struct {
		name      string
		req       *types.MsgGovUpdateParams
		expectErr bool
		errMsg    string
	}{
		{
			"valid accept list",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Plan: types.ParamUpdatePlan{
					Keys:   []string{"AcceptList"},
					Height: 9,
					Changes: types.Params{
						AcceptList: append(types.DefaultAcceptList, types.Denom{
							BaseDenom:   "base",
							SymbolDenom: "symbol",
							Exponent:    6,
						}),
					},
				},
			},
			false,
			"",
		},
		{
			"valid mandatory list",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Plan: types.ParamUpdatePlan{
					Keys:   []string{"MandatoryList"},
					Height: 9,
					Changes: types.Params{
						MandatoryList: types.DefaultMandatoryList,
					},
				},
			},
			false,
			"",
		},
		{
			"invalid mandatory list",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Plan: types.ParamUpdatePlan{
					Keys:   []string{"MandatoryList"},
					Height: 9,
					Changes: types.Params{
						MandatoryList: types.DenomList{
							{
								BaseDenom:   "test",
								SymbolDenom: "test",
								Exponent:    6,
							},
						},
					},
				},
			},
			true,
			"denom in MandatoryList not present in AcceptList",
		},
		{
			"valid reward band list",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Plan: types.ParamUpdatePlan{
					Keys:   []string{"RewardBands"},
					Height: 9,
					Changes: types.Params{
						RewardBands: append(types.DefaultRewardBands(), types.RewardBand{
							SymbolDenom: "symbol",
							RewardBand:  math.LegacyNewDecWithPrec(2, 2),
						}),
					},
				},
			},
			false,
			"",
		},
		{
			"invalid reward band list",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Plan: types.ParamUpdatePlan{
					Keys:   []string{"RewardBands"},
					Height: 9,
					Changes: types.Params{
						RewardBands: types.RewardBandList{
							{
								SymbolDenom: types.CheqdSymbol,
								RewardBand:  math.LegacyNewDecWithPrec(2, 0),
							},
							{
								SymbolDenom: types.AtomSymbol,
								RewardBand:  math.LegacyNewDecWithPrec(2, 2),
							},
						},
					},
				},
			},
			true,
			"oracle parameter RewardBand must be between [0, 1]",
		},
		{
			"multiple valid params",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Plan: types.ParamUpdatePlan{
					Keys: []string{
						"VotePeriod",
						"VoteThreshold",
						"RewardDistributionWindow",
						"SlashFraction",
						"SlashWindow",
						"MinValidPerWindow",
						"HistoricStampPeriod",
						"MedianStampPeriod",
						"MaximumPriceStamps",
						"MaximumMedianStamps",
					},
					Height: 9,
					Changes: types.Params{
						VotePeriod:               10,
						VoteThreshold:            math.LegacyNewDecWithPrec(40, 2),
						RewardDistributionWindow: types.BlocksPerWeek,
						SlashFraction:            math.LegacyNewDecWithPrec(2, 4),
						SlashWindow:              types.BlocksPerDay,
						MinValidPerWindow:        math.LegacyNewDecWithPrec(4, 2),
						HistoricStampPeriod:      10 * types.BlocksPerMinute,
						MedianStampPeriod:        5 * types.BlocksPerHour,
						MaximumPriceStamps:       40,
						MaximumMedianStamps:      30,
					},
				},
			},
			false,
			"",
		},
		{
			"invalid vote threshold",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Plan: types.ParamUpdatePlan{
					Keys:   []string{"VoteThreshold"},
					Height: 9,
					Changes: types.Params{
						VoteThreshold: math.LegacyNewDecWithPrec(10, 2),
					},
				},
			},
			true,
			"threshold must be bigger than 0.330000000000000000 and <= 1",
		},
		{
			"invalid slash window",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Plan: types.ParamUpdatePlan{
					Keys:   []string{"VotePeriod", "SlashWindow"},
					Height: 9,
					Changes: types.Params{
						VotePeriod:  5,
						SlashWindow: 4,
					},
				},
			},
			true,
			"oracle parameter SlashWindow must be greater than or equal with VotePeriod",
		},
		{
			"invalid key",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Plan: types.ParamUpdatePlan{
					Keys:    []string{"test"},
					Height:  9,
					Changes: types.Params{},
				},
			},
			true,
			"test is not an existing oracle param key",
		},
		{
			"bad authority",
			&types.MsgGovUpdateParams{
				Authority:   "cheqd1jvfm45nwpgvxg4m2r5l80j5my87uc54mkzsg6h",
				Title:       "test",
				Description: "test",
				Plan: types.ParamUpdatePlan{
					Keys:   []string{"RewardBands"},
					Height: 9,
					Changes: types.Params{
						RewardBands: types.RewardBandList{
							{
								SymbolDenom: types.CheqdSymbol,
								RewardBand:  math.LegacyNewDecWithPrec(2, 2),
							},
						},
					},
				},
			},
			true,
			"invalid authority",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := tc.req.ValidateBasic()
			if err == nil {
				_, err = s.msgServer.GovUpdateParams(s.ctx, tc.req)
				err1 := abci.EndBlocker(s.ctx, s.app.OracleKeeper, s.app.FeeabsKeeper)
				s.Require().NoError(err1)
			}
			if tc.expectErr {
				s.Require().ErrorContains(err, tc.errMsg)
			} else {
				s.Require().NoError(err)

				switch tc.name {
				case "valid accept list":
					acceptList := s.app.OracleKeeper.AcceptList(s.ctx)
					s.Require().Equal(acceptList, append(types.DefaultAcceptList, types.Denom{
						BaseDenom:   "base",
						SymbolDenom: "symbol",
						Exponent:    6,
					}).Normalize())

				case "valid mandatory list":
					mandatoryList := s.app.OracleKeeper.MandatoryList(s.ctx)
					s.Require().Equal(mandatoryList, types.DefaultMandatoryList.Normalize())

				case "valid reward band list":
					rewardBand := s.app.OracleKeeper.RewardBands(s.ctx)
					s.Require().Equal(rewardBand, append(types.DefaultRewardBands(), types.RewardBand{
						SymbolDenom: "symbol",
						RewardBand:  math.LegacyNewDecWithPrec(2, 2),
					}))

				case "multiple valid params":
					votePeriod := s.app.OracleKeeper.VotePeriod(s.ctx)
					voteThreshold := s.app.OracleKeeper.VoteThreshold(s.ctx)
					rewardDistributionWindow := s.app.OracleKeeper.RewardDistributionWindow(s.ctx)
					slashFraction := s.app.OracleKeeper.SlashFraction(s.ctx)
					slashWindow := s.app.OracleKeeper.SlashWindow(s.ctx)
					minValidPerWindow := s.app.OracleKeeper.MinValidPerWindow(s.ctx)
					historicStampPeriod := s.app.OracleKeeper.HistoricStampPeriod(s.ctx)
					medianStampPeriod := s.app.OracleKeeper.MedianStampPeriod(s.ctx)
					maximumPriceStamps := s.app.OracleKeeper.MaximumPriceStamps(s.ctx)
					maximumMedianStamps := s.app.OracleKeeper.MaximumMedianStamps(s.ctx)
					s.Require().Equal(votePeriod, uint64(10))
					s.Require().Equal(voteThreshold, math.LegacyNewDecWithPrec(40, 2))
					s.Require().Equal(rewardDistributionWindow, types.BlocksPerWeek)
					s.Require().Equal(slashFraction, math.LegacyNewDecWithPrec(2, 4))
					s.Require().Equal(slashWindow, types.BlocksPerDay)
					s.Require().Equal(minValidPerWindow, math.LegacyNewDecWithPrec(4, 2))
					s.Require().Equal(historicStampPeriod, 10*types.BlocksPerMinute)
					s.Require().Equal(medianStampPeriod, 5*types.BlocksPerHour)
					s.Require().Equal(maximumPriceStamps, uint64(40))
					s.Require().Equal(maximumMedianStamps, uint64(30))
				}
			}
		})
	}
}

func (s *IntegrationTestSuite) TestMsgServer_GovAddDenom() {
	govAccAddr := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()
	bandArgument := math.LegacyNewDecWithPrec(2, 3)
	foo := &types.Denom{
		SymbolDenom: "FOO",
		BaseDenom:   "FOO",
		Exponent:    6,
	}
	bar := &types.Denom{
		SymbolDenom: "BAR",
		BaseDenom:   "BAR",
		Exponent:    6,
	}
	foobar := &types.Denom{
		SymbolDenom: "FOOBAR",
		BaseDenom:   "FOOBAR",
		Exponent:    6,
	}
	reward := &types.Denom{
		SymbolDenom: "REWARD",
		BaseDenom:   "REWARD",
		Exponent:    6,
	}
	currencyPairProviders := types.CurrencyPairProvidersList{
		{
			BaseDenom:  "FOO",
			QuoteDenom: "BAR",
			PairAddress: []types.PairAddressProvider{
				{
					Address:         "address",
					AddressProvider: "provider",
				},
			},
			Providers: []string{
				"provider",
			},
		},
	}
	currencyDeviationThresholds := types.CurrencyDeviationThresholdList{
		{
			BaseDenom: "FOO",
			Threshold: "2.0",
		},
		{
			BaseDenom: "BAR",
			Threshold: "2.0",
		},
	}
	currencyPairProviders2 := types.CurrencyPairProvidersList{
		{
			BaseDenom:  "FOOBAR",
			QuoteDenom: "BAR",
			PairAddress: []types.PairAddressProvider{
				{
					Address:         "address2",
					AddressProvider: "provider2",
				},
			},
			Providers: []string{
				"provider2",
			},
		},
	}
	currencyDeviationThresholds2 := types.CurrencyDeviationThresholdList{
		{
			BaseDenom: "FOOBAR",
			Threshold: "2.0",
		},
	}

	testCases := []struct {
		name      string
		req       *types.MsgGovAddDenoms
		expectErr bool
		errMsg    string
	}{
		{
			"valid denom addition",
			&types.MsgGovAddDenoms{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,
				DenomList:   append(types.DenomList{}, *foo, *bar),
				Mandatory:   false,
			},
			false,
			"",
		},
		{
			"valid mandatory denom addition with currency pair providers and currency deviation thresholds",
			&types.MsgGovAddDenoms{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,
				DenomList:   append(types.DenomList{}, *foo, *bar),
				Mandatory:   true,

				CurrencyPairProviders:       currencyPairProviders,
				CurrencyDeviationThresholds: currencyDeviationThresholds,
			},
			false,
			"",
		},
		{
			"valid denom addition with reward band",
			&types.MsgGovAddDenoms{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,
				DenomList:   append(types.DenomList{}, *reward),
				Mandatory:   true,
				RewardBand:  &bandArgument,
			},
			false,
			"",
		},
		{
			"valid currency pair providers and currency deviation thresholds addition with no new denoms",
			&types.MsgGovAddDenoms{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,
				DenomList:   types.DenomList{},
				Mandatory:   true,

				CurrencyPairProviders:       currencyPairProviders2,
				CurrencyDeviationThresholds: currencyDeviationThresholds2,
			},
			false,
			"",
		},
		{
			"invalid multiple addition",
			&types.MsgGovAddDenoms{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,
				DenomList:   append(types.DenomList{}, *foo, *foo),
				Mandatory:   false,
			},
			true,
			"denom already exists in acceptList: FOO",
		},
		{
			"invalid multiple addition mandatory",
			&types.MsgGovAddDenoms{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,
				DenomList:   append(types.DenomList{}, *bar, *bar),
				Mandatory:   true,
			},
			true,
			"denom already exists in mandatoryList: BAR",
		},
		{
			"invalid existing addition",
			&types.MsgGovAddDenoms{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,
				DenomList: append(types.DenomList{}, types.Denom{
					SymbolDenom: "CHEQ",
					BaseDenom:   "CHEQ",
					Exponent:    6,
				}),
				Mandatory: false,
			},
			true,
			"denom already exists in acceptList: CHEQ",
		},
		{
			"invalid existing mandatory addition",
			&types.MsgGovAddDenoms{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,
				DenomList: append(types.DenomList{}, types.Denom{
					SymbolDenom: "USDC",
					BaseDenom:   "USDC",
					Exponent:    6,
				}),
				Mandatory: true,
			},
			true,
			"denom already exists in mandatoryList: USDC",
		},
		{
			"invalid empty denom",
			&types.MsgGovAddDenoms{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,
				DenomList:   append(types.DenomList{}, types.Denom{}),
				Mandatory:   true,
			},
			true,
			"invalid oracle param value",
		},
		{
			"invalid currency pair provider list",
			&types.MsgGovAddDenoms{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,
				DenomList:   append(types.DenomList{}, *foobar),
				Mandatory:   true,

				CurrencyPairProviders: types.CurrencyPairProvidersList{
					{
						BaseDenom:  "FOOBAR",
						QuoteDenom: "BAR",
						Providers:  []string{},
					},
				},
			},
			true,
			"oracle parameter CurrencyPairProviders must have at least 1 provider listed",
		},
		{
			"invalid currency deviation threshold list",
			&types.MsgGovAddDenoms{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,
				DenomList:   append(types.DenomList{}, *foobar),
				Mandatory:   true,

				CurrencyDeviationThresholds: types.CurrencyDeviationThresholdList{
					{
						BaseDenom: "FOOBAR",
					},
				},
			},
			true,
			"oracle parameter CurrencyDeviationThreshold must have Threshold",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := tc.req.ValidateBasic()
			if err == nil {
				_, err = s.msgServer.GovAddDenoms(s.ctx, tc.req)
				s.Require().NoError(err)
				err = abci.EndBlocker(s.ctx, s.app.OracleKeeper, s.app.FeeabsKeeper)
				s.Require().NoError(err)

			}
			if tc.expectErr {
				s.Require().ErrorContains(err, tc.errMsg)
			} else {
				s.Require().NoError(err)

				switch tc.name {
				case "valid denom addition":
					al := s.app.OracleKeeper.AcceptList(s.ctx)
					s.Require().True(
						al.Contains(foo.SymbolDenom) && al.Contains(bar.SymbolDenom),
					)

					rwb := s.app.OracleKeeper.RewardBands(s.ctx)
					band, err := rwb.GetBandFromDenom("foo")
					s.Require().Equal(band, math.LegacyNewDecWithPrec(2, 2))
					s.Require().NoError(err)
					band, err = rwb.GetBandFromDenom("bar")
					s.Require().Equal(band, math.LegacyNewDecWithPrec(2, 2))
					s.Require().NoError(err)

				case "valid mandatory denom addition with currency pair providers and currency deviation thresholds":
					al := s.app.OracleKeeper.AcceptList(s.ctx)
					s.Require().True(
						al.Contains(foo.SymbolDenom) && al.Contains(bar.SymbolDenom),
					)
					ml := s.app.OracleKeeper.MandatoryList(s.ctx)
					s.Require().True(
						ml.Contains(foo.SymbolDenom) && ml.Contains(bar.SymbolDenom),
					)

					rwb := s.app.OracleKeeper.RewardBands(s.ctx)
					band, err := rwb.GetBandFromDenom("foo")
					s.Require().Equal(band, math.LegacyNewDecWithPrec(2, 2))
					s.Require().NoError(err)
					band, err = rwb.GetBandFromDenom("bar")
					s.Require().Equal(band, math.LegacyNewDecWithPrec(2, 2))
					s.Require().NoError(err)

					cpp := s.app.OracleKeeper.CurrencyPairProviders(s.ctx)
					for i := range currencyPairProviders {
						s.Require().Contains(cpp, currencyPairProviders[i])
					}

					cdt := s.app.OracleKeeper.CurrencyDeviationThresholds(s.ctx)
					for i := range currencyDeviationThresholds {
						s.Require().Contains(cdt, currencyDeviationThresholds[i])
					}

				case "valid denom addition with reward band":
					al := s.app.OracleKeeper.AcceptList(s.ctx)
					s.Require().True(
						al.Contains("REWARD"),
					)

					rwb := s.app.OracleKeeper.RewardBands(s.ctx)
					band, err := rwb.GetBandFromDenom("REWARD")
					s.Require().Equal(band, math.LegacyNewDecWithPrec(2, 3))
					s.Require().NoError(err)

				case "valid currency pair providers and currency deviation thresholds addition with no new denoms":
					cpp := s.app.OracleKeeper.CurrencyPairProviders(s.ctx)
					for i := range currencyPairProviders2 {
						s.Require().Contains(cpp, currencyPairProviders2[i])
					}

					cdt := s.app.OracleKeeper.CurrencyDeviationThresholds(s.ctx)
					for i := range currencyDeviationThresholds2 {
						s.Require().Contains(cdt, currencyDeviationThresholds2[i])
					}
				}
			}
		})
	}
}

func (s *IntegrationTestSuite) TestMsgServer_GovRemoveCurrencyPairProviders() {
	govAccAddr := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()
	currentCurrencyPairProviders := types.CurrencyPairProvidersList{
		{
			BaseDenom:  "FOO",
			QuoteDenom: "BAR",
			PairAddress: []types.PairAddressProvider{
				{
					Address:         "address",
					AddressProvider: "provider",
				},
			},
			Providers: []string{
				"provider",
			},
		},
		{
			BaseDenom:  "FOOBAR",
			QuoteDenom: "BAR",
			PairAddress: []types.PairAddressProvider{
				{
					Address:         "address2",
					AddressProvider: "provider2",
				},
			},
			Providers: []string{
				"provider2",
			},
		},
		{
			BaseDenom:  types.CheqdSymbol,
			QuoteDenom: types.USDSymbol,
			Providers: []string{
				"binance",
				"coinbase",
			},
		},
		{
			BaseDenom:  "UNI",
			QuoteDenom: "ETH",
			PairAddress: []types.PairAddressProvider{
				{
					Address:         "address4",
					AddressProvider: "eth-uniswap",
				},
			},
			Providers: []string{
				"bitget",
				"eth-uniswap",
			},
		},
	}
	s.app.OracleKeeper.SetCurrencyPairProviders(s.ctx, currentCurrencyPairProviders)

	testCases := []struct {
		name      string
		req       *types.MsgGovRemoveCurrencyPairProviders
		expectErr bool
		errMsg    string
	}{
		{
			"remove nonexisting currency pair",
			&types.MsgGovRemoveCurrencyPairProviders{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,

				CurrencyPairProviders: types.CurrencyPairProvidersList{
					{
						BaseDenom:  "CURR1",
						QuoteDenom: "CURR2",
					},
				},
			},
			false,
			"",
		},
		{
			"remove nonexisting currency pair with existing base denom and one with existing quote denom",
			&types.MsgGovRemoveCurrencyPairProviders{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,

				CurrencyPairProviders: types.CurrencyPairProvidersList{
					{
						BaseDenom:  "OJO",
						QuoteDenom: "CURR2",
					},
					{
						BaseDenom:  "CURR1",
						QuoteDenom: "USD",
					},
				},
			},
			false,
			"",
		},
		{
			"remove existing currency pair",
			&types.MsgGovRemoveCurrencyPairProviders{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,

				CurrencyPairProviders: types.CurrencyPairProvidersList{
					{
						BaseDenom:  "FOO",
						QuoteDenom: "BAR",
					},
				},
			},
			false,
			"",
		},
		{
			"remove multiple existing currency pairs",
			&types.MsgGovRemoveCurrencyPairProviders{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,

				CurrencyPairProviders: types.CurrencyPairProvidersList{
					{
						BaseDenom:  "FOOBAR",
						QuoteDenom: "BAR",
					},
					{
						BaseDenom:  "UNI",
						QuoteDenom: "ETH",
					},
				},
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := tc.req.ValidateBasic()
			if err == nil {
				_, err = s.msgServer.GovRemoveCurrencyPairProviders(s.ctx, tc.req)
				s.Require().NoError(err)

				err = abci.EndBlocker(s.ctx, s.app.OracleKeeper, s.app.FeeabsKeeper)
				s.Require().NoError(err)
			}

			if tc.expectErr {
				s.Require().ErrorContains(err, tc.errMsg)
			} else {
				s.Require().NoError(err)

				switch tc.name {
				case "remove nonexisting currency pair":
					cpp := s.app.OracleKeeper.CurrencyPairProviders(s.ctx)
					s.Require().Equal(currentCurrencyPairProviders, cpp)

				case "remove nonexisting currency pair with existing quote":
					cpp := s.app.OracleKeeper.CurrencyPairProviders(s.ctx)
					s.Require().Equal(currentCurrencyPairProviders, cpp)

				case "remove existing currency pair":
					cpp := s.app.OracleKeeper.CurrencyPairProviders(s.ctx)
					s.Require().Equal(currentCurrencyPairProviders[1:], cpp)

				case "remove multiple existing currency pairs":
					cpp := s.app.OracleKeeper.CurrencyPairProviders(s.ctx)
					s.Require().Equal(types.CurrencyPairProvidersList{
						{
							BaseDenom:  types.CheqdSymbol,
							QuoteDenom: types.USDSymbol,
							Providers: []string{
								"binance",
								"coinbase",
							},
						},
					}, cpp)
				}
			}
		})
	}
}

func (s *IntegrationTestSuite) TestMsgServer_GovRemoveCurrencyDeviationThresholds() {
	govAccAddr := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()
	currentCurrencyDeviationThresholds := types.CurrencyDeviationThresholdList{
		{
			BaseDenom: "FOO",
			Threshold: "2",
		},
		{
			BaseDenom: "BAR",
			Threshold: "2",
		},
		{
			BaseDenom: types.CheqdSymbol,
			Threshold: "2",
		},
		{
			BaseDenom: "FOOBAR",
			Threshold: "2",
		},
	}

	s.app.OracleKeeper.SetCurrencyDeviationThresholds(s.ctx, currentCurrencyDeviationThresholds)

	testCases := []struct {
		name      string
		req       *types.MsgGovRemoveCurrencyDeviationThresholds
		expectErr bool
		errMsg    string
	}{
		{
			"remove nonexisting currency",
			&types.MsgGovRemoveCurrencyDeviationThresholds{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,

				Currencies: []string{"CURR1"},
			},
			false,
			"",
		},
		{
			"remove multiple nonexisting currencies",
			&types.MsgGovRemoveCurrencyDeviationThresholds{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,

				Currencies: []string{"CURR1", "CURR2"},
			},
			false,
			"",
		},
		{
			"remove existing currency",
			&types.MsgGovRemoveCurrencyDeviationThresholds{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,

				Currencies: []string{"FOO"},
			},
			false,
			"",
		},
		{
			"remove multiple existing currencies",
			&types.MsgGovRemoveCurrencyDeviationThresholds{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Height:      9,

				Currencies: []string{"BAR", "FOOBAR"},
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := tc.req.ValidateBasic()
			if err == nil {
				_, err = s.msgServer.GovRemoveCurrencyDeviationThresholds(s.ctx, tc.req)
				s.Require().NoError(err)

				err = abci.EndBlocker(s.ctx, s.app.OracleKeeper, s.app.FeeabsKeeper)
				s.Require().NoError(err)
			}

			if tc.expectErr {
				s.Require().ErrorContains(err, tc.errMsg)
			} else {
				s.Require().NoError(err)

				switch tc.name {
				case "remove nonexisting currency pair":
					cdt := s.app.OracleKeeper.CurrencyDeviationThresholds(s.ctx)
					s.Require().Equal(currentCurrencyDeviationThresholds, cdt)

				case "remove nonexisting currency pair with existing quote":
					cdt := s.app.OracleKeeper.CurrencyDeviationThresholds(s.ctx)
					s.Require().Equal(currentCurrencyDeviationThresholds, cdt)

				case "remove existing currency pair":
					cdt := s.app.OracleKeeper.CurrencyDeviationThresholds(s.ctx)
					s.Require().Equal(currentCurrencyDeviationThresholds[1:], cdt)

				case "remove multiple existing currency pairs":
					cdt := s.app.OracleKeeper.CurrencyDeviationThresholds(s.ctx)
					s.Require().Equal(types.CurrencyDeviationThresholdList{
						{
							BaseDenom: types.CheqdSymbol,
							Threshold: "2",
						},
					}, cdt)
				}
			}
		})
	}
}

func (s *IntegrationTestSuite) TestMsgServer_CancelUpdateGovParams() {
	govAccAddr := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()

	// No plan exists at height
	_, err := s.msgServer.GovCancelUpdateParamPlan(s.ctx,
		&types.MsgGovCancelUpdateParamPlan{
			Authority: govAccAddr,
			Height:    100,
		},
	)
	s.Require().ErrorContains(err, "No param update plan found at block height 100: invalid request")

	// Schedule plan
	_, err = s.msgServer.GovUpdateParams(s.ctx,
		&types.MsgGovUpdateParams{
			Authority:   govAccAddr,
			Title:       "test",
			Description: "test",
			Plan: types.ParamUpdatePlan{
				Keys:   []string{"VoteThreshold"},
				Height: 100,
				Changes: types.Params{
					VoteThreshold: math.LegacyNewDecWithPrec(40, 2),
				},
			},
		},
	)
	s.Require().NoError(err)

	// Plan exists now
	_, err = s.msgServer.GovCancelUpdateParamPlan(s.ctx,
		&types.MsgGovCancelUpdateParamPlan{
			Authority: govAccAddr,
			Height:    100,
		},
	)
	s.Require().NoError(err)

	plan := s.app.OracleKeeper.GetParamUpdatePlans(s.ctx)
	s.Require().Len(plan, 0)
}
