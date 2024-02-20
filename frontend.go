package http_lb

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

var _ GracefulShutdown = (*Frontend)(nil)

func NewFrontend(listenAddr string, tls *TLSOptions, reqForwarder RequestForwarder, logger *zap.Logger) *Frontend {
	return &Frontend{
		listenAddr:   listenAddr,
		mux:          http.NewServeMux(),
		tls:          tls,
		reqForwarder: reqForwarder,
		logger:       logger,
	}
}

type Frontend struct {
	listenAddr   string
	mux          *http.ServeMux
	tls          *TLSOptions
	reqForwarder RequestForwarder
	logger       *zap.Logger
	httpServer   http.Server
}

func (f *Frontend) Handler(rw http.ResponseWriter, r *http.Request) {
	err := f.reqForwarder.Forward(rw, r)
	if err != nil {
		f.logger.Error("failed to forward request", zap.String("ip", r.RemoteAddr),
			zap.String("path", r.RequestURI), zap.String("method", r.Method), zap.Error(err))
		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte(err.Error()))
		return
	}
	f.logger.Info("request forwarded", zap.String("ip", r.RemoteAddr),
		zap.String("path", r.RequestURI), zap.String("method", r.Method))
}

func (f *Frontend) Run() error {
	f.mux.HandleFunc("/", f.Handler)
	f.httpServer.Handler = f.mux
	f.httpServer.Addr = f.listenAddr
	if f.tls != nil {
		f.logger.Info("started listening tls", zap.String("listen", f.listenAddr),
			zap.String("certFile", f.tls.CertFile), zap.String("keyFile", f.tls.KeyFile))
		return f.httpServer.ListenAndServeTLS(f.tls.CertFile, f.tls.KeyFile)
	}
	f.logger.Info("started listening", zap.String("listen", f.listenAddr))
	return f.httpServer.ListenAndServe()
}

func (f *Frontend) Shutdown() error {
	return f.httpServer.Shutdown(context.Background())
}
