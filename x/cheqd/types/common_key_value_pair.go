package types

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Helpers

type PublicKeyJWK []*KeyValuePair

func PubKeyJWKToMap(pjwk PublicKeyJWK) map[string]string {
	rmap := make(map[string]string)
	for _, kv := range pjwk {
		rmap[kv.Key] = kv.Value
	}
	return rmap
}

func JSONToPubKeyJWK(jsonStr string) PublicKeyJWK {
	mjson := make(map[string]string)
	newPJWK := PublicKeyJWK{}
	err_ := json.Unmarshal([]byte(jsonStr), &mjson)
	if err_ != nil {
		panic(fmt.Errorf("internal error: Cannot unmarshal JSON string: %s", jsonStr))
	}
	for k, v := range mjson {
		newPJWK = append(newPJWK, &KeyValuePair{
			Key:   k,
			Value: v,
		})
	}
	return newPJWK
}

func IsUniqueKVSet(key PublicKeyJWK) bool {
	map_ := PubKeyJWKToMap(key)
	return len(map_) == len(key)
}

// Validation

func IsUniqueKeyValuePairSet() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.([]*KeyValuePair)
		if !ok {
			panic("IsUniqueKeyValuePairSet must be only applied on KeyValuePair array properties")
		}

		if !IsUniqueKVSet(casted) {
			return errors.New("the list of KeyValuePair should be without duplicates")
		}

		return nil
	})
}
