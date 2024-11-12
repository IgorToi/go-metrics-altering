// Package auth provides middleware for basic authorization.
package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidMAC(t *testing.T) {
	_, err := ValidMAC([]byte("A"), []byte("A"), []byte("key"))
	assert.NoError(t, err)
}
