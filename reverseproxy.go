package http_lb

import (
	"net/http/httputil"
	"net/url"
)

var _ ReverseProxyFactory = (*RPFactory)(nil)

func RewriteXForwarded(p *httputil.ProxyRequest) {
	p.Out.Header.Del("X-Forwarded-For")
	p.SetXForwarded()
}

func RewriteURL(url *url.URL) func(*httputil.ProxyRequest) {
	return func(p *httputil.ProxyRequest) {
		p.SetURL(url)
	}
}

func NewRPFactory(transportFactory TransportFactory, reWriters ...func(request *httputil.ProxyRequest)) *RPFactory {
	reverseProxy := &RPFactory{
		cache:            make(map[string]*httputil.ReverseProxy),
		reWriters:        reWriters,
		transportFactory: transportFactory,
	}
	return reverseProxy
}

type RPFactory struct {
	cache            map[string]*httputil.ReverseProxy
	reWriters        []func(request *httputil.ProxyRequest)
	transportFactory TransportFactory
}

func (rp *RPFactory) Create(backendAddr string) (*httputil.ReverseProxy, error) {
	if proxy, ok := rp.cache[backendAddr]; ok {
		return proxy, nil
	}
	proxy := &httputil.ReverseProxy{}
	if len(rp.reWriters) > 0 {
		proxy.Rewrite = func(request *httputil.ProxyRequest) {
			for _, reWriter := range rp.reWriters {
				reWriter(request)
			}
		}
	}
	transport, err := rp.transportFactory.Create(backendAddr)
	if err != nil {
		return nil, err
	}
	proxy.Transport = transport
	rp.cache[backendAddr] = proxy
	return proxy, nil
}
