package algorithms

import (
	"errors"
	http_lb "github.com/red-life/http-lb"
	"sync"
)

var _ http_lb.AddrsManager = (*BackendAddrsManager)(nil)

var (
	ErrBackendExists   = errors.New("backend already exists")
	ErrBackendNotExist = errors.New("backend doesn't exist")
)

func NewBackendAddrsManager(backendAddrs []string) *BackendAddrsManager {
	return &BackendAddrsManager{
		backendAddrs: http_lb.CopySlice(backendAddrs),
	}
}

type BackendAddrsManager struct {
	backendAddrs []string
	rwLock       sync.RWMutex
}

func (b *BackendAddrsManager) RegisterBackend(backendAddr string) error {
	if b.find(backendAddr) != -1 {
		return ErrBackendExists
	}
	b.rwLock.Lock()
	defer b.rwLock.Unlock()
	b.backendAddrs = append(b.backendAddrs, backendAddr)
	return nil
}

func (b *BackendAddrsManager) UnregisterBackend(backendAddr string) error {
	if i := b.find(backendAddr); i != -1 {
		b.rwLock.Lock()
		defer b.rwLock.Unlock()
		b.backendAddrs = append(b.backendAddrs[:i], b.backendAddrs[i+1:]...)
		return nil
	}
	return ErrBackendNotExist
}

func (b *BackendAddrsManager) GetBackends() []string {
	b.rwLock.RLock()
	defer b.rwLock.RUnlock()
	return http_lb.CopySlice(b.backendAddrs)
}

func (b *BackendAddrsManager) find(backendAddr string) int {
	b.rwLock.RLock()
	defer b.rwLock.RUnlock()
	for i, backend := range b.backendAddrs {
		if backend == backendAddr {
			return i
		}
	}
	return -1
}
