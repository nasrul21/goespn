package goespn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	client := &mockClient{}

	t.Run("Valid", func(t *testing.T) {
		req, err := client.NewRequest("some-method", "some-full-path")
		assert.Nil(t, err)
		assert.NotNil(t, req)
	})

	t.Run("MissingMethod", func(t *testing.T) {
		req, err := client.NewRequest("", "some-full-path")
		assert.Nil(t, req)
		assert.NotNil(t, err)
	})

	t.Run("MissingFullpath", func(t *testing.T) {
		req, err := client.NewRequest("some-path", "")
		assert.Nil(t, req)
		assert.NotNil(t, err)
	})

}
