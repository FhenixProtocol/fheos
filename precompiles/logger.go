package precompiles

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func initLogger() {

}

func logLevelFromString(logLevel string, defaultLevel logrus.Level) logrus.Level {
	logLevelUpper := strings.ToUpper(logLevel)
	switch logLevelUpper {
	case "ERROR":
		return logrus.ErrorLevel
	case "WARN":
		return logrus.WarnLevel
	case "INFO":
		return logrus.InfoLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "TRACE":
		return logrus.TraceLevel
	default:
		return defaultLevel
	}
}

func getLogLevel(defaultLevel logrus.Level) logrus.Level {
	LogLevelEnvVarName := "LOG_LEVEL"
	logLevelEnvVar := os.Getenv(LogLevelEnvVarName)
	logLevel := logLevelFromString(logLevelEnvVar, defaultLevel)

	if logLevel > defaultLevel {
		return defaultLevel
	}

	return logLevel
}

func getDefaultLogLevel() logrus.Level {
	// TODO: read from config when it is implemented
	return logrus.InfoLevel
}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	logger.SetLevel(getLogLevel(getDefaultLogLevel()))

	return logger
}
