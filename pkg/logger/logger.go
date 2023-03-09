package logger

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var (
	logger Logger

	once = &sync.Once{}
)

func GetGlobalLogger() Logger {
	once.Do(func() {
		logger = newLogger()
	})
	return logger
}

type Logger struct {
	zerolog.Logger
}

func newLogger() Logger {
	return Logger{
		zerolog.New(os.Stdout),
	}
}

func (l Logger) WithModule(name string) Logger {
	return Logger{l.Logger.With().Str("module", name).Logger()}
}
