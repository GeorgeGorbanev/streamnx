package deezer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectTrackID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "https://www.deezer.com/us/track/441481432",
			want:  "441481432",
		},
		{
			input: "https://deezer.com/us/track/441481432",
			want:  "441481432",
		},
		{
			input: "https://deezer.com/track/441481432",
			want:  "441481432",
		},
		{
			input: "https://www.deezer.com/track/3129746?host=6024526261",
			want:  "3129746",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := DetectTrackID(tt.input)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestDetectAlbumID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "https://www.deezer.com/us/album/53564982",
			want:  "53564982",
		},
		{
			input: "https://deezer.com/us/album/53564982",
			want:  "53564982",
		},
		{
			input: "https://deezer.com/album/53564982",
			want:  "53564982",
		},
		{
			input: "https://www.deezer.com/album/3129746?host=6024526261",
			want:  "3129746",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := DetectAlbumID(tt.input)
			require.Equal(t, tt.want, result)
		})
	}

}

func TestDetectCloakID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "https://link.deezer.com/s/30FbcgrctxIrNQImDnVEZ",
			want:  "30FbcgrctxIrNQImDnVEZ",
		},
		{
			input: "https://link.deezer.com/s/jHWZaLoRJutY3TiVA",
			want:  "jHWZaLoRJutY3TiVA",
		},
		{
			input: "not a deezer link",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := DetectCloakID(tt.input)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestDetectCloakDest(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr string
	}{
		{
			input: "https://link.deezer.com/?awf=https%3A%2F%2Fwww.deezer.com%2Ftrack%2F3129746%3Fhost%3D6024526261%26utm_campaign%3Dclipboard-generic%26utm_source%3Duser_sharing%26utm_content%3Dtrack-3129746%26deferredFl%3D1%26universal_link%3D1&dest=https%3A%2F%2Fwww.deezer.com%2Ftrack%2F3129746%3Fhost%3D6024526261%26utm_campaign%3Dclipboard-generic%26utm_source%3Duser_sharing%26utm_content%3Dtrack-3129746%26deferredFl%3D1%26universal_link%3D1&gwf=https%3A%2F%2Fwww.deezer.com%2Ftrack%2F3129746%3Fhost%3D6024526261%26utm_campaign%3Dclipboard-generic%26utm_source%3Duser_sharing%26utm_content%3Dtrack-3129746%26deferredFl%3D1%26universal_link%3D1&iwf=https%3A%2F%2Fwww.deezer.com%2Ftrack%2F3129746%3Fhost%3D6024526261%26utm_campaign%3Dclipboard-generic%26utm_source%3Duser_sharing%26utm_content%3Dtrack-3129746%26deferredFl%3D1%26universal_link%3D1",
			want:  "https://www.deezer.com/track/3129746?host=6024526261&utm_campaign=clipboard-generic&utm_source=user_sharing&utm_content=track-3129746&deferredFl=1&universal_link=1",
		},
		{
			input:   "not a url",
			wantErr: "failed to find cloak 'dest' param in url: not a url",
		},
		{
			input:   "https://link.deezer.com?no_dest=true",
			wantErr: "failed to find cloak 'dest' param in url: https://link.deezer.com?no_dest=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := DetectCloakDest(tt.input)

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}
