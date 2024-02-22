package algorithms

import (
	"github.com/red-life/http-lb"
	"go.uber.org/zap"
)

var _ http_lb.LoadBalancingAlgorithm = (*IPHash)(nil)

func NewIPHash(hash http_lb.HashingAlgorithm, serverPool http_lb.ServerPool, logger *zap.Logger) *IPHash {
	return &IPHash{
		hash:       hash,
		serverPool: serverPool,
		logger:     logger,
	}
}

type IPHash struct {
	hash       http_lb.HashingAlgorithm
	serverPool http_lb.ServerPool
	logger     *zap.Logger
}

func (i *IPHash) SelectServer(r http_lb.Request) (string, error) {
	servers := i.serverPool.Servers()
	if len(servers) <= 0 {
		i.logger.Error("no server is available")
		return "", http_lb.ErrNoServerAvailable
	}
	idx := int(i.hash(r.RemoteIP)) % len(servers)
	return servers[idx], nil
}
