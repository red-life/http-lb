package http_lb

import (
	"go.uber.org/zap"
	"sync"
)

var _ ServerPool = (*ServerPoolImplementation)(nil)

func NewServerPool(servers []string, logger *zap.Logger) *ServerPoolImplementation {
	return &ServerPoolImplementation{
		servers: CopySlice(servers),
		logger:  logger,
	}
}

type ServerPoolImplementation struct {
	servers []string
	rwLock  sync.RWMutex
	logger  *zap.Logger
}

func (b *ServerPoolImplementation) RegisterServer(server string) error {
	b.rwLock.RLock()
	if ContainsSlice(b.servers, server) {
		return ErrServerExists
	}
	b.rwLock.RUnlock()
	b.rwLock.Lock()
	defer b.rwLock.Unlock()
	b.servers = append(b.servers, server)
	b.logger.Debug("server registered", zap.String("server", server))
	return nil
}

func (b *ServerPoolImplementation) UnregisterServer(server string) error {
	if i := FindSlice(b.servers, server); i != -1 {
		b.rwLock.Lock()
		defer b.rwLock.Unlock()
		b.servers = append(b.servers[:i], b.servers[i+1:]...)
		b.logger.Debug("server unregistered", zap.String("server", server))
		return nil
	}
	return ErrServerNotExist
}

func (b *ServerPoolImplementation) Servers() []string {
	b.rwLock.RLock()
	defer b.rwLock.RUnlock()
	return CopySlice(b.servers)
}
