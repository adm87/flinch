package encoding

import (
	"hash/fnv"
	"reflect"
)

func HashType[T any]() uint64 {
	var zero T
	fullName := reflect.TypeOf(zero).String()
	h := fnv.New64a()
	h.Write([]byte(fullName))
	return h.Sum64()
}
