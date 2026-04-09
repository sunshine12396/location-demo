package stdlog

import (
	"github.com/example/location-demo/pkg/stdlog/graylog"
	"github.com/example/location-demo/pkg/stdlog/sentry"
)

type Option func(*Config)

func EnableGraylog(url string) Option {
	return func(cfg *Config) {
		cfg.enableGraylog = true
		cfg.graylogCore = graylog.New(graylog.Config{
			Address:    url,
			ServerName: cfg.ServerName,
			Version:    cfg.Version,
		}, newDefaultConfig())
	}
}

func EnableSentry(url string) Option {
	return func(cfg *Config) {
		cfg.enableSentry = true
		cfg.sentryCore = sentry.New(sentry.Config{
			DSN:         url,
			Version:     cfg.Version,
			ServerName:  cfg.ServerName,
			Environment: cfg.Environment,
		}, newDefaultConfig())
	}
}
