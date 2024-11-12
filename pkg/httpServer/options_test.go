package httpserver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	_, err := Address("localhost8080")
	assert.NoError(t, err)
}

func TestReadTimeout(t *testing.T) {
	_, err := ReadTimeout(time.Second * 3)
	assert.NoError(t, err)
}

func TestWriteTimeout(t *testing.T) {
	_, err := WriteTimeout(time.Second * 3)
	assert.NoError(t, err)
}

func TestShutdownTimeout(t *testing.T) {
	_, err := ShutdownTimeout(time.Second * 3)
	assert.NoError(t, err)
}
