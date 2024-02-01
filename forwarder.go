package http_lb

import (
	"net/http"
)

type Forwarder struct {
	lbAlgo       LoadBalancingAlgorithm
	reverseProxy ReverseProxy
}

func (f *Forwarder) Forward(rw http.ResponseWriter, r *http.Request) {
	request := Request{
		RemoteIP: r.RemoteAddr,
		URLPath:  r.URL.Path,
	}
	chosenBackendAddr := f.lbAlgo.ChooseBackend(request)
	f.reverseProxy.ServeHTTP(chosenBackendAddr, rw, r)
}
