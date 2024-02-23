package http_lb

import (
	"go.uber.org/zap"
	"sync"
)

var _ ServerPool = (*ServerPoolImplementation)(nil)

func NewServerPool(servers []string, logger *zap.Logger) *ServerPoolImplementation {
	serverPool := &ServerPoolImplementation{
		serverAddrsIndex: make(map[string]int),
		logger:           logger,
	}
	for _, server := range servers {
		serverPool.RegisterServer(server)
	}
	return serverPool
}

type ServerPoolImplementation struct {
	servers          []Server
	serverAddrsIndex map[string]int
	rwLock           sync.RWMutex
	logger           *zap.Logger
}

func (b *ServerPoolImplementation) SetServerStatus(server string, status ServerStatus) error {
	if !b.checkExistence(server) {
		return ErrServerNotExists
	}
	b.rwLock.Lock()
	defer b.rwLock.Unlock()
	b.servers[b.serverAddrsIndex[server]].Status = status
	return nil
}

func (b *ServerPoolImplementation) Servers() []Server {
	return CopySlice(b.servers)
}

func (b *ServerPoolImplementation) HealthyServers() []string {
	return b.getByStatus(Healthy)
}

func (b *ServerPoolImplementation) UnhealthyServers() []string {
	return b.getByStatus(Unhealthy)
}

func (b *ServerPoolImplementation) getByStatus(status ServerStatus) []string {
	servers := make([]string, 0)
	for _, server := range b.servers {
		if server.Status == status {
			servers = append(servers, server.Address)
		}
	}
	return servers
}

func (b *ServerPoolImplementation) RegisterServer(server string) error {
	if b.checkExistence(server) {
		return ErrServerExists
	}
	b.rwLock.Lock()
	defer b.rwLock.Unlock()
	b.servers = append(b.servers, Server{
		Address: server,
		Status:  Healthy,
	})
	b.serverAddrsIndex[server] = len(b.servers) - 1
	b.logger.Debug("server registered", zap.String("server", server))
	return nil
}

func (b *ServerPoolImplementation) UnregisterServer(server string) error {
	if !b.checkExistence(server) {
		return ErrServerNotExists
	}
	b.rwLock.Lock()
	defer b.rwLock.Unlock()
	idx := b.serverAddrsIndex[server]
	b.servers = append(b.servers[:idx], b.servers[idx+1:]...)
	delete(b.serverAddrsIndex, server)
	for s := range b.serverAddrsIndex {
		if b.serverAddrsIndex[s] > idx {
			b.serverAddrsIndex[s]--
		}
	}
	return nil
}

func (b *ServerPoolImplementation) checkExistence(server string) bool {
	b.rwLock.RLock()
	defer b.rwLock.RUnlock()
	_, ok := b.serverAddrsIndex[server]
	return ok
}
