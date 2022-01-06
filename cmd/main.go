package main

import (
	"context"

	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"

	"github.com/taehoio/apigateway/config"
	"github.com/taehoio/apigateway/server"
)

func main() {
	cfg := config.NewConfig(config.NewSetting())

	log := cfg.Logger()
	log.WithField("setting", cfg.Setting()).Info("Starting server...")

	if err := runServer(cfg); err != nil {
		log.Fatal(err)
	}
}

func runServer(cfg config.Config) error {
	if cfg.ShouldProfile() {
		if err := setUpProfiler(cfg.ServiceName()); err != nil {
			return err
		}
	}

	if cfg.ShouldTrace() {
		if err := setUpTracing(); err != nil {
			return err
		}
	}

	ctx := context.Background()
	srv, err := server.NewHTTPServer(ctx, cfg)
	if err != nil {
		return err
	}

	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func setUpProfiler(serviceName string) error {
	pc := profiler.Config{
		Service: serviceName,
	}
	if err := profiler.Start(pc); err != nil {
		return err
	}
	return nil
}

func setUpTracing() error {
	exporter, err := stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		return err
	}

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	return nil
}
