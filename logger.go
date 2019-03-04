package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/gemnasium/logrus-graylog-hook.v2"
)

// New creates new logger
func NewLogger(logstashServer string) *logrus.Logger {
	logger := logrus.New()

	logger.SetLevel(logrus.TraceLevel)

	hook := graylog.NewGraylogHook(logstashServer, map[string]interface{}{})
	logger.AddHook(hook)

	return logger
}
