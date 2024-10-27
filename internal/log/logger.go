package log

// Logger is an interface that is used to log messages.
type Logger interface {
	// Infof logs an info message with format.
	Infof(format string, args ...interface{})
	// Errorf logs an error message with format.
	Errorf(format string, args ...interface{})
	// Warnf logs a warning message with format.
	Warnf(format string, args ...interface{})
	// Debugf logs a debug message with format.
	Debugf(format string, args ...interface{})

	// Debugw logs a debug message with key-value pairs.
	Debugw(msg string, keysAndValues ...interface{})
	// Infow logs an info message with key-value pairs.
	Infow(msg string, keysAndValues ...interface{})
	// Warnw logs a warning message with key-value pairs.
	Warnw(msg string, keysAndValues ...interface{})
	// Errorw logs an error message with key-value pairs.
	Errorw(msg string, keysAndValues ...interface{})
}

// NewLogger creates a new instance of Logger.
func NewLogger() Logger {
	return newZapLogger()
}
