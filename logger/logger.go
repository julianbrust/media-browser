package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

func Init(level *string) *logrus.Logger {
	file := "/tmp/media-browser.log"
	log = logrus.New()

	_ = os.Remove(file)

	logFile, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)
	if *level != "" {
		log.SetLevel(GetLogLevel(*level))
	}

	return log
}

// GetLogLevel provides the logrus log level based on the config.
func GetLogLevel(confLevel string) logrus.Level {
	switch confLevel {
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}
