package config

import (
	"github.com/sirupsen/logrus"
)

type Config interface {
	Setting() Setting

	ServiceName() string
	HTTPServerPort() int
	ENV() string
	GracefulShutdownTimeoutMs() int
	ShouldProfile() bool
	ShouldTrace() bool
	ShouldUseGRPCClientTLS() bool
	CACertFile() string
	BaemincryptoGRPCServiceEndpoint() string
	BaemincryptoGRPCServiceURL() string
	UserGRPCServiceEndpoint() string
	UserGRPCServiceURL() string
	IsInGCP() bool
	IDToken() string
	Logger() *logrus.Logger
}

type DefaultConfig struct {
	Config

	setting Setting
	logger  *logrus.Logger
}

func NewConfig(setting Setting, logger *logrus.Logger) Config {
	return &DefaultConfig{
		setting: setting,
		logger:  logger,
	}
}

func (c DefaultConfig) Setting() Setting {
	return c.setting
}

func MockLogger() *logrus.Logger {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	return logrus.StandardLogger()
}

func MockConfig() Config {
	return NewConfig(MockSetting(), MockLogger())
}

func (c DefaultConfig) ServiceName() string {
	return c.Setting().serviceName
}

func (c DefaultConfig) HTTPServerPort() int {
	return c.Setting().httpServerPort
}

func (c DefaultConfig) ENV() string {
	return c.Setting().env
}

func (c DefaultConfig) GracefulShutdownTimeoutMs() int {
	return c.Setting().gracefulShutdownTimeoutMs
}

func (c DefaultConfig) ShouldProfile() bool {
	return c.Setting().shouldProfile
}

func (c DefaultConfig) ShouldTrace() bool {
	return c.Setting().shouldTrace
}

func (c DefaultConfig) ShouldUseGRPCClientTLS() bool {
	return c.Setting().shouldUseGRPCClientTLS
}

func (c DefaultConfig) CACertFile() string {
	return c.Setting().caCertFile
}

func (c DefaultConfig) BaemincryptoGRPCServiceEndpoint() string {
	return c.Setting().baemincryptoGRPCServiceEndpoint
}

func (c DefaultConfig) BaemincryptoGRPCServiceURL() string {
	return c.Setting().baemincryptoGRPCServiceURL
}

func (c DefaultConfig) UserGRPCServiceEndpoint() string {
	return c.Setting().userGRPCServiceEndpoint
}

func (c DefaultConfig) UserGRPCServiceURL() string {
	return c.Setting().userGRPCServiceURL
}

func (c DefaultConfig) IsInGCP() bool {
	return c.Setting().isInGCP
}

func (c DefaultConfig) IDToken() string {
	return c.Setting().idToken
}

func (c DefaultConfig) Logger() *logrus.Logger {
	return c.logger
}
