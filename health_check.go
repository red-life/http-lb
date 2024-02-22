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
	endPoint            string
	interval            time.Duration
	client              *http.Client
	serverPool          ServerPool
	expectedStatusCode  int
	unavailableBackends []string
	logger              *zap.Logger
	shutdownCh          chan struct{}
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
	foundUnavailableBackends := h.findUnavailableBackends()
	unavailableBackends := DifferenceSlices(foundUnavailableBackends, h.unavailableBackends)
	_ = h.unregister(unavailableBackends)
	if len(h.unavailableBackends) > 0 {
		availableBackends := DifferenceSlices(h.unavailableBackends, foundUnavailableBackends)
		_ = h.register(availableBackends)
	}
	h.unavailableBackends = foundUnavailableBackends
}

func (h *HealthCheck) findUnavailableBackends() []string {
	var unavailableBackends []string
	addrsToCheck := append(h.serverPool.Servers(), h.unavailableBackends...)
	for _, addr := range addrsToCheck {
		resp, err := h.client.Get(fmt.Sprintf("%s%s", addr, h.endPoint))
		if err == nil && resp.StatusCode == h.expectedStatusCode {
			h.logger.Info("server is up", zap.String("server", addr))
			continue
		}
		if err != nil {
			h.logger.Warn("server went down", zap.String("server", addr), zap.Error(err))
		} else if resp.StatusCode != h.expectedStatusCode {
			h.logger.Warn("backend went down", zap.Int("statusCode", resp.StatusCode),
				zap.String("addr", addr))
		}
		unavailableBackends = append(unavailableBackends, addr)
		continue

	}
	return unavailableBackends
}

func (h *HealthCheck) unregister(addrs []string) error {
	for _, addr := range addrs {
		err := h.serverPool.UnregisterServer(addr)
		if !errors.Is(err, ErrServerNotExist) {
			return err
		}
	}
	return nil
}

func (h *HealthCheck) register(addrs []string) error {
	for _, addr := range addrs {
		err := h.serverPool.RegisterServer(addr)
		if !errors.Is(err, ErrServerExists) {
			return err
		}
	}
	return nil
}
