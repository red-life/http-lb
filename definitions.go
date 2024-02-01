package http_lb

import (
	"time"
)

type LoadBalancingAlgorithm interface {
	ChooseBackend(r Request) string
}

type Request struct {
	RemoteIP string
	URLPath  string
}

type KeepAlive struct {
	MaxIdleConns     int
	IdleConnsTimeout time.Duration
}

type TransportOptions struct {
	Timeout   time.Duration
	KeepAlive *KeepAlive
}
