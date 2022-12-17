package config

import (
	"log"
	"os"
	"strconv"
)

type Setting struct {
	ServiceName                     string
	HTTPServerPort                  int
	Env                             string
	GracefulShutdownTimeoutMs       int
	ShouldProfile                   bool
	ShouldTrace                     bool
	ShouldUseGRPCClientTLS          bool
	CACertFile                      string
	BaemincryptoGRPCServiceEndpoint string
	BaemincryptoGRPCServiceURL      string
	UserGRPCServiceEndpoint         string
	UserGRPCServiceURL              string
	AuthGRPCServiceEndpoint         string
	AuthGRPCServiceURL              string
	OneononeGRPCServiceEndpoint     string
	OneononeGRPCServiceURL          string
	TexttospeechGRPCServiceEndpoint string
	TexttospeechGRPCServiceURL      string
	CarGRPCServiceEndpoint          string
	CarGRPCServiceURL               string
	IsInGCP                         bool
}

func NewSetting() Setting {
	return Setting{
		ServiceName:                     "apigateway",
		HTTPServerPort:                  mustAtoi(getEnv("HTTP_SERVER_PORT", "8080")),
		Env:                             getEnv("ENV", "development"),
		GracefulShutdownTimeoutMs:       mustAtoi(getEnv("GRACEFUL_SHUTDOWN_TIMEOUT_MS", "5000")),
		ShouldProfile:                   mustAtob(getEnv("SHOULD_PROFILE", "false")),
		ShouldTrace:                     mustAtob(getEnv("SHOULD_TRACE", "false")),
		ShouldUseGRPCClientTLS:          mustAtob(getEnv("SHOULD_USE_GRPC_CLIENT_TLS", "false")),
		CACertFile:                      getEnv("CA_CERT_FILE", "/etc/ssl/certs/ca-certificates.crt"),
		BaemincryptoGRPCServiceEndpoint: getEnv("BAEMINCRYPTO_GRPC_SERVICE_ENDPOINT", "baemincrypto-5hwa5dthla-an.a.run.app:443"),
		BaemincryptoGRPCServiceURL:      getEnv("BAEMINCRYPTO_GRPC_SERVICE_URL", "https://baemincrypto-5hwa5dthla-an.a.run.app"),
		UserGRPCServiceEndpoint:         getEnv("USER_GRPC_SERVICE_ENDPOINT", "user-5hwa5dthla-an.a.run.app:443"),
		UserGRPCServiceURL:              getEnv("USER_GRPC_SERVICE_URL", "https://user-5hwa5dthla-an.a.run.app"),
		AuthGRPCServiceEndpoint:         getEnv("AUTH_GRPC_SERVICE_ENDPOINT", "auth-5hwa5dthla-an.a.run.app:443"),
		AuthGRPCServiceURL:              getEnv("AUTH_GRPC_SERVICE_URL", "https://auth-5hwa5dthla-an.a.run.app"),
		OneononeGRPCServiceEndpoint:     getEnv("ONEONONE_GRPC_SERVICE_ENDPOINT", "oneonone-5hwa5dthla-an.a.run.app:443"),
		OneononeGRPCServiceURL:          getEnv("ONEONONE_GRPC_SERVICE_URL", "https://oneonone-5hwa5dthla-an.a.run.app"),
		TexttospeechGRPCServiceEndpoint: getEnv("TEXTTOSPEECH_GRPC_SERVICE_ENDPOINT", "texttospeech-5hwa5dthla-an.a.run.app:443"),
		TexttospeechGRPCServiceURL:      getEnv("TEXTTOSPEECH_GRPC_SERVICE_URL", "https://texttospeech-5hwa5dthla-an.a.run.app"),
		CarGRPCServiceEndpoint:          getEnv("CAR_GRPC_SERVICE_ENDPOINT", "car-5hwa5dthla-an.a.run.app:443"),
		CarGRPCServiceURL:               getEnv("CAR_GRPC_SERVICE_URL", "https://car-5hwa5dthla-an.a.run.app"),
		IsInGCP:                         mustAtob(getEnv("IS_IN_GCP", "false")),
	}
}

func MockSetting() Setting {
	return NewSetting()
}

func getEnv(key, defaultValue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		if defaultValue != "" {
			value = defaultValue
		} else {
			log.Fatalf("missing required environment variable: %v", key)
		}
	}
	return value
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Panic(err)
	}
	return i
}

func mustAtob(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Panic(err)
	}
	return b
}
