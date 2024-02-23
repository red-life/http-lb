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

func DifferenceSlices[T comparable](a, b []T) []T {
	mb := make(map[T]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []T
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func ContainsSlice[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func FindSlice[T comparable](slice []T, target T) int {
	for i, item := range slice {
		if item == target {
			return i
		}
	}
	return -1
}

func KeysMap[K comparable, V any](m map[K]V) []K {
	keys := make([]K, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
