package yandex

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) FetchTrack(_ context.Context, id string) (*Track, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Track), args.Error(1)
}

func (m *ClientMock) SearchTrack(_ context.Context, query string) (*Track, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Track), args.Error(1)
}

func (m *ClientMock) FetchAlbum(_ context.Context, id string) (*Album, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Album), args.Error(1)
}

func (m *ClientMock) SearchAlbum(_ context.Context, query string) (*Album, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Album), args.Error(1)
}
