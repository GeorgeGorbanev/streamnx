package streamnx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/internal/translator"
)

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
		mockAdapter func(m *adapterMock)
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
			mockAdapter: func(m *adapterMock) {
				m.
					On("FetchTrack", "1").
					Return(&Entity{ID: "1"}, nil).
					Once()
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
			mockAdapter: func(m *adapterMock) {
				m.
					On("FetchAlbum", "1").
					Return(&Entity{ID: "1"}, nil).
					Once()
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
			mockAdapter: func(m *adapterMock) {
				m.
					On("FetchTrack", "1").
					Return(nil, EntityNotFoundError).
					Once()
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
			mockAdapter: func(m *adapterMock) {
				m.
					On("FetchAlbum", "1").
					Return(nil, EntityNotFoundError).
					Once()
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
			mockAdapter: func(m *adapterMock) {},
			wantErr:     InvalidProviderError,
		},
		{
			name: "invalid entity type",
			args: args{
				p:  sampleProvider,
				et: EntityType("invalid"),
				id: "1",
			},
			mockAdapter: func(m *adapterMock) {},
			wantErr:     InvalidEntityTypeError,
		},
		{
			name: "cloak entity unsupported",
			args: args{
				p:  sampleProvider,
				et: Cloak,
				id: "test-cloak-code",
			},
			mockAdapter: func(m *adapterMock) {
				m.
					On("FetchCloak", "test-cloak-code").
					Return(nil, UnsupportedEntityTypeError).
					Once()
			},
			wantErr: UnsupportedEntityTypeError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			am := &adapterMock{}
			tt.mockAdapter(am)

			translatorMock := &translator.Mock{}

			registry, err := NewRegistry(
				ctx,
				Credentials{},
				WithTranslator(translatorMock),
				WithProviderAdapter(sampleProvider, am),
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

			am.AssertExpectations(t)
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
		mockAdapter func(m *adapterMock)
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
			mockAdapter: func(m *adapterMock) {
				m.
					On("SearchTrack", "artist", "name").
					Return(&Entity{ID: "1"}, nil).
					Once()
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
			mockAdapter: func(m *adapterMock) {
				m.
					On("SearchAlbum", "artist", "name").
					Return(&Entity{ID: "1"}, nil).
					Once()
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
			mockAdapter: func(m *adapterMock) {
				m.
					On("SearchTrack", "artist", "name").
					Return(nil, EntityNotFoundError).
					Once()
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
			mockAdapter: func(m *adapterMock) {
				m.
					On("SearchAlbum", "artist", "name").
					Return(nil, EntityNotFoundError).
					Once()
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
			mockAdapter: func(m *adapterMock) {},
			wantErr:     InvalidProviderError,
		},
		{
			name: "invalid entity type",
			args: args{
				p:      sampleProvider,
				et:     EntityType("invalid"),
				artist: "artist",
				name:   "name",
			},
			mockAdapter: func(m *adapterMock) {},
			wantErr:     InvalidEntityTypeError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			am := &adapterMock{}
			tt.mockAdapter(am)

			translatorMock := &translator.Mock{}

			registry, err := NewRegistry(
				ctx,
				Credentials{},
				WithTranslator(translatorMock),
				WithProviderAdapter(sampleProvider, am),
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

			am.AssertExpectations(t)
		})
	}
}

func TestRegistry_FetchCloak(t *testing.T) {
	ctx := context.Background()

	t.Run("unsupported provider", func(t *testing.T) {
		am := &adapterMock{}
		am.
			On("FetchCloak", "test-cloak-code").
			Return(nil, UnsupportedEntityTypeError).
			Once()

		translatorMock := &translator.Mock{}

		registry, err := NewRegistry(
			ctx,
			Credentials{},
			WithTranslator(translatorMock),
			WithProviderAdapter(Apple, am),
			WithProviderAdapter(Spotify, &adapterMock{}),
			WithProviderAdapter(Yandex, &adapterMock{}),
			WithProviderAdapter(Youtube, &adapterMock{}),
		)
		require.NoError(t, err)

		result, err := registry.FetchCloak(ctx, Apple, "test-cloak-code")
		require.ErrorIs(t, err, UnsupportedEntityTypeError)
		require.Nil(t, result)

		am.AssertExpectations(t)
	})
}
