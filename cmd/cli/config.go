package cli

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Algorithm   string      `mapstructure:"algorithm" validate:"oneof=round-robin sticky-round-robin weighted-round-robin url-hash ip-hash random"`
	LogLevel    string      `mapstructure:"log_level" validate:"oneof=debug info warn error"`
	Frontend    Frontend    `mapstructure:"frontend" validate:"required"`
	Backend     Backend     `mapstructure:"backend" validate:"required"`
	HealthCheck HealthCheck `mapstructure:"health_check" validate:"required"`
}

type Frontend struct {
	Listen string `mapstructure:"listen" validate:"required"`
	TLS    *TLS   `mapstructure:"tls"`
}

type TLS struct {
	CertFile string `mapstructure:"cert" validate:"filepath,required"`
	KeyFile  string `mapstructure:"key" validate:"filepath,required"`
}

type Backend struct {
	Servers   []string      `mapstructure:"servers" validate:"required"`
	Timeout   time.Duration `mapstructure:"timeout" validate:"min=1ms,required"`
	KeepAlive KeepAlive     `mapstructure:"keep_alive" validate:"required"`
}

type KeepAlive struct {
	MaxIdle     int           `mapstructure:"max_idle_connections" validate:"min=1,required"`
	IdleTimeout time.Duration `mapstructure:"idle_connection_timeout" validate:"min=1ms,required"`
}

type HealthCheck struct {
	Endpoint           string        `mapstructure:"endpoint" validate:"uri,required"`
	ExpectedStatusCode int           `mapstructure:"expected_status_code" validate:"min=100,max=599,required"`
	Interval           time.Duration `mapstructure:"interval" validate:"min=1ms,required"`
	Timeout            time.Duration `mapstructure:"timeout" validate:"min=1ms,required"`
}

func ParseAndValidateConfig(configFilePath string) (Config, error) {
	var config Config
	validate := validator.New()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}
	err = validate.Struct(config)
	if err != nil {
		return config, err
	}
	return config, nil
}
