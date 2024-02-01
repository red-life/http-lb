package algorithms

import "github.com/red-life/http-lb"

func NewURLHash(backendAddrs []string, hash func(string) uint) *URLHash {
	return &URLHash{
		backendAddrs: backendAddrs,
		hash:         hash,
	}
}

type URLHash struct {
	hash         func(string) uint
	backendAddrs []string
}

func (u *URLHash) ChooseBackend(r http_lb.Request) string {
	return u.backendAddrs[int(u.hash(r.URLPath))%len(u.backendAddrs)]
}
