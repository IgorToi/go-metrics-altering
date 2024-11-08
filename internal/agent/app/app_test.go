package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_printInfo(t *testing.T) {
	err := printInfo()
	assert.NoError(t, err)
}
