package bandcamp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompositeKey_ParseFromAlbumURL(t *testing.T) {
	tests := []struct {
		name     string
		albumURL string
		want     CompositeKey
		wantErr  error
	}{
		{
			name:     "valid album URL",
			albumURL: "https://autechre.bandcamp.com/album/amber",
			want:     CompositeKey{ID: "amber", ArtistSlug: "autechre"},
		},
		{
			name:     "invalid album URL",
			albumURL: "https://bandcamp.com/login",
			wantErr:  errors.New("invalid album url"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &CompositeKey{}
			err := k.ParseFromAlbumURL(tt.albumURL)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.Equal(t, tt.want, *k)
				require.NoError(t, err)
			}
		})
	}
}

func TestCompositeKey_ParseFromTrackURL(t *testing.T) {
	tests := []struct {
		name     string
		trackURL string
		want     CompositeKey
		wantErr  error
	}{
		{
			name:     "valid track URL",
			trackURL: "https://autechre.bandcamp.com/track/nil",
			want:     CompositeKey{ID: "nil", ArtistSlug: "autechre"},
		},
		{
			name:     "invalid track URL",
			trackURL: "https://bandcamp.com/login",
			wantErr:  errors.New("invalid track url"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &CompositeKey{}
			err := k.ParseFromTrackURL(tt.trackURL)

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.Equal(t, tt.want, *k)
				require.NoError(t, err)
			}
		})
	}
}

func TestCompositeKey_Marshal(t *testing.T) {
	ck := CompositeKey{
		ID:         "sampleid",
		ArtistSlug: "sampleartistslug",
	}

	result := ck.Marshal()

	require.Equal(t, "sampleartistslug:sampleid", result)
}

func TestCompositeKey_Unmarshal(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantResult CompositeKey
		wantErr    error
	}{
		{
			name:  "valid composite key",
			input: "autechre:amber",
			wantResult: CompositeKey{
				ID:         "amber",
				ArtistSlug: "autechre",
			},
			wantErr: nil,
		},
		{
			name:       "invalid composite key",
			input:      "autechreamber",
			wantResult: CompositeKey{},
			wantErr:    errors.New("invalid composite key: autechreamber"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ck := CompositeKey{}
			err := ck.Unmarshal(tt.input)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.wantResult, ck)
		})
	}
}
