package apple

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) FetchTrack(_ context.Context, id, storefront string) (*Entity, error) {
	args := m.Called(id, storefront)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Entity), args.Error(1)
}

func (m *ClientMock) SearchTrack(_ context.Context, artist, title string) (*Entity, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Entity), args.Error(1)
}

func (m *ClientMock) FetchAlbum(_ context.Context, id, storefront string) (*Entity, error) {
	args := m.Called(id, storefront)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Entity), args.Error(1)
}

func (m *ClientMock) SearchAlbum(_ context.Context, artist, title string) (*Entity, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Entity), args.Error(1)
}
