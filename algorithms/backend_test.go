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

func TestBackendAddrsManager_RegisterBackend(t *testing.T) {
	addrs := []string{"addr 1", "addr 2", "addr 3", "addr 4", "addr 5"}
	tests := []struct {
		input    string
		expected error
	}{
		{addrs[0], nil},
		{addrs[1], nil},
		{addrs[2], nil},
		{addrs[3], nil},
		{addrs[4], nil},
		{addrs[0], http_lb.ErrBackendExists},
		{addrs[1], http_lb.ErrBackendExists},
		{addrs[2], http_lb.ErrBackendExists},
	}
	logger, _ := zap.NewDevelopment()
	backendPool := algorithms.NewBackendPool([]string{}, logger)
	for _, test := range tests {
		err := backendPool.RegisterBackend(test.input)
		if err != test.expected {
			t.Fatalf("Expected err to be %s but got %s", test.expected, err)
		}
	}
	if !isEqual(addrs, backendPool.Backends()) {
		t.Fatalf("Expected all backend addrs to be registered but got %+v", backendPool.Backends())
	}
}

func TestBackendAddrsManager_UnregisterBackend(t *testing.T) {
	addrs := []string{"addr 1", "addr 2", "addr 3", "addr 4", "addr 5"}
	tests := []struct {
		input    string
		expected error
	}{
		{addrs[0], nil},
		{addrs[1], nil},
		{addrs[2], nil},
		{addrs[3], nil},
		{addrs[4], nil},
		{addrs[0], http_lb.ErrBackendNotExist},
		{addrs[1], http_lb.ErrBackendNotExist},
		{addrs[2], http_lb.ErrBackendNotExist},
	}
	logger, _ := zap.NewDevelopment()
	backendPool := algorithms.NewBackendPool(addrs, logger)
	for _, test := range tests {
		err := backendPool.UnregisterBackend(test.input)
		if err != test.expected {
			t.Fatalf("Expected err to be %s but got %s\n", test.expected, err)
		}
	}
	if len(backendPool.Backends()) != 0 {
		t.Fatalf("Expected backend addrs' length to be 0 but got %d\n", len(backendPool.Backends()))
	}
}
