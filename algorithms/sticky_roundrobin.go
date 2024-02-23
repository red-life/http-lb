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

func NewStickyRoundRobin(serverPool http_lb.ServerPool, logger *zap.Logger) *StickyRoundRobin {
	return &StickyRoundRobin{
		cache:      cache.New(DefaultExpiration, DefaultCleanupInterval),
		serverPool: serverPool,
		logger:     logger,
	}
}

type StickyRoundRobin struct {
	counter    int
	cache      *cache.Cache
	serverPool http_lb.ServerPool
	lock       sync.Mutex
	logger     *zap.Logger
}

func (s *StickyRoundRobin) SelectServer(r http_lb.Request) (string, error) {
	servers := s.serverPool.HealthyServers()
	if len(servers) <= 0 {
		s.logger.Error("no server is available")
		return "", http_lb.ErrNoHealthyServerAvailable
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	requestIP := r.RemoteIP
	if server, ok := s.cache.Get(requestIP); ok {
		if http_lb.ContainsSlice(servers, server.(string)) {
			s.logger.Debug("server found in cache",
				zap.String("server", server.(string)), zap.String("cacheKey", requestIP))
			return server.(string), nil
		} else {
			s.logger.Debug("invalidate cached server",
				zap.String("server", server.(string)), zap.String("cacheKey", requestIP))
			s.cache.Delete(requestIP)
		}
	}
	if s.counter > len(servers)-1 {
		s.counter = 0
	}
	selectedServer := servers[s.counter]
	s.cache.SetDefault(requestIP, selectedServer)
	s.counter++
	return selectedServer, nil
}
