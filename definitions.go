package http_lb

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type LoadBalancingAlgorithm interface {
	ChooseBackend(r Request) string
}

func RewriteHeaders(headers map[string]string) func(*httputil.ProxyRequest) {
	return func(p *httputil.ProxyRequest) {
		for key, value := range headers {
			p.Out.Header.Set(key, value)
		}
	}
}

func RewriteXForwarded(p *httputil.ProxyRequest) {
	p.Out.Header.Del("X-Forwarded-For")
	p.Out.Header.Del("X-Forwarded-Host")
	p.Out.Header.Del("X-Forwarded-Proto")
	p.SetXForwarded()
}

func RewriteURL(url *url.URL) func(*httputil.ProxyRequest) {
	return func(p *httputil.ProxyRequest) {
		p.SetURL(url)
	}
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
