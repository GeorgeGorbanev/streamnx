package streamnx

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type adapterMock struct {
	mock.Mock
}

func (m *adapterMock) FetchTrack(_ context.Context, id string) (*Entity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Entity), args.Error(1)
}

func (m *adapterMock) SearchTrack(_ context.Context, artist, title string) (*Entity, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Entity), args.Error(1)
}

func (m *adapterMock) FetchAlbum(_ context.Context, id string) (*Entity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Entity), args.Error(1)
}

func (m *adapterMock) SearchAlbum(_ context.Context, artist, title string) (*Entity, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Entity), args.Error(1)
}

func (m *adapterMock) FetchCloak(_ context.Context, cloakCode string) (*Entity, error) {
	args := m.Called(cloakCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Entity), args.Error(1)
}
