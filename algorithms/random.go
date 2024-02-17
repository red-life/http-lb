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

func (r *Random) ChooseBackend(_ http_lb.Request) string {
	addrs := r.addrMng.GetBackends()
	randomNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(len(addrs))))
	return addrs[int(randomNumber.Int64())]
}
