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

func (u *URLHash) ChooseBackend(r http_lb.Request) (string, error) {
	addrs := u.addrMng.GetBackends()
	if len(addrs) <= 0 {
		return "", http_lb.ErrNoServerAvailable
	}
	idx := int(u.hash(r.URLPath)) % len(addrs)
	return addrs[idx], nil
}
