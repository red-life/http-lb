package algorithms

import (
	"github.com/red-life/http-lb"
	"go.uber.org/zap"
)

var _ http_lb.LoadBalancingAlgorithm = (*URLHash)(nil)

func NewURLHash(hash http_lb.HashingAlgorithm, backendPool http_lb.BackendPool, logger *zap.Logger) *URLHash {
	return &URLHash{
		hash:        hash,
		backendPool: backendPool,
		logger:      logger,
	}
}

type URLHash struct {
	hash        http_lb.HashingAlgorithm
	backendPool http_lb.BackendPool
	logger      *zap.Logger
}

func (u *URLHash) SelectBackend(r http_lb.Request) (string, error) {
	addrs := u.backendPool.Backends()
	if len(addrs) <= 0 {
		u.logger.Error("no backend available")
		return "", http_lb.ErrNoServerAvailable
	}
	idx := int(u.hash(r.URLPath)) % len(addrs)
	return addrs[idx], nil
}
