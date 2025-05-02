package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func New(env string) *Logger {
	var logger = logrus.New()

	if env == "dev" {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: false,
			ForceColors:   true,
		})
		logger.SetOutput(os.Stdout)

	} else if env == "prod" {
		logger.SetLevel(logrus.InfoLevel)
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			PrettyPrint:     true,
			DataKey:         "data",
		})

		file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		if err != nil {
			logger.Error("failed to log to file, using default stderr")
		} else {
			logger.SetOutput(file)
		}

	}

	return &Logger{log: logger}
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.log.WithFields(fields)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}
