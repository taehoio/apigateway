package config

import (
	"log"
	"os"
	"strconv"
)

type Setting struct {
	serviceName                     string
	httpServerPort                  string
	env                             string
	gracefulShutdownTimeoutMs       int
	shouldProfile                   bool
	shouldTrace                     bool
	shouldUseGRPCClientTLS          bool
	caCertFile                      string
	baemincryptoGRPCServiceEndpoint string
}

func NewSetting() Setting {
	return Setting{
		serviceName:                     "apigateway",
		httpServerPort:                  getEnv("HTTP_SERVER_PORT", "8080"),
		env:                             getEnv("ENV", "development"),
		gracefulShutdownTimeoutMs:       mustAtoi(getEnv("GRACEFUL_SHUTDOWN_TIMEOUT_MS", "10000")),
		shouldProfile:                   mustAtob(getEnv("SHOULD_PROFILE", "false")),
		shouldTrace:                     mustAtob(getEnv("SHOULD_TRACE", "false")),
		shouldUseGRPCClientTLS:          mustAtob(getEnv("SHOULD_USE_GRPC_CLIENT_TLS", "false")),
		caCertFile:                      getEnv("CA_CERT_FILE", "/etc/ssl/certs/ca-certificates.crt"),
		baemincryptoGRPCServiceEndpoint: getEnv("BAEMINCRYPTO_GRPC_SERVICE_ENDPOINT", "localhost:50051"),
	}
}

func (s Setting) ServiceName() string {
	return s.serviceName
}

func (s Setting) HTTPServerPort() string {
	return s.httpServerPort
}

func (s Setting) ENV() string {
	return s.env
}

func (s Setting) GracefulShutdownTimeoutMs() int {
	return s.gracefulShutdownTimeoutMs
}

func (s Setting) ShouldProfile() bool {
	return s.shouldProfile
}

func (s Setting) ShouldTrace() bool {
	return s.shouldTrace
}

func (s Setting) ShouldUseGRPCClientTLS() bool {
	return s.shouldUseGRPCClientTLS
}

func (s Setting) CACertFile() string {
	return s.caCertFile
}

func (s Setting) BaemincryptoGRPCServiceEndpoint() string {
	return s.baemincryptoGRPCServiceEndpoint
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
