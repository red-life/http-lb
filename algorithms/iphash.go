package algorithms

import (
	"github.com/red-life/http-lb"
)

var _ http_lb.LoadBalancingAlgorithm = (*IPHash)(nil)

func NewIPHash(hash http_lb.HashingAlgorithm, addrMng http_lb.AddrsManager) *IPHash {
	return &IPHash{
		hash:    hash,
		addrMng: addrMng,
	}
}

type IPHash struct {
	hash    http_lb.HashingAlgorithm
	addrMng http_lb.AddrsManager
}

func (i *IPHash) ChooseBackend(r http_lb.Request) (string, error) {
	addrs := i.addrMng.GetBackends()
	if len(addrs) <= 0 {
		return "", http_lb.ErrNoServerAvailable
	}
	idx := int(i.hash(r.RemoteIP)) % len(addrs)
	return addrs[idx], nil
}
