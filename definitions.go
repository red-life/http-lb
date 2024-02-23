package http_lb

import (
	"net/http"
	"net/http/httputil"
)

type HashingAlgorithm func(string) uint

type ServerStatus int

func (s ServerStatus) String() string {
	if s == Healthy {
		return "up"
	}
	return "down"
}

const (
	Healthy ServerStatus = iota + 1
	Unhealthy
)

type Server struct {
	Address string
	Status  ServerStatus
}

type ServerPool interface {
	RegisterServer(string) error
	UnregisterServer(string) error
	SetServerStatus(string, ServerStatus) error
	Servers() []Server
	HealthyServers() []string
	UnhealthyServers() []string
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
