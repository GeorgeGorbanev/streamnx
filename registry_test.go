package streaminx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type adapterMock struct {
	fetchTrack  map[string]*Entity
	searchTrack map[string]map[string]*Entity
	fetchAlbum  map[string]*Entity
	searchAlbum map[string]map[string]*Entity
}

func (a *adapterMock) FetchTrack(_ context.Context, id string) (*Entity, error) {
	entity, ok := a.fetchTrack[id]
	if !ok {
		return nil, EntityNotFoundError
	}
	return entity, nil
}

func (a *adapterMock) SearchTrack(_ context.Context, artistName, trackName string) (*Entity, error) {
	if tracks, ok := a.searchTrack[artistName]; ok {
		track, ok := tracks[trackName]
		if !ok {
			return nil, EntityNotFoundError
		}
		return track, nil
	}
	return nil, EntityNotFoundError
}

func (a *adapterMock) FetchAlbum(_ context.Context, id string) (*Entity, error) {
	entity, ok := a.fetchAlbum[id]
	if !ok {
		return nil, EntityNotFoundError
	}
	return entity, nil
}

func (a *adapterMock) SearchAlbum(_ context.Context, artistName, albumName string) (*Entity, error) {
	if albums, ok := a.searchAlbum[artistName]; ok {
		album, ok := albums[albumName]
		if !ok {
			return nil, EntityNotFoundError
		}
		return album, nil

	}
	return nil, EntityNotFoundError
}

func TestRegistry_Fetch(t *testing.T) {
	sampleProvider := Apple

	type args struct {
		p  *Provider
		et EntityType
		id string
	}

	tests := []struct {
		name        string
		args        args
		adapterMock adapterMock
		want        *Entity
		wantErr     error
	}{
		{
			name: "track found",
			args: args{
				p:  sampleProvider,
				et: Track,
				id: "1",
			},
			adapterMock: adapterMock{
				fetchTrack: map[string]*Entity{
					"1": {ID: "1"},
				},
			},
			want: &Entity{ID: "1"},
		},
		{
			name: "album found",
			args: args{
				p:  sampleProvider,
				et: Album,
				id: "1",
			},
			adapterMock: adapterMock{
				fetchAlbum: map[string]*Entity{
					"1": {ID: "1"},
				},
			},
			want: &Entity{ID: "1"},
		},
		{
			name: "track not found",
			args: args{
				p:  sampleProvider,
				et: Track,
				id: "1",
			},
			want:    nil,
			wantErr: EntityNotFoundError,
		},
		{
			name: "album not found",
			args: args{
				p:  sampleProvider,
				et: Album,
				id: "1",
			},
			want:    nil,
			wantErr: EntityNotFoundError,
		},
		{
			name: "invalid provider",
			args: args{
				p:  &Provider{},
				et: Track,
				id: "1",
			},
			wantErr: InvalidProviderError,
		},
		{
			name: "invalid entity type",
			args: args{
				p:  sampleProvider,
				et: EntityType("invalid"),
				id: "1",
			},
			wantErr: InvalidEntityTypeError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			registry, err := NewRegistry(
				ctx,
				Credentials{},
				WithTranslator(&translatorMock{}),
				WithProviderAdapter(sampleProvider, &tt.adapterMock),
				WithProviderAdapter(Spotify, &adapterMock{}),
				WithProviderAdapter(Yandex, &adapterMock{}),
				WithProviderAdapter(Youtube, &adapterMock{}),
			)
			require.NoError(t, err)

			result, err := registry.Fetch(ctx, tt.args.p, tt.args.et, tt.args.id)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestRegistry_Search(t *testing.T) {
	sampleProvider := Apple

	type args struct {
		p      *Provider
		et     EntityType
		artist string
		name   string
	}

	tests := []struct {
		name        string
		args        args
		adapterMock adapterMock
		want        *Entity
		wantErr     error
	}{
		{
			name: "track found",
			args: args{
				p:      sampleProvider,
				et:     Track,
				artist: "artist",
				name:   "name",
			},
			adapterMock: adapterMock{
				searchTrack: map[string]map[string]*Entity{
					"artist": {
						"name": {
							ID: "1",
						},
					},
				},
			},
			want: &Entity{ID: "1"},
		},
		{
			name: "album found",
			args: args{
				p:      sampleProvider,
				et:     Album,
				artist: "artist",
				name:   "name",
			},
			adapterMock: adapterMock{
				searchAlbum: map[string]map[string]*Entity{
					"artist": {
						"name": {
							ID: "1",
						},
					},
				},
			},
			want: &Entity{ID: "1"},
		},
		{
			name: "track not found",
			args: args{
				p:      sampleProvider,
				et:     Track,
				artist: "artist",
				name:   "name",
			},
			want:    nil,
			wantErr: EntityNotFoundError,
		},
		{
			name: "album not found",
			args: args{
				p:      sampleProvider,
				et:     Album,
				artist: "artist",
				name:   "name",
			},
			want:    nil,
			wantErr: EntityNotFoundError,
		},
		{
			name: "invalid provider",
			args: args{
				p:      &Provider{},
				et:     Track,
				artist: "artist",
				name:   "name",
			},
			wantErr: InvalidProviderError,
		},
		{
			name: "invalid entity type",
			args: args{
				p:      sampleProvider,
				et:     EntityType("invalid"),
				artist: "artist",
				name:   "name",
			},
			wantErr: InvalidEntityTypeError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			registry, err := NewRegistry(
				ctx,
				Credentials{},
				WithTranslator(&translatorMock{}),
				WithProviderAdapter(sampleProvider, &tt.adapterMock),
				WithProviderAdapter(Spotify, &adapterMock{}),
				WithProviderAdapter(Yandex, &adapterMock{}),
				WithProviderAdapter(Youtube, &adapterMock{}),
			)
			require.NoError(t, err)

			result, err := registry.Search(ctx, tt.args.p, tt.args.et, tt.args.artist, tt.args.name)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}
