package stencil

// Logger interface used to get logging from stencil internals.
type Logger interface {
	Info(string)
	Error(string)
}

type wrappedLogger struct {
	l Logger
}

func (w wrappedLogger) Info(msg string) {
	if w.l != nil {
		w.l.Info(msg)
	}
}

func (w wrappedLogger) Error(msg string) {
	if w.l != nil {
		w.l.Info(msg)
	}
}

func wrapLogger(l Logger) Logger {
	return wrappedLogger{l}
}
