package helpers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// SetLogLevel sets the log level for logrus by string; possible values are debug, error, fatal, panic, info, trace
func SetLogLevel(loglevelString string) {
	loglevel := GetLogLevel(loglevelString)
	logrus.SetLevel(loglevel)
	logrus.Debugf("LogLevel: %s", loglevel.String())
}

func GetLogLevel(logLevelString string) logrus.Level {
	l, err := logrus.ParseLevel(logLevelString)
	if err != nil {
		return logrus.ErrorLevel
	}
	return l
}

func SetLogFormatter(logFormatString string) {
	logrus.SetFormatter(GetLogFormatter(logFormatString))
	logrus.Debugf("LogFormat: %s", logFormatString)
}

func GetLogFormatter(logFormatString string) logrus.Formatter {
	if logFormatString == "json" {
		return &logrus.JSONFormatter{}
	}
	return &logrus.TextFormatter{
		ForceColors:               true,
		FullTimestamp:             true,
		QuoteEmptyFields:          true,
		EnvironmentOverrideColors: true,
	}
}

func LogHTTPRequest(r *http.Request) {
	logrus.Debugf("Request: %s %s %s", r.Method, r.RequestURI, r.Proto)
}
