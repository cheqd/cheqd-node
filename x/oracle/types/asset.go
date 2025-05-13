package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CheqdDenom       string = "ncheq"
	CheqdSymbol      string = "CHEQ"
	CheqdExponent           = uint32(9)
	AtomDenom        string = "ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9"
	AtomSymbol       string = "ATOM"
	AtomExponent            = uint32(6)
	BitcoinDenom     string = "ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9"
	BitcoinSymbol    string = "BTC"
	BitcoinExponent         = uint32(18)
	EthereumDenom    string = "ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9"
	EthereumSymbol   string = "ETH"
	EthereumExponent        = uint32(18)
	USDTDenom        string = "ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9"
	USDTSymbol       string = "USDT"
	USDTExponent            = uint32(18)
	USDSymbol        string = "USD"
	BlocksPerMinute         = uint64(10)
	BlocksPerHour           = BlocksPerMinute * 60
	BlocksPerDay            = BlocksPerHour * 24
	BlocksPerWeek           = BlocksPerDay * 7
	BlocksPerMonth          = BlocksPerDay * 30
	BlocksPerYear           = BlocksPerDay * 365
	MicroUnit               = int64(1e6)
)

type (
	// ExchangeRatePrevote defines a structure to store a validator's prevote on
	// the rate of USD in the denom asset.
	ExchangeRatePrevote struct {
		Hash        VoteHash       `json:"hash"`         // Vote hex hash to protect centralize data source problem
		Denom       string         `json:"denom"`        // Ticker symbol of denomination exchanged against USD
		Voter       sdk.ValAddress `json:"voter"`        // Voter validator address
		SubmitBlock int64          `json:"submit_block"` // Block height at submission
	}

	// ExchangeRateVote defines a structure to store a validator's vote on the
	// rate of USD in the denom asset.
	ExchangeRateVote struct {
		ExchangeRate math.LegacyDec `json:"exchange_rate"` // Exchange rate of a denomination against USD
		Denom        string         `json:"denom"`         // Ticker symbol of denomination exchanged against USD
		Voter        sdk.ValAddress `json:"voter"`         // Voter validator address
	}

	// VoteHash defines a hash value to hide vote exchange rate which is formatted
	// as a HEX string:
	// SHA256("{salt}:{symbol}:{exchangeRate},...,{symbol}:{exchangeRate}:{voter}")
	VoteHash []byte
)
