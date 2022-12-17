package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	googlemetadata "cloud.google.com/go/compute/metadata"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/taehoio/apigateway/config"
	authv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/auth/v1"
	baemincryptov1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/baemincrypto/v1"
	carv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/car/v1"
	oneononev1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/oneonone/v1"
	texttospeechv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/texttospeech/v1"
	userv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/user/v1"
)

func getIDTokenInGCP(serviceURL string) (string, error) {
	tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", serviceURL)
	return googlemetadata.Get(tokenURL)
}

func withMarshalerOption() runtime.ServeMuxOption {
	return runtime.WithMarshalerOption(
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
	)
}

func withMetadata(cfg config.Config, serviceNameURLMap map[string]string) runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		serviceName := strings.Split(req.URL.Path, "/")[1]
		serviceURL := serviceNameURLMap[serviceName]

		md := metadata.MD{}

		if cfg.Setting().IsInGCP {
			idToken, err := getIDTokenInGCP(serviceURL)
			if err != nil {
				logrus.StandardLogger().Error(err)
			}
			md.Append("Authorization", "Bearer "+idToken)
		}

		return md
	})
}

func withSecureOption(cfg config.Config) (grpc.DialOption, error) {
	secureOpt := grpc.WithTransportCredentials(insecure.NewCredentials())

	if cfg.Setting().ShouldUseGRPCClientTLS {
		creds, err := credentials.NewClientTLSFromFile(cfg.Setting().CACertFile, "")
		if err != nil {
			return nil, err
		}

		secureOpt = grpc.WithTransportCredentials(creds)
	}

	return secureOpt, nil
}

func registerBaemincryptoService(ctx context.Context, gwMux *runtime.ServeMux, endpoint string, opts ...grpc.DialOption) error {
	baemincryptov1Conn, err := grpc.Dial(
		endpoint,
		opts...,
	)
	if err != nil {
		return err
	}

	if err := baemincryptov1.RegisterBaemincryptoServiceHandler(
		ctx,
		gwMux,
		baemincryptov1Conn,
	); err != nil {
		return err
	}

	return nil
}

func registerUserService(ctx context.Context, gwMux *runtime.ServeMux, endpoint string, opts ...grpc.DialOption) error {
	userv1Conn, err := grpc.Dial(
		endpoint,
		opts...,
	)
	if err != nil {
		return err
	}

	if err := userv1.RegisterUserServiceHandler(
		ctx,
		gwMux,
		userv1Conn,
	); err != nil {
		return err
	}

	return nil
}

func registerAuthService(ctx context.Context, gwMux *runtime.ServeMux, endpoint string, opts ...grpc.DialOption) error {
	authv1Conn, err := grpc.Dial(
		endpoint,
		opts...,
	)
	if err != nil {
		return err
	}

	if err := authv1.RegisterAuthServiceHandler(
		ctx,
		gwMux,
		authv1Conn,
	); err != nil {
		return err
	}

	return nil
}

func registerOneononeService(ctx context.Context, gwMux *runtime.ServeMux, endpoint string, opts ...grpc.DialOption) error {
	oneononev1Conn, err := grpc.Dial(
		endpoint,
		opts...,
	)
	if err != nil {
		return err
	}

	if err := oneononev1.RegisterOneononeServiceHandler(
		ctx,
		gwMux,
		oneononev1Conn,
	); err != nil {
		return err
	}

	return nil
}

func registerTexttospeechService(ctx context.Context, gwMuc *runtime.ServeMux, endpoint string, opts ...grpc.DialOption) error {
	texttospeechv1Conn, err := grpc.Dial(
		endpoint,
		opts...,
	)
	if err != nil {
		return err
	}

	if err := texttospeechv1.RegisterTexttospeechServiceHandler(
		ctx,
		gwMuc,
		texttospeechv1Conn,
	); err != nil {
		return err
	}

	return nil
}

func registerCarService(ctx context.Context, gwMuc *runtime.ServeMux, endpoint string, opts ...grpc.DialOption) error {
	carv1Conn, err := grpc.Dial(
		endpoint,
		opts...,
	)
	if err != nil {
		return err
	}

	if err := carv1.RegisterCarServiceHandler(
		ctx,
		gwMuc,
		carv1Conn,
	); err != nil {
		return err
	}

	return nil
}

func grpcGWMux(ctx context.Context, cfg config.Config) (*runtime.ServeMux, error) {
	serviceNameURLMap := map[string]string{
		"baemincrypto": cfg.Setting().BaemincryptoGRPCServiceURL,
		"user":         cfg.Setting().UserGRPCServiceURL,
		"auth":         cfg.Setting().AuthGRPCServiceURL,
		"oneonone":     cfg.Setting().OneononeGRPCServiceURL,
		"texttospeech": cfg.Setting().TexttospeechGRPCServiceURL,
		"car":          cfg.Setting().CarGRPCServiceURL,
	}

	gwMux := runtime.NewServeMux(
		withMarshalerOption(),
		withMetadata(cfg, serviceNameURLMap),
	)

	secureOpt, err := withSecureOption(cfg)
	if err != nil {
		return nil, err
	}

	if err := registerBaemincryptoService(
		ctx,
		gwMux,
		cfg.Setting().BaemincryptoGRPCServiceEndpoint,
		secureOpt,
		grpc.WithUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
		),
	); err != nil {
		return nil, err
	}

	if err := registerUserService(
		ctx,
		gwMux,
		cfg.Setting().UserGRPCServiceEndpoint,
		secureOpt,
		grpc.WithUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
		),
	); err != nil {
		return nil, err
	}

	if err := registerAuthService(
		ctx,
		gwMux,
		cfg.Setting().AuthGRPCServiceEndpoint,
		secureOpt,
		grpc.WithUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
		),
	); err != nil {
		return nil, err
	}

	if err := registerOneononeService(
		ctx,
		gwMux,
		cfg.Setting().OneononeGRPCServiceEndpoint,
		secureOpt,
		grpc.WithUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
		),
	); err != nil {
		return nil, err
	}

	if err := registerTexttospeechService(
		ctx,
		gwMux,
		cfg.Setting().TexttospeechGRPCServiceEndpoint,
		secureOpt,
		grpc.WithUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
		),
	); err != nil {
		return nil, err
	}

	if err := registerCarService(
		ctx,
		gwMux,
		cfg.Setting().CarGRPCServiceEndpoint,
		secureOpt,
		grpc.WithUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
		),
	); err != nil {
		return nil, err
	}

	return gwMux, nil
}
