package http_lb

import (
	"errors"
)

var (
	ErrBackendExists     = errors.New("backend already exists")
	ErrBackendNotExist   = errors.New("backend doesn't exist")
	ErrNoServerAvailable = errors.New("no server is available")
)
