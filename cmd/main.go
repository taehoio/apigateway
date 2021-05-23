package main

import (
	"context"

	"cloud.google.com/go/profiler"
	"github.com/sirupsen/logrus"
	"github.com/taehoio/apigateway/config"
	"github.com/taehoio/apigateway/server"
)

func main() {
	cfg := config.NewConfig(config.NewSetting())

	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := logrus.StandardLogger()

	if err := setUpProfiler(cfg); err != nil {
		log.Fatal(err)
	}

	log.WithField("setting", cfg.Setting()).Info("Starting server...")

	ctx := context.Background()
	srv, err := server.NewHTTPServer(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func setUpProfiler(cfg config.Config) error {
	if !shouldProfile(cfg) {
		return nil
	}

	pc := profiler.Config{
		Service: cfg.Setting().ServiceName,
	}
	if err := profiler.Start(pc); err != nil {
		return err
	}
	return nil
}

func shouldProfile(cfg config.Config) bool {
	profilingEnvs := []string{"production", "staging"}
	return in(profilingEnvs, cfg.Setting().Env)
}

func in(arr []string, s string) bool {
	for _, sia := range arr {
		if sia == s {
			return true
		}
	}
	return false
}
