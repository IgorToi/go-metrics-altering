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

func Test_isIPInTrustedSubnet(t *testing.T) {
	r, err := IsIPInTrustedSubnet("192.168.1.10", "192.168.1.0/24")
	assert.True(t, r)
	assert.NoError(t, err)

	r, err = IsIPInTrustedSubnet("192.168.2", "192.168.1.0/24")
	assert.False(t, r)
	assert.Error(t, err)

	r, err = IsIPInTrustedSubnet("192.168.2", "192.168.1.0/A")
	assert.False(t, r)
	assert.Error(t, err)
}
