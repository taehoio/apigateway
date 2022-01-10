package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/taehoio/apigateway/config"
)

func newRouter(ctx context.Context, cfg config.Config) (*mux.Router, error) {
	rtr := mux.NewRouter()

	ec := newEcho()
	rtr.HandleFunc("/", ec.ServeHTTP)

	gwMux, err := grpcGWMux(ctx, cfg)
	if err != nil {
		return nil, err
	}
	rtr.Handle("/{serviceName}/{version}/{rest:.*}", gwMux)

	return rtr, nil
}

func NewHTTPServer(ctx context.Context, cfg config.Config) (*http.Server, error) {
	rtr, err := newRouter(ctx, cfg)
	if err != nil {
		return nil, err
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", rtr)

	httpHandler := otelhttp.NewHandler(httpMux, "server")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Setting().HTTPServerPort),
		Handler: httpHandler,
	}

	return srv, nil
}
