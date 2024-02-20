package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	http_lb "github.com/red-life/http-lb"
	"github.com/red-life/http-lb/cmd/cli"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func Shutdown(services []http_lb.GracefulShutdown, logger *zap.Logger) {
	for _, s := range services {
		if err := s.Shutdown(); err != nil {
			logger.Error("error occurred while shutting the service", zap.Error(err))
		}
	}
}

func isDev() bool {
	return os.Getenv("development") != ""
}

func getLogger(logLevel string) *zap.Logger {
	rawJSON := []byte(fmt.Sprintf(`{
		"level": "%s",
		"development": %v,
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
			"messageKey": "msg",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`, logLevel, isDev()))
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger := zap.Must(cfg.Build())
	return logger
}

func main() {
	configFile := flag.String("c", "config.yaml", "yaml configuration file path")
	flag.Parse()
	config, err := cli.ParseAndValidateConfig(*configFile)
	if err != nil {
		panic(fmt.Sprintf("failed to parse and validate configuration file: %s", err))
	}
	logger := getLogger(config.LogLevel)
	defer logger.Sync()
	frontend, healthCheck := cli.Factory(config, logger)
	healthCheck.Run()
	go func() {
		logger.Error("frontend stopped", zap.Error(frontend.Run()))
	}()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	Shutdown([]http_lb.GracefulShutdown{frontend, healthCheck}, logger)
	logger.Info("stopped")
}
