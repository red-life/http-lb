package algorithms

import (
	"github.com/patrickmn/go-cache"
	"github.com/red-life/http-lb"
	"sync"
	"time"
)

const DefaultExpiration = 30 * time.Minute
const DefaultCleanupInterval = 15 * time.Minute

func NewStickyRoundRobin(backendAddrs []string) *StickyRoundRobin {
	return &StickyRoundRobin{
		backendAddrs: backendAddrs,
		cache:        cache.New(DefaultExpiration, DefaultCleanupInterval),
	}
}

type StickyRoundRobin struct {
	mutex        sync.Mutex
	counter      int
	backendAddrs []string
	cache        *cache.Cache
}

func (s *StickyRoundRobin) ChooseBackend(r http_lb.Request) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	requestIP := r.RemoteIP
	if backend, ok := s.cache.Get(requestIP); ok {
		return backend.(string)
	}
	if s.counter > len(s.backendAddrs)-1 {
		s.counter = 0
	}
	chosenBackend := s.backendAddrs[s.counter]
	s.cache.SetDefault(requestIP, chosenBackend)
	s.counter++
	return chosenBackend
}
