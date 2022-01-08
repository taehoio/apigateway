package server

import (
	"context"
	"fmt"
	"net/http"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/gorilla/mux"
	"go.opencensus.io/plugin/ochttp"

	"github.com/taehoio/apigateway/config"
)

func newRouter(ctx context.Context, cfg config.Config) (*mux.Router, error) {
	rtr := mux.NewRouter()

	ec := newEcho()
	rtr.HandleFunc("/", ec.ServeHTTP)

	baemincryptoGatewayMux, err := newBaemincryptoServiceGRPCGatewayMux(ctx, cfg)
	if err != nil {
		return nil, err
	}
	rtr.Handle("/baemincrypto/{version}/{rest:.*}", baemincryptoGatewayMux)

	userGatewayMux, err := newUserServiceGRPCGatewayMux(ctx, cfg)
	if err != nil {
		return nil, err
	}
	rtr.Handle("/user/{version}/{rest:.*}", userGatewayMux)

	return rtr, nil
}

func handlerWithTracingPropagation(httpMux *http.ServeMux) *ochttp.Handler {
	return &ochttp.Handler{
		Propagation: &propagation.HTTPFormat{},
		Handler:     httpMux,
	}
}

func NewHTTPServer(ctx context.Context, cfg config.Config) (*http.Server, error) {
	rtr, err := newRouter(ctx, cfg)
	if err != nil {
		return nil, err
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", rtr)

	httpHandler := handlerWithTracingPropagation(httpMux)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPServerPort()),
		Handler: httpHandler,
	}

	return srv, nil
}
