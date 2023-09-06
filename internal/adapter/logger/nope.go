package logger

import "go.uber.org/zap"

type NilLogger struct {
}

func NewNilLogger() *NilLogger {
	return &NilLogger{}
}

func (l *NilLogger) Debug(msg string, fields ...zap.Field) {
}

func (l *NilLogger) Info(msg string, fields ...zap.Field) {
}

func (l *NilLogger) Error(msg string, fields ...zap.Field) {
}
