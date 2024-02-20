package cli

import (
	http_lb "github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"go.uber.org/zap"
)

func Factory(config Config, logger *zap.Logger) (*http_lb.Frontend, *http_lb.HealthCheck) {
	addrsMng := AddrsManagerFactory(config.Backend, logger)
	lbAlgo := LoadBalancingAlgorithmFactory(addrsMng, http_lb.Hash, logger)(config.Algorithm)
	reverseProxy := RevereProxyFactory(config.Backend)
	forwarder := RequestForwarderFactory(lbAlgo, reverseProxy, logger)
	frontend := FrontendFactory(config.Frontend, forwarder, logger)
	healthCheck := HealthCheckFactory(config.HealthCheck, addrsMng, logger)
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

func HealthCheckFactory(healthCheck HealthCheck, addrMng http_lb.AddrsManager, logger *zap.Logger) *http_lb.HealthCheck {
	return http_lb.NewHealthCheck(healthCheck.Endpoint, healthCheck.Interval, healthCheck.Timeout, addrMng, healthCheck.ExpectedStatusCode, logger)
}

func RequestForwarderFactory(lbAlgo http_lb.LoadBalancingAlgorithm,
	reverseProxy http_lb.ReverseProxy, logger *zap.Logger) http_lb.RequestForwarder {
	return http_lb.NewForwarder(lbAlgo, reverseProxy, logger)
}

func RevereProxyFactory(configBackends []Backend) http_lb.ReverseProxy {
	var backends []http_lb.Backend
	for _, b := range configBackends {
		var keepAlive *http_lb.KeepAlive
		if b.KeepAlive != nil {
			keepAlive = &http_lb.KeepAlive{
				MaxIdleConns:     b.KeepAlive.MaxIdle,
				IdleConnsTimeout: b.KeepAlive.IdleTimeout,
			}
		}
		transport := http_lb.CreateTransport(http_lb.TransportOptions{
			Timeout:   b.Timeout,
			KeepAlive: keepAlive,
		})
		backend := http_lb.Backend{
			Addr:      b.Address,
			Transport: transport,
		}
		backends = append(backends, backend)
	}
	return http_lb.NewReverseProxy(backends)
}

func LoadBalancingAlgorithmFactory(addrMng http_lb.AddrsManager,
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

func AddrsManagerFactory(backendsConfig []Backend, logger *zap.Logger) http_lb.AddrsManager {
	var backends []string
	for _, b := range backendsConfig {
		backends = append(backends, b.Address)
	}
	return algorithms.NewBackendAddrsManager(backends, logger)
}
