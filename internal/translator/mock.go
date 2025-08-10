package translator

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) TranslateEnToRu(_ context.Context, text string) (string, error) {
	args := m.Called(text)
	return args.String(0), args.Error(1)
}

func (m *Mock) Close() error {
	args := m.Called()
	return args.Error(0)
}
