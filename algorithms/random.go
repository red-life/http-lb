package algorithms

import (
	"crypto/rand"
	"github.com/red-life/http-lb"
	"math/big"
)

func NewRandom(backendAddrs []string) *Random {
	return &Random{backendAddrs: backendAddrs}
}

type Random struct {
	backendAddrs []string
}

func (r *Random) ChooseBackend(_ http_lb.Request) string {
	randomNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(len(r.backendAddrs))))
	return r.backendAddrs[int(randomNumber.Int64())]
}
