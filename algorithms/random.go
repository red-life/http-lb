package algorithms

import (
	"crypto/rand"
	"github.com/red-life/http-lb"
	"go.uber.org/zap"
	"math/big"
)

var _ http_lb.LoadBalancingAlgorithm = (*Random)(nil)

func NewRandom(backendPool http_lb.BackendPool, logger *zap.Logger) *Random {
	return &Random{
		backendPool: backendPool,
		logger:      logger,
	}
}

type Random struct {
	backendPool http_lb.BackendPool
	logger      *zap.Logger
}

func (r *Random) ChooseBackend(_ http_lb.Request) (string, error) {
	addrs := r.backendPool.Backends()
	if len(addrs) <= 0 {
		r.logger.Error("no backend available")
		return "", http_lb.ErrNoServerAvailable
	}
	randomNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(len(addrs))))
	idx := int(randomNumber.Int64())
	return addrs[idx], nil
}
