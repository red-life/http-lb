package http_lb

import (
	"net"
	"net/http"
	"time"
)

var _ TransportFactory = (*TransFactory)(nil)

func NewTransportFactory(timeout time.Duration, maxIdleConns int, idleConnsTimeout time.Duration) *TransFactory {
	return &TransFactory{
		timeout:          timeout,
		maxIdleConns:     maxIdleConns,
		idleConnsTimeout: idleConnsTimeout,
	}
}

type TransFactory struct {
	timeout          time.Duration
	maxIdleConns     int
	idleConnsTimeout time.Duration
}

func (t *TransFactory) Create() *http.Transport {
	transport := &http.Transport{}
	dialer := &net.Dialer{
		Timeout: t.timeout,
	}
	transport.MaxIdleConns = t.maxIdleConns
	transport.IdleConnTimeout = t.idleConnsTimeout
	transport.DialContext = dialer.DialContext
	return transport
}
