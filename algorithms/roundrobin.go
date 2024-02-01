package algorithms

import "sync"

func NewRoundRobin(backendAddrs []string) *RoundRobin {
	return &RoundRobin{backendAddrs: backendAddrs}
}

type RoundRobin struct {
	mutex        sync.Mutex
	counter      int
	backendAddrs []string
}

func (r *RoundRobin) ChooseBackend() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	defer func() { r.counter++ }()
	if r.counter > len(r.backendAddrs)-1 {
		r.counter = 0
	}
	return r.backendAddrs[r.counter]
}
