package algorithms

import (
	"github.com/red-life/http-lb"
	"go.uber.org/zap"
	"sync"
)

var _ http_lb.LoadBalancingAlgorithm = (*RoundRobin)(nil)

func NewRoundRobin(backendPool http_lb.BackendPool, logger *zap.Logger) *RoundRobin {
	return &RoundRobin{
		backendPool: backendPool,
		logger:      logger,
	}
}

type RoundRobin struct {
	counter     int
	backendPool http_lb.BackendPool
	lock        sync.Mutex
	logger      *zap.Logger
}

func (r *RoundRobin) SelectBackend(_ http_lb.Request) (string, error) {
	addrs := r.backendPool.Backends()
	if len(addrs) <= 0 {
		r.logger.Error("no backend available")
		return "", http_lb.ErrNoServerAvailable
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	defer func() { r.counter++ }()
	if r.counter > len(addrs)-1 {
		r.counter = 0
	}
	return addrs[r.counter], nil
}
