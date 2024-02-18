package algorithms_test

import (
	"github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"go.uber.org/zap"
	"testing"
)

func TestURLHash_ChooseBackend(t *testing.T) {
	backendAddrs := []string{
		"addr 1",
		"addr 2",
		"addr 3",
		"addr 4",
		"addr 5",
		"addr 6",
	}
	logger, _ := zap.NewDevelopment()
	addrMng := algorithms.NewBackendAddrsManager(backendAddrs, logger)
	urlHash := algorithms.NewURLHash(http_lb.Hash, addrMng, logger)
	tests := []struct {
		input    http_lb.Request
		expected string
	}{
		{http_lb.Request{URLPath: "/"}, backendAddrs[int(http_lb.Hash("/"))%len(backendAddrs)]},
		{http_lb.Request{URLPath: "/home"}, backendAddrs[int(http_lb.Hash("/home"))%len(backendAddrs)]},
		{http_lb.Request{URLPath: "/auth/login"}, backendAddrs[int(http_lb.Hash("/auth/login"))%len(backendAddrs)]},
		{http_lb.Request{URLPath: "/api/v1"}, backendAddrs[int(http_lb.Hash("/api/v1"))%len(backendAddrs)]},
	}
	for i, test := range tests {
		chosenBackend, err := urlHash.ChooseBackend(test.input)
		if err != nil {
			t.Fatalf("Expected err to be nil but got %s\n", err)
		}
		if test.expected != chosenBackend {
			t.Errorf("Failed on %d with Path %s: Expected '%s' but got '%s'\n", i, test.input.URLPath,
				test.expected, chosenBackend)
		}
	}

}
