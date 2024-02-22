package algorithms

import (
	"crypto/rand"
	"github.com/red-life/http-lb"
	"go.uber.org/zap"
	"math/big"
)

var _ http_lb.LoadBalancingAlgorithm = (*Random)(nil)

func NewRandom(serverPool http_lb.ServerPool, logger *zap.Logger) *Random {
	return &Random{
		serverPool: serverPool,
		logger:     logger,
	}
}

type Random struct {
	serverPool http_lb.ServerPool
	logger     *zap.Logger
}

func (r *Random) SelectServer(_ http_lb.Request) (string, error) {
	servers := r.serverPool.Servers()
	if len(servers) <= 0 {
		r.logger.Error("no server is available")
		return "", http_lb.ErrNoServerAvailable
	}
	randomNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(len(servers))))
	idx := int(randomNumber.Int64())
	return servers[idx], nil
}
