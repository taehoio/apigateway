package config

import (
	"github.com/sirupsen/logrus"
)

type Config interface {
	Setting() Setting
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

func (c DefaultConfig) Logger() *logrus.Logger {
	return c.logger
}
