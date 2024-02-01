package algorithms

import "net/http"

type LoadBalancingAlgorithm interface {
	ChooseBackend(r http.Request) string
}
