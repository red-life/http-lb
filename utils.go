package http_lb

import (
	"hash/fnv"
)

func Hash(input string) uint {
	hash := fnv.New32()
	hash.Write([]byte(input))
	return uint(hash.Sum32())
}
