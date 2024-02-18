package http_lb

import "net/http"

func NewFrontend(listenAddr string, tls *TLSOptions) *Frontend {
	return &Frontend{
		listenAddr: listenAddr,
		mux:        http.NewServeMux(),
		tls:        tls,
	}
}

type TLSOptions struct {
	CertFile string
	KeyFile  string
}

type Frontend struct {
	listenAddr   string
	mux          *http.ServeMux
	tls          *TLSOptions
	reqForwarder RequestForwarder
}

func (f *Frontend) Handler(rw http.ResponseWriter, r *http.Request) {
	err := f.reqForwarder.Forward(rw, r)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte(err.Error()))
	}
}

func (f *Frontend) Run() error {
	f.mux.HandleFunc("/", f.Handler)
	if f.tls != nil {
		return http.ListenAndServeTLS(f.listenAddr, f.tls.CertFile, f.tls.KeyFile, f.mux)
	}
	return http.ListenAndServe(f.listenAddr, f.mux)
}
