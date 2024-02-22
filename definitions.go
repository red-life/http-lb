package http_lb

import (
	"net/http"
	"net/http/httputil"
)

type HashingAlgorithm func(string) uint

type ServerPool interface {
	RegisterServer(string) error
	UnregisterServer(string) error
	Servers() []string
}

type LoadBalancingAlgorithm interface {
	SelectServer(Request) (string, error)
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

type TLSOptions struct {
	CertFile string
	KeyFile  string
}
