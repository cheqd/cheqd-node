package utils

import (
	"fmt"

	"github.com/tendermint/tendermint/types"
)

func GetTxHash(txBytes []byte) string {
	// return base64.StdEncoding.EncodeToString(tmhash.Sum(txBytes))
	return fmt.Sprintf("%X", types.Tx(txBytes).Hash())
}
