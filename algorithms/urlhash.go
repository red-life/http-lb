package algorithms

import "github.com/red-life/http-lb"

var _ http_lb.LoadBalancingAlgorithm = (*URLHash)(nil)

func NewURLHash(hash http_lb.HashingAlgorithm, addrMng http_lb.AddrsManager) *URLHash {
	return &URLHash{
		hash:    hash,
		addrMng: addrMng,
	}
}

type URLHash struct {
	hash    http_lb.HashingAlgorithm
	addrMng http_lb.AddrsManager
}

func (u *URLHash) ChooseBackend(r http_lb.Request) string {
	addrs := u.addrMng.GetBackends()
	return addrs[int(u.hash(r.URLPath))%len(addrs)]
}
