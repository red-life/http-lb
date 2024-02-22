package algorithms_test

import (
	"github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"go.uber.org/zap"
	"testing"
)

func TestURLHash_SelectServer(t *testing.T) {
	servers := []string{
		"server 1",
		"server 2",
		"server 3",
		"server 4",
		"server 5",
		"server 6",
	}
	logger, _ := zap.NewDevelopment()
	serverPool := http_lb.NewServerPool(servers, logger)
	urlHash := algorithms.NewURLHash(http_lb.Hash, serverPool, logger)
	tests := []struct {
		input    http_lb.Request
		expected string
	}{
		{http_lb.Request{URLPath: "/"}, servers[int(http_lb.Hash("/"))%len(servers)]},
		{http_lb.Request{URLPath: "/home"}, servers[int(http_lb.Hash("/home"))%len(servers)]},
		{http_lb.Request{URLPath: "/auth/login"}, servers[int(http_lb.Hash("/auth/login"))%len(servers)]},
		{http_lb.Request{URLPath: "/api/v1"}, servers[int(http_lb.Hash("/api/v1"))%len(servers)]},
	}
	for i, test := range tests {
		selectedServer, err := urlHash.SelectServer(test.input)
		if err != nil {
			t.Fatalf("Expected err to be nil but got %s\n", err)
		}
		if test.expected != selectedServer {
			t.Errorf("Failed on %d with Path %s: Expected '%s' but got '%s'\n", i, test.input.URLPath,
				test.expected, selectedServer)
		}
	}

}
