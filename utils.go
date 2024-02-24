package http_lb

import (
	"hash/fnv"
)

func Hash(input string) uint {
	hash := fnv.New32()
	hash.Write([]byte(input))
	return uint(hash.Sum32())
}

func CopySlice[T any](slice []T) []T {
	copySlice := make([]T, len(slice))
	copy(copySlice, slice)
	return copySlice
}

func ContainsSlice[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}
