package agent

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_retryURL(t *testing.T) {
	var r http.Request
	var client http.Client

	err := retryURL(client, &r)
	assert.NotNil(t, err)
}
