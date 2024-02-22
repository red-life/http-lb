package cli

import (
	http_lb "github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"go.uber.org/zap"
	"net/http"
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

func HealthCheckFactory(healthCheck HealthCheck, serverPool http_lb.ServerPool, logger *zap.Logger) *http_lb.HealthCheck {
	client := &http.Client{Timeout: healthCheck.Timeout}
	return http_lb.NewHealthCheck(healthCheck.Endpoint, healthCheck.Interval, serverPool, healthCheck.ExpectedStatusCode, client, logger)
}

func RequestForwarderFactory(lbAlgo http_lb.LoadBalancingAlgorithm, rpFactory http_lb.ReverseProxyFactory, logger *zap.Logger) http_lb.RequestForwarder {
	return http_lb.NewForwarder(lbAlgo, rpFactory, logger)
}

func LoadBalancingAlgorithmFactory(serverPool http_lb.ServerPool,
	hash http_lb.HashingAlgorithm, logger *zap.Logger) func(algorithmName string) http_lb.LoadBalancingAlgorithm {
	return func(algorithmName string) http_lb.LoadBalancingAlgorithm {
		if algorithmName == "round-robin" {
			return algorithms.NewRoundRobin(serverPool, logger)
		} else if algorithmName == "sticky-round-robin" {
			return algorithms.NewStickyRoundRobin(serverPool, logger)
		} else if algorithmName == "url-hash" {
			return algorithms.NewURLHash(hash, serverPool, logger)
		} else if algorithmName == "ip-hash" {
			return algorithms.NewURLHash(hash, serverPool, logger)
		} else if algorithmName == "random" {
			return algorithms.NewRandom(serverPool, logger)
		}
		logger.Panic("unknown load balancing algorithm", zap.String("algorithmName", algorithmName))
		return nil
	}
}

func ServerPoolFactory(backendsConfig Backend, logger *zap.Logger) http_lb.ServerPool {
	var servers []string
	for _, server := range backendsConfig.Servers {
		servers = append(servers, server)
	}
	return algorithms.NewServerPool(servers, logger)
}
