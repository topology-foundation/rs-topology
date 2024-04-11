// The code inside the `log` directory is a modified version of the [`go-ethereum` log](https://github.com/ethereum/go-ethereum/tree/bd91810462187086b2715fd343aa427e181d89a2/log)

package log

import (
	"log/slog"
	"os"
	"sync/atomic"

	"github.com/topology-gg/gram/config"
)

var root atomic.Value

func init() {
	root.Store(&logger{slog.New(EmptyHandler())})
}

// SetDefault sets the default global logger
func SetDefault(config *config.LogConfig) {
	level := StringLevel(config.LogLevel)
	l := NewLogger(NewHandler(os.Stdout, level, true))

	root.Store(l)
	if lg, ok := l.(*logger); ok {
		slog.SetDefault(lg.inner)
	}
}

// Root returns the root logger
func Root() Logger {
	return root.Load().(Logger)
}

// Debug is a convenient alias for Root().Debug
func Debug(msg string, ctx ...interface{}) {
	Root().Write(slog.LevelDebug, msg, ctx...)
}

// Info is a convenient alias for Root().Info
func Info(msg string, ctx ...interface{}) {
	Root().Write(slog.LevelInfo, msg, ctx...)
}

// Warn is a convenient alias for Root().Warn
func Warn(msg string, ctx ...interface{}) {
	Root().Write(slog.LevelWarn, msg, ctx...)
}

// Error is a convenient alias for Root().Error
func Error(msg string, ctx ...interface{}) {
	Root().Write(slog.LevelError, msg, ctx...)
}
