package algorithms

import (
	"github.com/red-life/http-lb"
	"go.uber.org/zap"
	"sync"
)

var _ http_lb.LoadBalancingAlgorithm = (*RoundRobin)(nil)

func NewRoundRobin(serverPool http_lb.ServerPool, logger *zap.Logger) *RoundRobin {
	return &RoundRobin{
		serverPool: serverPool,
		logger:     logger,
	}
}

type RoundRobin struct {
	counter    int
	serverPool http_lb.ServerPool
	lock       sync.Mutex
	logger     *zap.Logger
}

func (r *RoundRobin) SelectServer(_ http_lb.Request) (string, error) {
	servers := r.serverPool.HealthyServers()
	if len(servers) <= 0 {
		r.logger.Error("no server is available")
		return "", http_lb.ErrNoHealthyServerAvailable
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	defer func() { r.counter++ }()
	if r.counter > len(servers)-1 {
		r.counter = 0
	}
	return servers[r.counter], nil
}
