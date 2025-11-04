package types

import (
	"fmt"

	"cosmossdk.io/errors"
	"github.com/cometbft/cometbft/crypto/tmhash"
)

// Oracle sentinel errors
var (
	ErrInvalidExchangeRate   = errors.Register(ModuleName, 2, "invalid exchange rate")
	ErrNoPrevote             = errors.Register(ModuleName, 3, "no prevote")
	ErrNoVote                = errors.Register(ModuleName, 4, "no vote")
	ErrNoVotingPermission    = errors.Register(ModuleName, 5, "unauthorized voter")
	ErrInvalidHash           = errors.Register(ModuleName, 6, "invalid hash")
	ErrInvalidHashLength     = errors.Register(ModuleName, 7, fmt.Sprintf("invalid hash length; should equal %d", tmhash.TruncatedSize)) //nolint: lll
	ErrVerificationFailed    = errors.Register(ModuleName, 8, "hash verification failed")
	ErrRevealPeriodMissMatch = errors.Register(ModuleName, 9, "reveal period of submitted vote does not match with registered prevote") //nolint: lll
	ErrInvalidSaltLength     = errors.Register(ModuleName, 10, "invalid salt length; must be 64")
	ErrInvalidSaltFormat     = errors.Register(ModuleName, 11, "invalid salt format")
	ErrNoAggregatePrevote    = errors.Register(ModuleName, 12, "no aggregate prevote")
	ErrNoAggregateVote       = errors.Register(ModuleName, 13, "no aggregate vote")
	ErrUnknownDenom          = errors.Register(ModuleName, 14, "unknown denom")
	ErrNegativeOrZeroRate    = errors.Register(ModuleName, 15, "invalid exchange rate; should be positive")
	ErrExistingPrevote       = errors.Register(ModuleName, 16, "prevote already submitted for this voting period")
	ErrBallotNotSorted       = errors.Register(ModuleName, 17, "ballot must be sorted before this operation")
	ErrInvalidOraclePrice    = errors.Register(ModuleName, 18, "invalid or unavailable oracle price")
	ErrNoHistoricPrice       = errors.Register(ModuleName, 19, "no historic price for this denom at this block")
	ErrNoMedian              = errors.Register(ModuleName, 20, "no median for this denom at this block")
	ErrNoMedianDeviation     = errors.Register(ModuleName, 21, "no median deviation for this denom at this block")
	ErrNoRewardBand          = errors.Register(ModuleName, 22, "unable to find the reward band the given asset")
	ErrNoValidatorRewardSet  = errors.Register(ModuleName, 23, "unable to find the latest validator reward set")
	ErrNoGovAuthority        = errors.Register(ModuleName, 24, "invalid gov authority to perform these changes")
	ErrInvalidRequest        = errors.Register(ModuleName, 25, "invalid request")
	ErrInvalidParamValue     = errors.Register(ModuleName, 26, "invalid oracle param value")
	ErrEncodeInjVoteExt      = errors.Register(ModuleName, 27, "failed to encode injected vote extension tx")
	ErrNonEqualInjVotesLen   = errors.Register(ModuleName, 28, "number of exchange rate votes in vote extension and extended commit info are not equal") //nolint: lll
	ErrNonEqualInjVotesRates = errors.Register(ModuleName, 29, "injected exchange rate votes and generated exchange votes are not equal")                //nolint: lll
	ErrNoCommitInfo          = errors.Register(ModuleName, 30, "no commit info in process proposal request")
	ErrInvalidWmaStrategy    = errors.Register(ModuleName, 31, "invalid WMA strategy")
)
