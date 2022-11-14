//go:build upgrade

package upgrade

var (
	ExistingDidDocCreatePayloads, ExistingDidDocUpdatePayloads, ExistingDidDocDeactivatePayloads          []string
	ExistingSignInputCreatePayloads, ExistingSignInputUpdatePayloads, ExistingSignInputDeactivatePayloads []string
	ExistingResourceCreatePayloads                                                                        []string
	ExpectedDidDocCreateRecords, ExpectedDidDocUpdateRecords, ExpectedDidDocDeactivateRecords             []string
	ExpectedResourceCreateRecords                                                                         []string
)

// Pre
var (
	CURRENT_HEIGHT    int64
	VOTING_END_HEIGHT int64
	UPGRADE_HEIGHT    int64
	HEIGHT_ERROR      error
)
