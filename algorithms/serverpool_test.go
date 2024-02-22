package algorithms_test

import (
	"github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"go.uber.org/zap"
	"testing"
)

func find(slice []string, value string) int {
	for idx, val := range slice {
		if val == value {
			return idx
		}
	}
	return -1
}

func isEqual(a, b []string) bool {
	for _, val := range a {
		if find(b, val) == -1 {
			return false
		}
	}
	return true
}

func TestServerPool_RegisterServer(t *testing.T) {
	servers := []string{"server 1", "server 2", "server 3", "server 4", "server 5"}
	tests := []struct {
		input    string
		expected error
	}{
		{servers[0], nil},
		{servers[1], nil},
		{servers[2], nil},
		{servers[3], nil},
		{servers[4], nil},
		{servers[0], http_lb.ErrServerExists},
		{servers[1], http_lb.ErrServerExists},
		{servers[2], http_lb.ErrServerExists},
	}
	logger, _ := zap.NewDevelopment()
	serverPool := algorithms.NewServerPool([]string{}, logger)
	for _, test := range tests {
		err := serverPool.RegisterServer(test.input)
		if err != test.expected {
			t.Fatalf("Expected err to be %s but got %s", test.expected, err)
		}
	}
	if !isEqual(servers, serverPool.Servers()) {
		t.Fatalf("Expected all servers to be registered but got %+v", serverPool.Servers())
	}
}

func TestServerPool_UnregisterServer(t *testing.T) {
	servers := []string{"server 1", "server 2", "server 3", "server 4", "server 5"}
	tests := []struct {
		input    string
		expected error
	}{
		{servers[0], nil},
		{servers[1], nil},
		{servers[2], nil},
		{servers[3], nil},
		{servers[4], nil},
		{servers[0], http_lb.ErrServerNotExist},
		{servers[1], http_lb.ErrServerNotExist},
		{servers[2], http_lb.ErrServerNotExist},
	}
	logger, _ := zap.NewDevelopment()
	serverPool := algorithms.NewServerPool(servers, logger)
	for _, test := range tests {
		err := serverPool.UnregisterServer(test.input)
		if err != test.expected {
			t.Fatalf("Expected err to be %s but got %s\n", test.expected, err)
		}
	}
	if len(serverPool.Servers()) != 0 {
		t.Fatalf("Expected servers length to be 0 but got %d\n", len(serverPool.Servers()))
	}
}
