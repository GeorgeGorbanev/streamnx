package soundcloud

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
			trackURL: "https://soundcloud.com/forss/flickermood",
			want:     CompositeKey{ID: "flickermood", UserSlug: "forss"},
		},
		{
			name:     "invalid track URL",
			trackURL: "https://soundcloud.com/forss/sets/soulhack",
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

func TestCompositeKey_ParseFromAlbumURL(t *testing.T) {
	tests := []struct {
		name     string
		albumURL string
		want     CompositeKey
		wantErr  error
	}{
		{
			name:     "valid album URL",
			albumURL: "https://soundcloud.com/forss/sets/soulhack",
			want:     CompositeKey{ID: "soulhack", UserSlug: "forss"},
		},
		{
			name:     "invalid album URL",
			albumURL: "https://soundcloud.com/forss/flickermood",
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

func TestCompositeKey_Marshal(t *testing.T) {
	ck := CompositeKey{
		ID:       "flickermood",
		UserSlug: "forss",
	}

	result := ck.Marshal()

	require.Equal(t, "forss:flickermood", result)
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
			name:         "valid composite key",
			compositeKey: CompositeKey{},
			input:        "forss:flickermood",
			wantResult: CompositeKey{
				ID:       "flickermood",
				UserSlug: "forss",
			},
		},
		{
			name:         "invalid composite key",
			compositeKey: CompositeKey{},
			input:        "forssflickermood",
			wantResult:   CompositeKey{},
			wantErr:      errors.New("invalid composite key: forssflickermood"),
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
