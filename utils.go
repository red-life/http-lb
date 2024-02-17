package http_lb

import (
	"hash/fnv"
	"net"
	"net/http"
	"time"
)

func Hash(input string) uint {
	hash := fnv.New32()
	hash.Write([]byte(input))
	return uint(hash.Sum32())
}

func CreateTransport(options TransportOptions) *http.Transport {
	tr := &http.Transport{}
	dialer := &net.Dialer{
		Timeout: options.Timeout,
	}
	if options.KeepAlive != nil {
		tr.MaxIdleConns = options.KeepAlive.MaxIdleConns
		tr.IdleConnTimeout = options.KeepAlive.IdleConnsTimeout
	} else {
		tr.DisableKeepAlives = true
		dialer.KeepAlive = -1
	}
	tr.DialContext = dialer.DialContext
	return tr
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

func HttpGet(url string, timeout time.Duration) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	return client.Get(url)
}
