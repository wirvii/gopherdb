package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// zapLogger is a logger implementation using zap.
type zapLogger struct {
	logger *zap.SugaredLogger
}

// newZapLogger creates a new instance of ZapLogger.
func newZapLogger() *zapLogger {
	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeTime = timeEncoder
	conf.EncoderConfig.CallerKey = "invoker"
	conf.EncoderConfig.MessageKey = "message"
	conf.EncoderConfig.LevelKey = "level"
	conf.EncoderConfig.TimeKey = "timestamp"

	l, _ := conf.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.PanicLevel),
	)

	return &zapLogger{
		logger: l.Sugar(),
	}
}

// Infof logs an info message with format.
func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Errorf logs an error message with format.
func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Warnf logs a warning message with format.
func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Debugf logs a debug message with format.
func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Debugw logs a debug message with key-value pairs.
func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

// Infow logs an info message with key-value pairs.
func (l *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

// Warnw logs a warning message with key-value pairs.
func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

// Errorw logs an error message with key-value pairs.
func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}
