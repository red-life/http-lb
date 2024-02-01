package http_lb

type LoadBalancingAlgorithm interface {
	ChooseBackend(r Request) string
}

type Request struct {
	RemoteIP string
	URLPath  string
}
