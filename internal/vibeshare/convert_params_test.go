package vibeshare

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/stretchr/testify/require"
)

func TestConvertParamsMarshal(t *testing.T) {
	params := convertParams{
		ID:     "123",
		Source: music.Spotify,
		Target: music.Yandex,
	}
	require.Equal(t, []string{"spotify", "123", "yandex"}, params.marshal())
}

func TestConvertParamsUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    convertParams
		wantErr bool
	}{
		{
			name:  "when params are valid",
			input: []string{"spotify", "123", "yandex"},
			want: convertParams{
				ID:     "123",
				Source: music.Spotify,
				Target: music.Yandex,
			},
			wantErr: false,
		},
		{
			name:    "when invalid length",
			input:   []string{"spotify", "123"},
			wantErr: true,
		},
		{
			name:    "when invalid source provider",
			input:   []string{"invalid", "123", "yandex"},
			wantErr: true,
		},
		{
			name:    "when invalid target provider",
			input:   []string{"spotify", "123", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := convertParams{}
			err := params.unmarshal(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, params)
			}
		})
	}
}
