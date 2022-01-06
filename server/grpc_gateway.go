package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	googlemetadata "cloud.google.com/go/compute/metadata"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/taehoio/apigateway/config"
	baemincryptov1 "github.com/taehoio/idl/gen/go/services/baemincrypto/v1"
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

func withMetadata(cfg config.Config) runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		md := metadata.MD{}

		if cfg.IsInGCP() {
			idToken, err := getIDTokenInGCP(strings.Join([]string{
				cfg.BaemincryptoGRPCServiceURL(),
			}, ","))
			if err != nil {
				logrus.StandardLogger().Error(err)
			}
			md.Append("Authorization", "Bearer "+idToken)
		} else {
			idToken := cfg.IDToken()
			md.Append("Authorization", "Bearer "+idToken)
		}

		return md
	})
}

func withSecureOption(cfg config.Config) (grpc.DialOption, error) {
	secureOpt := grpc.WithTransportCredentials(insecure.NewCredentials())

	if cfg.ShouldUseGRPCClientTLS() {
		creds, err := credentials.NewClientTLSFromFile(cfg.CACertFile(), "")
		if err != nil {
			return nil, err
		}

		secureOpt = grpc.WithTransportCredentials(creds)
	}

	return secureOpt, nil
}

func withTracingStatsHandler() grpc.DialOption {
	return grpc.WithStatsHandler(&ocgrpc.ClientHandler{})
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

func newGRPCGatewayMux(ctx context.Context, cfg config.Config) (*runtime.ServeMux, error) {
	gwMux := runtime.NewServeMux(
		withMarshalerOption(),
		withMetadata(cfg),
	)

	secureOpt, err := withSecureOption(cfg)
	if err != nil {
		return nil, err
	}

	if err := registerBaemincryptoService(
		ctx,
		gwMux,
		cfg.BaemincryptoGRPCServiceEndpoint(),
		secureOpt,
		withTracingStatsHandler(),
	); err != nil {
		return nil, err
	}

	return gwMux, nil
}
