package youtubemusic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAlbumID(t *testing.T) {
	tests := []struct {
		name    string
		idType  albumIDType
		rawID   string
		want    string
		wantErr bool
	}{
		{name: "browse", idType: albumIDTypeBrowse, rawID: "MPREb_sample-123", want: "b:MPREb_sample-123"},
		{name: "playlist", idType: albumIDTypePlaylist, rawID: "OLAK5uy_sample-123", want: "p:OLAK5uy_sample-123"},
		{name: "invalid type", idType: "browse", rawID: "MPREb_sample-123", wantErr: true},
		{name: "empty id", idType: albumIDTypeBrowse, wantErr: true},
		{name: "delimiter in id", idType: albumIDTypeBrowse, rawID: "MPRE:sample", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newAlbumID(tt.idType, tt.rawID)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
			gotType, gotRawID, err := parseAlbumID(got)
			require.NoError(t, err)
			require.Equal(t, tt.idType, gotType)
			require.Equal(t, tt.rawID, gotRawID)
		})
	}
}

func TestParseAlbumIDRejectsInvalidDump(t *testing.T) {
	for _, id := range []string{"", "MPREb_sample", "b:", "browse:MPREb_sample", "b:MPRE:sample"} {
		t.Run(id, func(t *testing.T) {
			_, _, err := parseAlbumID(id)
			require.Error(t, err)
		})
	}
}
