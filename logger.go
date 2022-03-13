package oauth2

import "log"

type Logger interface {
	Info(format string, v ...interface{})
}

type DefaultLogger struct {
}

func (l *DefaultLogger) Info(format string, v ...interface{}) {
	log.Println(format, v)
}
