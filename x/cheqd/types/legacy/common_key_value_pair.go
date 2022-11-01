package legacy

import (
	"encoding/json"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Helpers

type PublicKeyJWK []*KeyValuePair

func PubKeyJWKToMap(key PublicKeyJWK) map[string]string {
	map_ := make(map[string]string)
	for _, kv := range key {
		map_[kv.Key] = kv.Value
	}
	return map_
}

func JSONToPubKeyJWK(jsonStr string) PublicKeyJWK {
	map_ := make(map[string]string)
	res := PublicKeyJWK{}
	err := json.Unmarshal([]byte(jsonStr), &map_)
	if err != nil {
		panic(fmt.Errorf("internal error: Cannot unmarshal JSON string: %s", jsonStr))
	}
	for k, v := range map_ {
		res = append(res, &KeyValuePair{
			Key:   k,
			Value: v,
		})
	}
	return res
}

func PubKeyJWKToJson(key PublicKeyJWK) (string, error) {
	map_ := PubKeyJWKToMap(key)
	json_, err_ := json.Marshal(map_)
	if err_ != nil {
		return "", errors.New("can't marshal PublicKeyJWK map to JSON")
	}

	return string(json_), nil
}

func IsUniqueKeyValuePairListByKey(key PublicKeyJWK) bool {
	map_ := PubKeyJWKToMap(key)
	return len(map_) == len(key)
}

// Validation

func (p KeyValuePair) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Key, validation.Required),
		validation.Field(&p.Value, validation.Required),
	)
}

func IsUniqueKeyValuePairListByKeyRule() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.([]*KeyValuePair)
		if !ok {
			panic("IsUniqueKeyValuePairListByKeyRule must be only applied on KeyValuePair array properties")
		}

		if !IsUniqueKeyValuePairListByKey(casted) {
			return errors.New("the list of KeyValuePair should be without duplicates")
		}

		return nil
	})
}
