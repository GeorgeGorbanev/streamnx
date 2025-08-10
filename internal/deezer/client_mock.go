package deezer

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

func (m *ClientMock) SearchTrack(_ context.Context, artist, title string) (*Track, error) {
	args := m.Called(artist, title)
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

func (m *ClientMock) SearchAlbum(_ context.Context, artist, title string) (*Album, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Album), args.Error(1)
}

func (m *ClientMock) FollowCloak(_ context.Context, cloakCode string) (string, error) {
	args := m.Called(cloakCode)
	return args.String(0), args.Error(1)
}
