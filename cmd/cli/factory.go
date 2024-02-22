package cli

import (
	http_lb "github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"go.uber.org/zap"
)

func Factory(config Config, logger *zap.Logger) (*http_lb.Frontend, *http_lb.HealthCheck) {
	serverPool := ServerPoolFactory(config.Backend, logger)
	lbAlgo := LoadBalancingAlgorithmFactory(serverPool, http_lb.Hash, logger)(config.Algorithm)
	transport := http_lb.NewTransportFactory(config.Backend.Timeout, config.Backend.KeepAlive.MaxIdle, config.Backend.KeepAlive.IdleTimeout)
	reverseProxy := http_lb.NewRPFactory(transport)
	forwarder := RequestForwarderFactory(lbAlgo, reverseProxy, logger)
	frontend := FrontendFactory(config.Frontend, forwarder, logger)
	healthCheck := HealthCheckFactory(config.HealthCheck, serverPool, logger)
	return frontend, healthCheck
}

func FrontendFactory(frontend Frontend, reqForwarder http_lb.RequestForwarder, logger *zap.Logger) *http_lb.Frontend {
	var tls *http_lb.TLSOptions
	if frontend.TLS != nil {
		tls = &http_lb.TLSOptions{
			CertFile: frontend.TLS.CertFile,
			KeyFile:  frontend.TLS.KeyFile,
		}
	}
	return http_lb.NewFrontend(frontend.Listen, tls, reqForwarder, logger)
}

func HealthCheckFactory(healthCheck HealthCheck, serverPool http_lb.BackendPool, logger *zap.Logger) *http_lb.HealthCheck {
	return http_lb.NewHealthCheck(healthCheck.Endpoint, healthCheck.Interval, healthCheck.Timeout, serverPool, healthCheck.ExpectedStatusCode, logger)
}

func RequestForwarderFactory(lbAlgo http_lb.LoadBalancingAlgorithm, rpFactory http_lb.ReverseProxyFactory, logger *zap.Logger) http_lb.RequestForwarder {
	return http_lb.NewForwarder(lbAlgo, rpFactory, logger)
}

func LoadBalancingAlgorithmFactory(addrMng http_lb.BackendPool,
	hash http_lb.HashingAlgorithm, logger *zap.Logger) func(algorithmName string) http_lb.LoadBalancingAlgorithm {
	return func(algorithmName string) http_lb.LoadBalancingAlgorithm {
		if algorithmName == "round-robin" {
			return algorithms.NewRoundRobin(addrMng, logger)
		} else if algorithmName == "sticky-round-robin" {
			return algorithms.NewStickyRoundRobin(addrMng, logger)
		} else if algorithmName == "url-hash" {
			return algorithms.NewURLHash(hash, addrMng, logger)
		} else if algorithmName == "ip-hash" {
			return algorithms.NewURLHash(hash, addrMng, logger)
		} else if algorithmName == "random" {
			return algorithms.NewRandom(addrMng, logger)
		}
		logger.Panic("unknown load balancing algorithm", zap.String("algorithmName", algorithmName))
		return nil
	}
}

func ServerPoolFactory(backendsConfig Backend, logger *zap.Logger) http_lb.BackendPool {
	var backends []string
	for _, addr := range backendsConfig.Addresses {
		backends = append(backends, addr)
	}
	return algorithms.NewBackendPool(backends, logger)
}
