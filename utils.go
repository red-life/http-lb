package http_lb

import (
	"hash/fnv"
	"net"
	"net/http"
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
