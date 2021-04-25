package logging

// Logger implementation is responsible for providing structured and levels
// logging functions.
type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})

	// WithFields should return a logger which is annotated with the given
	// fields. These fields should be added to every logging call on the
	// returned logger.
	WithFields(m map[string]interface{}) Logger
}
