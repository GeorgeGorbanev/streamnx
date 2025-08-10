package youtube

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) GetVideo(_ context.Context, id string) (*Video, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Video), args.Error(1)
}

func (m *ClientMock) SearchVideo(_ context.Context, term string) (*SearchResponse, error) {
	args := m.Called(term)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SearchResponse), args.Error(1)
}

func (m *ClientMock) GetPlaylist(_ context.Context, id string) (*Playlist, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Playlist), args.Error(1)
}

func (m *ClientMock) SearchPlaylist(_ context.Context, term string) (*SearchResponse, error) {
	args := m.Called(term)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SearchResponse), args.Error(1)
}

func (m *ClientMock) GetPlaylistItems(_ context.Context, id string) ([]Video, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Video), args.Error(1)
}
