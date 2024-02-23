package http_lb

import (
	"errors"
)

var (
	ErrServerExists             = errors.New("server already exists")
	ErrServerNotExist           = errors.New("server doesn't exist")
	ErrNoHealthyServerAvailable = errors.New("no healthy server is available")
)
