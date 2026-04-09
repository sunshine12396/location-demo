package grpc

import "github.com/example/location-demo/pkg/stdlog"

type Option func(*Config)

func WithStdLogger(logger stdlog.Logger) Option {
	return func(cfg *Config) {
		cfg.writer = logger
	}
}
