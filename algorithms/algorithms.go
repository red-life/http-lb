package algorithms

import (
	"github.com/red-life/http-lb"
)

type LoadBalancingAlgorithm interface {
	ChooseBackend(r http_lb.Request) string
}
