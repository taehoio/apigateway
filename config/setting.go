package config

import (
	"log"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Setting struct {
	serviceName                     string
	httpServerPort                  int
	env                             string
	shouldProfile                   bool
	shouldTrace                     bool
	shouldUseGRPCClientTLS          bool
	caCertFile                      string
	baemincryptoGRPCServiceEndpoint string
	baemincryptoGRPCServiceURL      string
	isInGCP                         bool
	idToken                         string
	logger                          *logrus.Logger
}

func NewSetting() Setting {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	return Setting{
		serviceName:                     "apigateway",
		httpServerPort:                  mustAtoi(getEnv("HTTP_SERVER_PORT", "8080")),
		env:                             getEnv("ENV", "development"),
		shouldProfile:                   mustAtob(getEnv("SHOULD_PROFILE", "false")),
		shouldTrace:                     mustAtob(getEnv("SHOULD_TRACE", "false")),
		shouldUseGRPCClientTLS:          mustAtob(getEnv("SHOULD_USE_GRPC_CLIENT_TLS", "false")),
		caCertFile:                      getEnv("CA_CERT_FILE", "/etc/ssl/certs/ca-certificates.crt"),
		baemincryptoGRPCServiceEndpoint: getEnv("BAEMINCRYPTO_GRPC_SERVICE_ENDPOINT", "baemincrypto-5hwa5dthla-an.a.run.app:443"),
		baemincryptoGRPCServiceURL:      getEnv("BAEMINCRYPTO_GRPC_SERVICE_URL", "https://baemincrypto-5hwa5dthla-an.a.run.app"),
		isInGCP:                         mustAtob(getEnv("IS_IN_GCP", "false")),
		idToken:                         getEnv("ID_TOKEN", "NOT_USED_IN_GCP"),
		logger:                          logrus.StandardLogger(),
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
