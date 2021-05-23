package server

import (
	"context"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/taehoio/apigateway/config"
	baemincryptov1 "github.com/taehoio/idl/gen/go/services/baemincrypto/v1"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"
)

func newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.HEAD("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"service":   "apigateway",
			"version":   "v1",
			"host":      c.Request().Host,
			"timestamp": time.Now().UTC().String(),
		})
	})

	return e
}

func newGRPCGateway(ctx context.Context, cfg config.Config) (*runtime.ServeMux, error) {
	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		),
	)

	secureOpt := grpc.WithInsecure()
	if cfg.Setting().ShouldUseGRPCClientTLS() {
		creds, err := credentials.NewClientTLSFromFile(cfg.Setting().CACertFile(), "")
		if err != nil {
			return nil, err
		}
		secureOpt = grpc.WithTransportCredentials(creds)
	}

	baemincryptov1Conn, err := grpc.Dial(
		cfg.Setting().BaemincryptoGRPCServiceEndpoint(),
		secureOpt,
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
	)
	if err != nil {
		return nil, err
	}
	if err := baemincryptov1.RegisterBaemincryptoServiceHandler(
		ctx,
		gwMux,
		baemincryptov1Conn,
	); err != nil {
		return nil, err
	}

	return gwMux, nil
}

func NewHTTPServer(ctx context.Context, cfg config.Config) (*http.Server, error) {
	r := mux.NewRouter()

	e := newEcho()
	r.HandleFunc("/", e.ServeHTTP)

	gwMux, err := newGRPCGateway(ctx, cfg)
	if err != nil {
		return nil, err
	}
	r.Handle("/{serviceName}/{version}/{rest:.*}", gwMux)

	httpMux := http.NewServeMux()
	httpMux.Handle("/", r)

	httpHandler := &ochttp.Handler{
		Propagation: &propagation.HTTPFormat{},
		Handler:     httpMux,
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Setting().HTTPServerPort(),
		Handler: httpHandler,
	}

	return srv, nil
}
