package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p ParamUpdatePlan) String() string {
	due := p.DueAt()
	return fmt.Sprintf(`Oracle Param Update Plan
  Keys: %s
  %s
  Changes: %s.`, p.Keys, due, p.Changes)
}

// ValidateBasic does basic validation of a ParamUpdatePlan
func (p ParamUpdatePlan) ValidateBasic() error {
	if p.Height <= 0 {
		return ErrInvalidRequest.Wrap("height must be greater than 0")
	}

	for _, key := range p.Keys {
		switch key {
		case string(KeyVotePeriod):
			if err := validateVotePeriod(p.Changes.VotePeriod); err != nil {
				return err
			}

		case string(KeyVoteThreshold):
			if err := validateVoteThreshold(p.Changes.VoteThreshold); err != nil {
				return err
			}

		case string(KeyRewardBands):
			if err := validateRewardBands(p.Changes.RewardBands); err != nil {
				return err
			}

		case string(KeyRewardDistributionWindow):
			if err := validateRewardDistributionWindow(p.Changes.RewardDistributionWindow); err != nil {
				return err
			}

		case string(KeyAcceptList):
			if err := validateDenomList(p.Changes.AcceptList); err != nil {
				return err
			}

		case string(KeyMandatoryList):
			if err := validateDenomList(p.Changes.MandatoryList); err != nil {
				return err
			}

		case string(KeySlashFraction):
			if err := validateSlashFraction(p.Changes.SlashFraction); err != nil {
				return err
			}

		case string(KeySlashWindow):
			if err := validateSlashWindow(p.Changes.SlashWindow); err != nil {
				return err
			}

		case string(KeyMinValidPerWindow):
			if err := validateMinValidPerWindow(p.Changes.MinValidPerWindow); err != nil {
				return err
			}

		case string(KeyHistoricStampPeriod):
			if err := validateHistoricStampPeriod(p.Changes.HistoricStampPeriod); err != nil {
				return err
			}

		case string(KeyMedianStampPeriod):
			if err := validateMedianStampPeriod(p.Changes.MedianStampPeriod); err != nil {
				return err
			}

		case string(KeyMaximumPriceStamps):
			if err := validateMaximumPriceStamps(p.Changes.MaximumPriceStamps); err != nil {
				return err
			}

		case string(KeyMaximumMedianStamps):
			if err := validateMaximumMedianStamps(p.Changes.MaximumMedianStamps); err != nil {
				return err
			}

		case string(KeyCurrencyPairProviders):
			if err := validateCurrencyPairProviders(p.Changes.CurrencyPairProviders); err != nil {
				return err
			}

		case string(KeyCurrencyDeviationThresholds):
			if err := validateCurrencyDeviationThresholds(p.Changes.CurrencyDeviationThresholds); err != nil {
				return err
			}

		case string(KeyUsdcIbcDenom):
			if err := validateString(p.Changes.UsdcIbcDenom); err != nil {
				return err
			}

		default:
			return fmt.Errorf("%s is not an existing oracle param key", key)
		}
	}

	return nil
}

// ShouldExecute returns true if the Plan is ready to execute given the current context
func (p ParamUpdatePlan) ShouldExecute(ctx sdk.Context) bool {
	return p.Height == ctx.BlockHeight()
}

// DueAt is a string representation of when this plan is due to be executed
func (p ParamUpdatePlan) DueAt() string {
	return fmt.Sprintf("height: %d", p.Height)
}
