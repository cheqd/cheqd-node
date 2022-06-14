package utils

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/tendermint/tendermint/types"
)

func GetTxHash(txBytes []byte) string {
	// return base64.StdEncoding.EncodeToString(tmhash.Sum(txBytes))
	return fmt.Sprintf("%X", types.Tx(txBytes).Hash())
}

func ValidateUUID(u string) error {
	_, err := uuid.Parse(u)
	return err
}
