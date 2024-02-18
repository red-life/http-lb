package http_lb

import (
	"net/http"
)

var _ RequestForwarder = (*Forwarder)(nil)

type Forwarder struct {
	lbAlgo       LoadBalancingAlgorithm
	reverseProxy ReverseProxy
}

func (f *Forwarder) Forward(rw http.ResponseWriter, r *http.Request) error {
	request := Request{
		RemoteIP: r.RemoteAddr,
		URLPath:  r.URL.Path,
	}
	chosenBackendAddr, err := f.lbAlgo.ChooseBackend(request)
	if err != nil {
		return err
	}
	f.reverseProxy.ServeHTTP(chosenBackendAddr, rw, r)
	return nil
}
