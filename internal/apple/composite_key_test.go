package apple

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompositeKey_ParseFromTrackURL(t *testing.T) {
	tests := []struct {
		name     string
		trackURL string
		want     CompositeKey
		wantErr  error
	}{
		{
			name:     "valid track URL",
			trackURL: "https://music.apple.com/ru/album/song-name/1234567890?i=987654321",
			want:     CompositeKey{ID: "987654321", Storefront: "ru"},
		},
		{
			name:     "invalid storefront URL",
			trackURL: "https://music.apple.com/invalid/album/song-name/1234567890?i=987654321",
			wantErr:  CompositeKeyError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &CompositeKey{}
			err := k.ParseFromTrackURL(tt.trackURL)

			if tt.wantErr != nil {
				require.ErrorAs(t, err, &tt.wantErr)
			} else {
				require.Equal(t, tt.want, *k)
				require.NoError(t, err)
			}
		})
	}
}

func TestCompositeKey_ParseFromAlbumURL(t *testing.T) {
	tests := []struct {
		name     string
		albumURL string
		want     CompositeKey
		wantErr  error
	}{
		{
			name:     "valid album URL",
			albumURL: "https://music.apple.com/us/album/album-name/123456789",
			want:     CompositeKey{ID: "123456789", Storefront: "us"},
		},
		{
			name:     "invalid storefront URL",
			albumURL: "https://music.apple.com/invalid/album/album-name/123456789",
			wantErr:  CompositeKeyError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &CompositeKey{}
			err := k.ParseFromAlbumURL(tt.albumURL)

			if tt.wantErr != nil {
				require.ErrorAs(t, err, &tt.wantErr)
			} else {
				require.Equal(t, tt.want, *k)
				require.NoError(t, err)
			}
		})
	}
}

func TestCompositeKey_Marshal(t *testing.T) {
	ck := CompositeKey{
		ID:         "123",
		Storefront: "us",
	}

	result := ck.Marshal()

	require.Equal(t, "us-123", result)
}

func TestCompositeKey_Unmarshal(t *testing.T) {
	tests := []struct {
		name         string
		compositeKey CompositeKey
		input        string
		wantResult   CompositeKey
		wantErr      error
	}{
		{
			name: "valid composite key",
			compositeKey: CompositeKey{
				ID:         "",
				Storefront: "",
			},
			input: "us-123",
			wantResult: CompositeKey{
				ID:         "123",
				Storefront: "us",
			},
			wantErr: nil,
		},
		{
			name: "invalid composite key",
			compositeKey: CompositeKey{
				ID:         "",
				Storefront: "",
			},
			input:      "us123",
			wantResult: CompositeKey{},
			wantErr:    errors.New("invalid composite key: us123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.compositeKey.Unmarshal(tt.input)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.wantResult, tt.compositeKey)
		})
	}
}
