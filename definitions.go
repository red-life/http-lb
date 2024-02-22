package http_lb

import (
	"net/http"
	"net/http/httputil"
	"time"
)

type HashingAlgorithm func(string) uint

type BackendPool interface {
	RegisterBackend(string) error
	UnregisterBackend(string) error
	Backends() []string
}

type LoadBalancingAlgorithm interface {
	SelectBackend(Request) (string, error)
}

type RequestForwarder interface {
	Forward(http.ResponseWriter, *http.Request) error
}

type HealthChecker interface {
	Run()
}

type GracefulShutdown interface {
	Shutdown() error
}

type ReverseProxyFactory interface {
	Create(string) (*httputil.ReverseProxy, error)
}

type TransportFactory interface {
	Create() *http.Transport
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
}

type TLSOptions struct {
	CertFile string
	KeyFile  string
}
