package logger

import (
	"context"
	"github.com/erkkipm/sso_auth/pkg/logger/handlers/slogpretty"
	"io"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// SetupLogger возвращает логгер и функцию закрытия файла лога.
func SetupLogger(env string, name string) (*slog.Logger, func()) {
	fileName := name + ".log"
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	fileHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug})

	var log *slog.Logger
	switch env {
	case envLocal:
		opts := slogpretty.PrettyHandlerOptions{SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug}}
		consoleHandler := opts.NewPrettyHandler(os.Stdout)
		log = slog.New(&multiHandler{handlers: []slog.Handler{consoleHandler, fileHandler}})
	case envDev:
		mw := io.MultiWriter(os.Stdout, logFile)
		log = slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		mw := io.MultiWriter(os.Stdout, logFile)
		log = slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		mw := io.MultiWriter(os.Stdout, logFile)
		log = slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	log.Info("Приложение запущено", slog.String("app-version", "v0.0.1-beta"))
	return log, func() { logFile.Close() }
}

// multiHandler транслирует записи всем вложенным обработчикам.
type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}

func SetupSlog() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}