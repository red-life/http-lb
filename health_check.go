package http_lb

import (
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
				for _, server := range h.serverPool.Servers() {
					go h.check(server)
				}
			}
		}
	}()
}

func (h *HealthCheck) Shutdown() error {
	h.shutdownCh <- struct{}{}
	return nil
}

func (h *HealthCheck) check(server Server) {
	status := Unhealthy
	resp, err := h.client.Get(fmt.Sprintf("%s%s", server.Address, h.endPoint))
	if err == nil && resp.StatusCode == h.expectedStatusCode {
		status = Healthy
	}
	err = h.serverPool.SetServerStatus(server.Address, status)
	if err != nil {
		h.logger.Error("unexpected error for setting the server status", zap.Error(err), zap.String("server", server.Address))
	}
	if server.Status != status {
		h.logger.Warn(fmt.Sprintf("server went %s", status), zap.String("server", server.Address))
		return
	}
	h.logger.Warn(fmt.Sprintf("server is still %s", status), zap.String("server", server.Address))
}
