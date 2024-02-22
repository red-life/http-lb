package algorithms

import (
	"github.com/red-life/http-lb"
	"go.uber.org/zap"
)

var _ http_lb.LoadBalancingAlgorithm = (*URLHash)(nil)

func NewURLHash(hash http_lb.HashingAlgorithm, serverPool http_lb.ServerPool, logger *zap.Logger) *URLHash {
	return &URLHash{
		hash:       hash,
		serverPool: serverPool,
		logger:     logger,
	}
}

type URLHash struct {
	hash       http_lb.HashingAlgorithm
	serverPool http_lb.ServerPool
	logger     *zap.Logger
}

func (u *URLHash) SelectServer(r http_lb.Request) (string, error) {
	servers := u.serverPool.Servers()
	if len(servers) <= 0 {
		u.logger.Error("no server available")
		return "", http_lb.ErrNoServerAvailable
	}
	idx := int(u.hash(r.URLPath)) % len(servers)
	return servers[idx], nil
}
