package soundcloud

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) FetchTrack(_ context.Context, userSlug, trackSlug string) (*Track, error) {
	args := m.Called(userSlug, trackSlug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Track), args.Error(1)
}

func (m *ClientMock) FetchAlbum(_ context.Context, userSlug, setSlug string) (*Album, error) {
	args := m.Called(userSlug, setSlug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Album), args.Error(1)
}

func (m *ClientMock) SearchTrack(_ context.Context, artist, title string) (*Track, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Track), args.Error(1)
}

func (m *ClientMock) SearchAlbum(_ context.Context, artist, title string) (*Album, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Album), args.Error(1)
}
