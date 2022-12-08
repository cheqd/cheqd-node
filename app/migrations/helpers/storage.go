package helpers

import (
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ByteStr []byte

func ReadAllKeys(store types.KVStore, prefix []byte) []ByteStr {
	keys := []ByteStr{}

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer CloseIteratorOrPanic(iterator)

	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, ByteStr(iterator.Key()))
	}

	return keys
}

func CloseIteratorOrPanic(iterator sdk.Iterator) {
	err := iterator.Close()
	if err != nil {
		panic(err.Error())
	}
}
