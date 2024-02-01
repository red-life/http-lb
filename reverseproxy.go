package http_lb

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(backends []Backend) *CustomReverseProxy {
	reverseProxy := &CustomReverseProxy{}
	for _, b := range backends {
		parsedURL, _ := url.Parse(b.Addr)
		rp := &httputil.ReverseProxy{}
		rp.Rewrite = func(request *httputil.ProxyRequest) {
			RewriteHeaders(b.Headers)(request)
			RewriteXForwarded(request)
			RewriteURL(parsedURL)
		}
		reverseProxy.cache[b.Addr] = rp
	}
	return reverseProxy
}

type CustomReverseProxy struct {
	cache map[string]*httputil.ReverseProxy
}

func (rp *CustomReverseProxy) ServeHTTP(backendAddr string, rw http.ResponseWriter, r *http.Request) {
	rp.cache[backendAddr].ServeHTTP(rw, r)
}
