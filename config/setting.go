package config

import (
	"log"
	"os"
	"strconv"
)

type Setting struct {
	ServiceName               string
	HTTPServerPort            string
	Env                       string
	GracefulShutdownTimeoutMs int

	IsGRPCInsecure                  bool
	CertFile                        string
	BaemincryptoGRPCServiceEndpoint string
}

func NewSetting() Setting {
	return Setting{
		ServiceName:               "apigateway",
		HTTPServerPort:            getEnv("HTTP_SERVER_PORT", "8080"),
		Env:                       getEnv("ENV", "development"),
		GracefulShutdownTimeoutMs: mustAtoi(getEnv("GRACEFUL_SHUTDOWN_TIMEOUT_MS", "10000")),

		IsGRPCInsecure:                  mustAtob(getEnv("IS_GRPC_INSECURE", "true")),
		CertFile:                        getEnv("CERT_FILE", "/etc/ssl/certs/ca-certificates.crt"),
		BaemincryptoGRPCServiceEndpoint: getEnv("BAEMINCRYPTO_GRPC_SERVICE_ENDPOINT", "localhost:50051"),
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
