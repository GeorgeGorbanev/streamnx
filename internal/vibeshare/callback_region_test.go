package vibeshare

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/yandex"

	"github.com/stretchr/testify/require"
)

func TestRegionParamsMarshal(t *testing.T) {
	params := regionParams{
		EntityID: "123",
		Region:   yandex.RegionRussia,
	}
	require.Equal(t, []string{"123", "ru"}, params.marshal())
}

func TestRegionParamsUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    regionParams
		wantErr bool
	}{
		{
			name:  "when params are valid",
			input: []string{"123", "ru"},
			want: regionParams{
				EntityID: "123",
				Region:   yandex.RegionRussia,
			},
			wantErr: false,
		},
		{
			name:    "when invalid length",
			input:   []string{"123"},
			wantErr: true,
		},
		{
			name:    "when invalid locale",
			input:   []string{"123", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := regionParams{}
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
