package http_lb

import (
	"hash/fnv"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
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
