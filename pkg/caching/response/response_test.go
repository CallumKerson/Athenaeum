package response

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponseToBytes(t *testing.T) {
	testResponse := &Response{
		Value:      nil,
		Expiration: time.Time{},
		LastAccess: time.Time{},
	}

	assert.NotEmpty(t, testResponse.Bytes())
}

func TestBytesToResponse(t *testing.T) {
	testResponse := Response{
		Value:      []byte("value 1"),
		Expiration: time.Time{},
		LastAccess: time.Time{},
	}

	got := BytesToResponse(testResponse.Bytes())
	assert.Equal(t, "value 1", string(got.Value))
}
