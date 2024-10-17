package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"google.golang.org/protobuf/proto"
)

func IndexOf(array []string, searchElement string, fromIndex int) int {
	for i, v := range array[fromIndex:] {
		if v == searchElement {
			return fromIndex + i
		}
	}

	return -1
}

func Contains(vs []string, t string) bool {
	return IndexOf(vs, t, 0) >= 0
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Subtract(minuend []string, subtrahend []string) []string {
	m := map[string]bool{}

	for _, v := range minuend {
		m[v] = true
	}

	for _, v := range subtrahend {
		delete(m, v)
	}

	result := make([]string, 0, len(m))

	for k := range m {
		result = append(result, k)
	}

	return result
}

// Unique returns a copy of the passed array with duplicates removed
func Unique(array []string) []string {
	m := map[string]bool{}

	for _, v := range array {
		m[v] = true
	}

	result := make([]string, 0, len(m))

	for k := range m {
		result = append(result, k)
	}

	return result
}

func IsUnique(list []string) bool {
	set := map[string]bool{}

	for _, did := range list {
		set[did] = true
	}

	return len(list) == len(set)
}

func ToInterfaces(list []string) []interface{} {
	res := make([]interface{}, len(list))

	for i := range list {
		res[i] = list[i]
	}

	return res
}

func ReplaceInSlice(list []string, old, new string) {
	for i := range list {
		if list[i] == old {
			list[i] = new
		}
	}
}

func UniqueSorted(ls []string) []string {
	tmp := Unique(ls)
	sort.Strings(tmp)
	return tmp
}

func StrBytes(p string) []byte {
	return []byte(p)
}

// Generic function to validate protobuf-supported fields in a JSON string
func ValidateProtobufFields(jsonString string) error {
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &input); err != nil {
		return errors.New("input should be a valid JSON string")
	}

	for key, value := range input {
		switch value.(type) {
		case string, int, int32, int64, float32, float64, bool, proto.Message:
			continue
		default:
			return fmt.Errorf("field %s is not protobuf-supported", key)
		}
	}

	return nil
}
