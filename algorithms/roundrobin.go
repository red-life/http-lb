package algorithms

import (
	"github.com/red-life/http-lb"
	"sync"
)

var _ http_lb.LoadBalancingAlgorithm = (*RoundRobin)(nil)

func NewRoundRobin(addrMng http_lb.AddrsManager) *RoundRobin {
	return &RoundRobin{
		addrMng: addrMng,
	}
}

type RoundRobin struct {
	counter int
	addrMng http_lb.AddrsManager
	lock    sync.Mutex
}

func (r *RoundRobin) ChooseBackend(_ http_lb.Request) string {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer func() { r.counter++ }()
	addrs := r.addrMng.GetBackends()
	if r.counter > len(addrs)-1 {
		r.counter = 0
	}
	return addrs[r.counter]
}
