package http_lb

import (
	"go.uber.org/zap"
	"sync"
)

var _ ServerPool = (*ServerPoolImplementation)(nil)

func NewServerPool(servers []string, logger *zap.Logger) *ServerPoolImplementation {
	s := make(map[string]ServerStatus)
	for _, server := range servers {
		s[server] = Healthy
	}
	return &ServerPoolImplementation{
		servers: s,
		logger:  logger,
	}
}

type ServerPoolImplementation struct {
	servers map[string]ServerStatus
	rwLock  sync.RWMutex
	logger  *zap.Logger
}

func (b *ServerPoolImplementation) SetServerStatus(server string, status ServerStatus) error {
	if _, ok := b.servers[server]; !ok {
		return ErrServerNotExist
	}
	b.servers[server] = status
	return nil
}

func (b *ServerPoolImplementation) Servers() []string {
	return KeysMap(b.servers)
}

func (b *ServerPoolImplementation) HealthyServers() []string {
	healthyServers := make([]string, len(b.servers))
	for server, status := range b.servers {
		if status == Healthy {
			healthyServers = append(healthyServers, server)
		}
	}
	return healthyServers
}

func (b *ServerPoolImplementation) UnhealthyServers() []string {
	healthyServers := make([]string, len(b.servers))
	for server, status := range b.servers {
		if status == Unhealthy {
			healthyServers = append(healthyServers, server)
		}
	}
	return healthyServers
}

func (b *ServerPoolImplementation) RegisterServer(server string) error {
	b.rwLock.RLock()
	if _, ok := b.servers[server]; ok {
		return ErrServerExists
	}
	b.rwLock.RUnlock()
	b.rwLock.Lock()
	defer b.rwLock.Unlock()
	b.servers[server] = Healthy
	b.logger.Debug("server registered", zap.String("server", server))
	return nil
}

func (b *ServerPoolImplementation) UnregisterServer(server string) error {
	b.rwLock.Lock()
	defer b.rwLock.Unlock()
	if _, ok := b.servers[server]; !ok {
		return ErrServerNotExist
	}
	delete(b.servers, server)
	return nil
}
