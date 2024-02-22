package http_lb

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var _ HealthChecker = (*HealthCheck)(nil)
var _ GracefulShutdown = (*HealthCheck)(nil)

func NewHealthCheck(endPoint string, interval time.Duration, serverPool ServerPool,
	expectedStatusCode int, client *http.Client, logger *zap.Logger) *HealthCheck {
	return &HealthCheck{
		endPoint:           endPoint,
		interval:           interval,
		serverPool:         serverPool,
		expectedStatusCode: expectedStatusCode,
		client:             client,
		logger:             logger,
		shutdownCh:         make(chan struct{}, 1),
	}
}

type HealthCheck struct {
	endPoint           string
	interval           time.Duration
	client             *http.Client
	serverPool         ServerPool
	expectedStatusCode int
	unavailableServers []string
	logger             *zap.Logger
	shutdownCh         chan struct{}
}

func (h *HealthCheck) Run() {
	go func() {
		ticker := time.NewTicker(h.interval)
		for range ticker.C {
			select {
			case <-h.shutdownCh:
				return
			default:
				h.run()
			}
		}
	}()
}

func (h *HealthCheck) Shutdown() error {
	h.shutdownCh <- struct{}{}
	return nil
}

func (h *HealthCheck) run() {
	foundUnavailableServers := h.findUnavailableServers()
	unavailableServers := DifferenceSlices(foundUnavailableServers, h.unavailableServers)
	_ = h.unregister(unavailableServers)
	if len(h.unavailableServers) > 0 {
		availableServers := DifferenceSlices(h.unavailableServers, foundUnavailableServers)
		_ = h.register(availableServers)
	}
	h.unavailableServers = foundUnavailableServers
}

func (h *HealthCheck) findUnavailableServers() []string {
	var unavailableServers []string
	serversToCheck := append(h.serverPool.Servers(), h.unavailableServers...)
	for _, server := range serversToCheck {
		resp, err := h.client.Get(fmt.Sprintf("%s%s", server, h.endPoint))
		if err == nil && resp.StatusCode == h.expectedStatusCode {
			h.logger.Info("server is up", zap.String("server", server))
			continue
		}
		if err != nil {
			h.logger.Warn("server went down", zap.String("server", server), zap.Error(err))
		} else if resp.StatusCode != h.expectedStatusCode {
			h.logger.Warn("server went down", zap.Int("statusCode", resp.StatusCode),
				zap.String("server", server))
		}
		unavailableServers = append(unavailableServers, server)
		continue

	}
	return unavailableServers
}

func (h *HealthCheck) unregister(servers []string) error {
	for _, server := range servers {
		err := h.serverPool.UnregisterServer(server)
		if !errors.Is(err, ErrServerNotExist) {
			return err
		}
	}
	return nil
}

func (h *HealthCheck) register(servers []string) error {
	for _, server := range servers {
		err := h.serverPool.RegisterServer(server)
		if !errors.Is(err, ErrServerExists) {
			return err
		}
	}
	return nil
}
