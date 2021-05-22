package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/taehoio/apigateway/config"
	"github.com/taehoio/apigateway/server"
)

func main() {
	cfg := config.NewConfig(config.NewSetting())

	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := logrus.StandardLogger()

	port := cfg.Setting().HTTPServerPort
	log.WithField("port", port).Info("server starting...")

	ctx := context.Background()
	srv, err := server.NewHTTPServer(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
