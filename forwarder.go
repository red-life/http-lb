package http_lb

import (
	"go.uber.org/zap"
	"net/http"
)

var _ RequestForwarder = (*Forwarder)(nil)

func NewForwarder(lbAlgo LoadBalancingAlgorithm, reverseProxy ReverseProxy, logger *zap.Logger) *Forwarder {
	return &Forwarder{
		lbAlgo:       lbAlgo,
		reverseProxy: reverseProxy,
		logger:       logger,
	}
}

type Forwarder struct {
	lbAlgo       LoadBalancingAlgorithm
	reverseProxy ReverseProxy
	logger       *zap.Logger
}

func (f *Forwarder) Forward(rw http.ResponseWriter, r *http.Request) error {
	request := Request{
		RemoteIP: r.RemoteAddr,
		URLPath:  r.URL.Path,
	}
	chosenBackendAddr, err := f.lbAlgo.SelectBackend(request)
	if err != nil {
		return err
	}
	f.logger.Info("backend chose", zap.String("addr", chosenBackendAddr))
	f.reverseProxy.ServeHTTP(chosenBackendAddr, rw, r)
	return nil
}
