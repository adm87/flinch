package manifest

import "hash/fnv"

func HashFNV(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}
