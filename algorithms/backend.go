package algorithms

import (
	http_lb "github.com/red-life/http-lb"
	"go.uber.org/zap"
	"sync"
)

var _ http_lb.AddrsManager = (*BackendAddrsManager)(nil)

func NewBackendAddrsManager(backendAddrs []string, logger *zap.Logger) *BackendAddrsManager {
	return &BackendAddrsManager{
		backendAddrs: http_lb.CopySlice(backendAddrs),
		logger:       logger,
	}
}

type BackendAddrsManager struct {
	backendAddrs []string
	rwLock       sync.RWMutex
	logger       *zap.Logger
}

func (b *BackendAddrsManager) RegisterBackend(backendAddr string) error {
	if b.find(backendAddr) != -1 {
		return http_lb.ErrBackendExists
	}
	b.rwLock.Lock()
	defer b.rwLock.Unlock()
	b.backendAddrs = append(b.backendAddrs, backendAddr)
	b.logger.Debug("backend registered", zap.String("addr", backendAddr))
	return nil
}

func (b *BackendAddrsManager) UnregisterBackend(backendAddr string) error {
	if i := b.find(backendAddr); i != -1 {
		b.rwLock.Lock()
		defer b.rwLock.Unlock()
		b.backendAddrs = append(b.backendAddrs[:i], b.backendAddrs[i+1:]...)
		b.logger.Debug("backend unregistered", zap.String("addr", backendAddr))
		return nil
	}
	return http_lb.ErrBackendNotExist
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
