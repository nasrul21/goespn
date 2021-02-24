package espn

import (
	"errors"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	mock.Mock
}

func (m *mockClient) NewRequest(method string, fullPath string) (*http.Request, error) {
	if method == "" || fullPath == "" {
		return nil, errors.New("invalid arguments")
	}
	return &http.Request{}, nil
}
