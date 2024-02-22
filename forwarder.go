package http_lb

import (
	"go.uber.org/zap"
	"net/http"
)

var _ RequestForwarder = (*Forwarder)(nil)

func NewForwarder(lbAlgo LoadBalancingAlgorithm, rpFactory ReverseProxyFactory, logger *zap.Logger) *Forwarder {
	return &Forwarder{
		lbAlgo:    lbAlgo,
		rpFactory: rpFactory,
		logger:    logger,
	}
}

type Forwarder struct {
	lbAlgo    LoadBalancingAlgorithm
	rpFactory ReverseProxyFactory
	logger    *zap.Logger
}

func (f *Forwarder) Forward(rw http.ResponseWriter, r *http.Request) error {
	request := Request{
		RemoteIP: r.RemoteAddr,
		URLPath:  r.URL.Path,
	}
	selectedServer, err := f.lbAlgo.SelectServer(request)
	if err != nil {
		return err
	}
	f.logger.Info("server selected", zap.String("server", selectedServer))
	rp, err := f.rpFactory.Create(selectedServer)
	if err != nil {
		return err
	}
	rp.ServeHTTP(rw, r)
	return nil
}
