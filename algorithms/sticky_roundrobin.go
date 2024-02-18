package algorithms

import (
	"github.com/patrickmn/go-cache"
	"github.com/red-life/http-lb"
	"sync"
	"time"
)

var _ http_lb.LoadBalancingAlgorithm = (*StickyRoundRobin)(nil)

const DefaultExpiration = 30 * time.Minute
const DefaultCleanupInterval = 15 * time.Minute

func NewStickyRoundRobin(addrMng http_lb.AddrsManager) *StickyRoundRobin {
	return &StickyRoundRobin{
		cache:   cache.New(DefaultExpiration, DefaultCleanupInterval),
		addrMng: addrMng,
	}
}

type StickyRoundRobin struct {
	counter int
	cache   *cache.Cache
	addrMng http_lb.AddrsManager
	lock    sync.Mutex
}

func (s *StickyRoundRobin) ChooseBackend(r http_lb.Request) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	requestIP := r.RemoteIP
	if backend, ok := s.cache.Get(requestIP); ok { // the backend might be unregistered but still available in cache
		return backend.(string), nil
	}
	addrs := s.addrMng.GetBackends()
	if len(addrs) <= 0 {
		return "", http_lb.ErrNoServerAvailable
	}
	if s.counter > len(addrs)-1 {
		s.counter = 0
	}
	chosenBackend := addrs[s.counter]
	s.cache.SetDefault(requestIP, chosenBackend)
	s.counter++
	return chosenBackend, nil
}
