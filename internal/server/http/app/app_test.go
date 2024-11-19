package httpapp

import (
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

func Test_printInfo(t *testing.T) {
	err := printInfo()
	assert.NoError(t, err)
}
