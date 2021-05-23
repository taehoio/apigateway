package main

import (
	"context"

	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/sirupsen/logrus"
	"github.com/taehoio/apigateway/config"
	"github.com/taehoio/apigateway/server"
	"go.opencensus.io/trace"
)

func main() {
	cfg := config.NewConfig(config.NewSetting())

	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := logrus.StandardLogger()

	if cfg.Setting().ShouldProfile {
		if err := setUpProfiler(cfg.Setting().ServiceName); err != nil {
			log.Fatal(err)
		}
	}

	if cfg.Setting().ShouldTrace {
		if err := setUpTracing(); err != nil {
			log.Fatal(err)
		}
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
