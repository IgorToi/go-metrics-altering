package logging

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f, _ := prepareFields(fields...)
		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)
		checkLevel(logger, lvl, msg)
	})
}

func checkLevel(logger *zap.Logger, lvl logging.Level, msg string) error {
	switch lvl {
	case logging.LevelDebug:
		logger.Debug(msg)
	case logging.LevelInfo:
		logger.Info(msg)
	case logging.LevelWarn:
		logger.Warn(msg)
	case logging.LevelError:
		logger.Error(msg)
	default:
		return errors.New("unknown level")
	}
	return nil
}

func prepareFields(fields ...any) ([]zapcore.Field, error) {
	f := make([]zap.Field, 0, len(fields)/2)

	for i := 0; i < len(fields); i += 2 {
		key := fields[i]
		value := fields[i+1]

		switch v := value.(type) {
		case string:
			f = append(f, zap.String(key.(string), v))
		case int:
			f = append(f, zap.Int(key.(string), v))
		case bool:
			f = append(f, zap.Bool(key.(string), v))
		default:
			f = append(f, zap.Any(key.(string), v))
		}
	}
	return f, nil
}
