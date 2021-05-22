package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/taehoio/apigateway/config"
	baemincryptov1 "github.com/taehoio/idl/gen/go/services/baemincrypto/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Logger())

	e.HEAD("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"Hello": "World"})
	})

	return e
}

func newGRPCGateway(ctx context.Context, cfg config.Config) (*runtime.ServeMux, error) {
	gwmux := runtime.NewServeMux(
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
	options := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	if err := baemincryptov1.RegisterBaemincryptoServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		cfg.Setting().BaemincryptoServiceEndpoint,
		options,
	); err != nil {
		return nil, err
	}

	return gwmux, nil
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

	srv := &http.Server{
		Addr:    ":" + cfg.Setting().HTTPServerPort,
		Handler: httpMux,
	}

	return srv, nil
}
