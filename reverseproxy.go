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

func NewRPFactory(transportFactory TransportFactory) *RPFactory {
	reverseProxy := &RPFactory{
		cache:            make(map[string]*httputil.ReverseProxy),
		transportFactory: transportFactory,
	}
	return reverseProxy
}

type RPFactory struct {
	cache            map[string]*httputil.ReverseProxy
	transportFactory TransportFactory
}

func (rp *RPFactory) Create(server string) (*httputil.ReverseProxy, error) {
	if proxy, ok := rp.cache[server]; ok {
		return proxy, nil
	}
	proxy := &httputil.ReverseProxy{}
	parsedUrl, _ := url.Parse(server)
	rewriteURL := RewriteURL(parsedUrl)
	proxy.Rewrite = func(request *httputil.ProxyRequest) {
		RewriteXForwarded(request)
		rewriteURL(request)
	}
	proxy.Transport = rp.transportFactory.Create()
	rp.cache[server] = proxy
	return proxy, nil
}
