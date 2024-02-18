package algorithms

import (
	"github.com/red-life/http-lb"
	"go.uber.org/zap"
	"sync"
)

var _ http_lb.LoadBalancingAlgorithm = (*RoundRobin)(nil)

func NewRoundRobin(addrMng http_lb.AddrsManager, logger *zap.Logger) *RoundRobin {
	return &RoundRobin{
		addrMng: addrMng,
		logger:  logger,
	}
}

type RoundRobin struct {
	counter int
	addrMng http_lb.AddrsManager
	lock    sync.Mutex
	logger  *zap.Logger
}

func (r *RoundRobin) ChooseBackend(_ http_lb.Request) (string, error) {
	addrs := r.addrMng.GetBackends()
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
