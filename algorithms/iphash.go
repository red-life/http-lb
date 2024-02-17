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

func (i *IPHash) ChooseBackend(r http_lb.Request) string {
	addrs := i.addrMng.GetBackends()
	return addrs[int(i.hash(r.RemoteIP))%len(addrs)]
}
