package logging

import (
	"context"
	"reflect"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	l "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type Logger interface {
	Log(ctx context.Context, level l.Level, msg string, fields ...any)
}

func TestInterceptorLogger(t *testing.T) {
	logger := zap.NewExample()
	ls := InterceptorLogger(logger)

	ts := reflect.TypeOf(ls)
	i := reflect.TypeOf((*Logger)(nil)).Elem()
	assert.True(t, ts.Implements(i))
}

func Test_checkLevel(t *testing.T) {
	logger := zap.NewExample()
	err := checkLevel(logger, logging.LevelDebug, "test")
	assert.NoError(t, err)

	err = checkLevel(logger, logging.LevelInfo, "test")
	assert.NoError(t, err)

	err = checkLevel(logger, logging.LevelWarn, "test")
	assert.NoError(t, err)

	err = checkLevel(logger, logging.LevelError, "test")
	assert.NoError(t, err)

	err = checkLevel(logger, 45, "test")
	assert.NotNil(t, err)
}

func Test_prepareFields(t *testing.T) {
	_, err := prepareFields("1", "2")
	assert.NoError(t, err)

	_, err = prepareFields("1", 2)
	assert.NoError(t, err)

	_, err = prepareFields("true", false)
	assert.NoError(t, err)

	_, err = prepareFields("1", "2")
	assert.NoError(t, err)
}
