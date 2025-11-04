package helpers

const (
	// OracleJitterTolerance is the tolerance for oracle price fluctuations
	OracleJitterTolerance = int64(5e3) // 5000 ncheq
	// BalanceJitterTolerance is the tolerance for account balance fluctuations; this is higher than OracleJitterTolerance to account for multiple txs or ICQ twap updates (generally higher slippage - ICQ is less predictable and takes precedence, if moving averages not yet computed)
	BalanceJitterTolerance = int64(2e8) // 200_000_000 ncheq
)
