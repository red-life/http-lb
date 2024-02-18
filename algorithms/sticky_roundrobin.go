package algorithms

import (
	"github.com/patrickmn/go-cache"
	"github.com/red-life/http-lb"
	"go.uber.org/zap"
	"sync"
	"time"
)

var _ http_lb.LoadBalancingAlgorithm = (*StickyRoundRobin)(nil)

const DefaultExpiration = 30 * time.Minute
const DefaultCleanupInterval = 15 * time.Minute

func NewStickyRoundRobin(addrMng http_lb.AddrsManager, logger *zap.Logger) *StickyRoundRobin {
	return &StickyRoundRobin{
		cache:   cache.New(DefaultExpiration, DefaultCleanupInterval),
		addrMng: addrMng,
		logger:  logger,
	}
}

type StickyRoundRobin struct {
	counter int
	cache   *cache.Cache
	addrMng http_lb.AddrsManager
	lock    sync.Mutex
	logger  *zap.Logger
}

func (s *StickyRoundRobin) ChooseBackend(r http_lb.Request) (string, error) {
	addrs := s.addrMng.GetBackends()
	if len(addrs) <= 0 {
		s.logger.Error("no backend available")
		return "", http_lb.ErrNoServerAvailable
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	requestIP := r.RemoteIP
	if backend, ok := s.cache.Get(requestIP); ok {
		if http_lb.ContainsSlice(addrs, backend.(string)) {
			s.logger.Debug("backend addr found in cache",
				zap.String("addr", backend.(string)), zap.String("cacheKey", requestIP))
			return backend.(string), nil
		} else {
			s.logger.Debug("invalidate cached backend addr",
				zap.String("addr", backend.(string)), zap.String("cacheKey", requestIP))
			s.cache.Delete(requestIP)
		}
	}
	if s.counter > len(addrs)-1 {
		s.counter = 0
	}
	chosenBackend := addrs[s.counter]
	s.cache.SetDefault(requestIP, chosenBackend)
	s.counter++
	return chosenBackend, nil
}
