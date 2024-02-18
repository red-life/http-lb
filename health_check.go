package http_lb

import (
	"errors"
	"fmt"
	"time"
)

var _ HealthChecker = (*HealthCheck)(nil)

func NewHealthCheck(endPoint string, interval time.Duration, timeout time.Duration, addrsMng AddrsManager, expectedStatusCode int) *HealthCheck {
	return &HealthCheck{
		endPoint:           endPoint,
		interval:           interval,
		timeout:            timeout,
		addrsMng:           addrsMng,
		expectedStatusCode: expectedStatusCode,
	}
}

type HealthCheck struct {
	endPoint            string
	interval            time.Duration
	timeout             time.Duration
	addrsMng            AddrsManager
	expectedStatusCode  int
	unavailableBackends []string
}

func (h *HealthCheck) Run() {
	go func() {
		for {
			foundUnavailableBackends := h.findUnavailableBackends()
			unavailableBackends := DifferenceSlices(foundUnavailableBackends, h.unavailableBackends)
			availableBackends := DifferenceSlices(h.unavailableBackends, foundUnavailableBackends)
			h.unavailableBackends = unavailableBackends
			_ = h.register(availableBackends)
			_ = h.unregister(unavailableBackends)
			time.Sleep(h.interval)
		}
	}()
}

func (h *HealthCheck) findUnavailableBackends() []string {
	var unavailableBackends []string
	for _, addr := range h.addrsMng.GetBackends() {
		resp, err := HttpGet(fmt.Sprintf("%s/%s", addr, h.endPoint), h.timeout)
		if err != nil || resp.StatusCode != h.expectedStatusCode {
			unavailableBackends = append(unavailableBackends, addr)
		}
	}
	return unavailableBackends
}

func (h *HealthCheck) unregister(addrs []string) error {
	for _, addr := range addrs {
		err := h.addrsMng.UnregisterBackend(addr)
		if !errors.Is(err, ErrBackendNotExist) {
			return err
		}
	}
	return nil
}

func (h *HealthCheck) register(addrs []string) error {
	for _, addr := range addrs {
		err := h.addrsMng.RegisterBackend(addr)
		if !errors.Is(err, ErrBackendExists) {
			return err
		}
	}
	return nil
}
