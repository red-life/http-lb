package algorithms

import (
	"github.com/red-life/http-lb"
	"go.uber.org/zap"
)

var _ http_lb.LoadBalancingAlgorithm = (*IPHash)(nil)

func NewIPHash(hash http_lb.HashingAlgorithm, addrMng http_lb.AddrsManager, logger *zap.Logger) *IPHash {
	return &IPHash{
		hash:    hash,
		addrMng: addrMng,
		logger:  logger,
	}
}

type IPHash struct {
	hash    http_lb.HashingAlgorithm
	addrMng http_lb.AddrsManager
	logger  *zap.Logger
}

func (i *IPHash) ChooseBackend(r http_lb.Request) (string, error) {
	addrs := i.addrMng.GetBackends()
	if len(addrs) <= 0 {
		i.logger.Error("no backend available")
		return "", http_lb.ErrNoServerAvailable
	}
	idx := int(i.hash(r.RemoteIP)) % len(addrs)
	return addrs[idx], nil
}
