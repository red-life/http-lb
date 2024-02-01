package algorithms

import "github.com/red-life/http-lb"

func NewIPHash(backendAddrs []string, hash func(string) uint) *IPHash {
	return &IPHash{
		backendAddrs: backendAddrs,
		hash:         hash,
	}
}

type IPHash struct {
	hash         func(string) uint
	backendAddrs []string
}

func (i *IPHash) ChooseBackend(r http_lb.Request) string {
	return i.backendAddrs[int(i.hash(r.RemoteIP))%len(i.backendAddrs)]
}
