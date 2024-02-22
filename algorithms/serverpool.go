package algorithms

import (
	http_lb "github.com/red-life/http-lb"
	"go.uber.org/zap"
	"sync"
)

var _ http_lb.ServerPool = (*ServerPool)(nil)

func NewServerPool(servers []string, logger *zap.Logger) *ServerPool {
	return &ServerPool{
		servers: http_lb.CopySlice(servers),
		logger:  logger,
	}
}

type ServerPool struct {
	servers []string
	rwLock  sync.RWMutex
	logger  *zap.Logger
}

func (b *ServerPool) RegisterServer(server string) error {
	if b.find(server) != -1 {
		return http_lb.ErrServerExists
	}
	b.rwLock.Lock()
	defer b.rwLock.Unlock()
	b.servers = append(b.servers, server)
	b.logger.Debug("server registered", zap.String("server", server))
	return nil
}

func (b *ServerPool) UnregisterServer(server string) error {
	if i := b.find(server); i != -1 {
		b.rwLock.Lock()
		defer b.rwLock.Unlock()
		b.servers = append(b.servers[:i], b.servers[i+1:]...)
		b.logger.Debug("server unregistered", zap.String("server", server))
		return nil
	}
	return http_lb.ErrServerNotExist
}

func (b *ServerPool) Servers() []string {
	b.rwLock.RLock()
	defer b.rwLock.RUnlock()
	return http_lb.CopySlice(b.servers)
}

func (b *ServerPool) find(server string) int {
	b.rwLock.RLock()
	defer b.rwLock.RUnlock()
	for i, s := range b.servers {
		if s == server {
			return i
		}
	}
	return -1
}
