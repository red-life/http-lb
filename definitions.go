package http_lb

import (
	"net/http"
	"time"
)

type LoadBalancingAlgorithm interface {
	ChooseBackend(Request) string
}

type ReverseProxy interface {
	ServeHTTP(string, http.ResponseWriter, *http.Request)
}

type RequestForwarder interface {
	Forward(http.ResponseWriter, *http.Request)
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

type Backend struct {
	Addr      string
	Transport *http.Transport
	Headers   map[string]string
}
