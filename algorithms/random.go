package algorithms

import (
	"crypto/rand"
	"github.com/red-life/http-lb"
	"math/big"
)

var _ http_lb.LoadBalancingAlgorithm = (*Random)(nil)

func NewRandom(addrMng http_lb.AddrsManager) *Random {
	return &Random{
		addrMng: addrMng,
	}
}

type Random struct {
	addrMng http_lb.AddrsManager
}

func (r *Random) ChooseBackend(_ http_lb.Request) (string, error) {
	addrs := r.addrMng.GetBackends()
	if len(addrs) <= 0 {
		return "", http_lb.ErrNoServerAvailable
	}
	randomNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(len(addrs))))
	idx := int(randomNumber.Int64())
	return addrs[idx], nil
}
