package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var log *logrus.Logger

// LoggerConfig specifies the service logger configuration
type LoggerConfig struct {
	// Exported
	Sanitize bool
	Backend  loggerBackend
	Fields   logFields
}

// type Logger log.Logger

// defaultLoggerConfig sets the default configuration for a service logger
var defaultLoggerConfig = &LoggerConfig{
	Backend: LogrusPrefixed,
	Fields: logFields{
		"service": defaultServiceName,
	},
}

// alias type logrus.Fields with logFields
type logFields = logrus.Fields

// define a valid list of loggers
type loggerType string

const (
	LogNone        loggerBackend = ""
	LogDefault     loggerBackend = "log"
	Logrus         loggerBackend = "logrus"
	LogrusPrefixed loggerBackend = "logrus-prefixed"
)

// String method returns the Logger type into a string type.
func (l *loggerType) String() string {
	return fmt.Sprintf("%v", l)
}

// newLoggerWithConfig instanciates a new logger with a LoggerConfig struct
func newLoggerWithConfig(lc *LoggerConfig) *logrus.Logger {
	fmt.Println("logger backend=", lc.Backend.String())

	var logger *logrus.Logger

	switch lc.Backend {
	case Logrus:
		logger = logrus.New()
		logger.Level = logrus.DebugLevel
		logger.Debugln("Logger started at", time.Now().Format("2015-06-16-0431 UTC"))

	case LogrusPrefixed:
		logger = logrus.New()
		logger.Formatter = new(prefixed.TextFormatter)
		logger.Level = logrus.DebugLevel
		logger.WithFields(lc.Fields).Debugln("Logger started at", time.Now().Format("2015-06-16-0431 UTC"))

	case LogNone:
		fallthrough

	default:
		fmt.Println("no logger instanciated...")
	}

	return logger

}
