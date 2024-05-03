package yandex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindRegionByDomainZone(t *testing.T) {
	tests := []struct {
		code string
		want *Region
	}{
		{
			code: "ru",
			want: RegionRussia,
		},
		{
			code: "by",
			want: RegionBelarus,
		},
		{
			code: "kz",
			want: RegionKazakhstan,
		},
		{
			code: "uz",
			want: RegionUzbekistan,
		},
		{
			code: "unknown",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := FindRegionByDomainZone(tt.code)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestRegion_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		region *Region
		want   bool
	}{
		{
			name:   "valid region by",
			region: RegionBelarus,
			want:   true,
		},
		{
			name:   "valid region kz",
			region: RegionKazakhstan,
			want:   true,
		},
		{
			name:   "valid region ru",
			region: RegionRussia,
			want:   true,
		},
		{
			name:   "valid region uz",
			region: RegionUzbekistan,
			want:   true,
		},
		{
			name: "invalid region",
			region: &Region{
				DomainZone: "invalid",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.region.IsValid())
		})
	}
}

func TestRegion_LocalizeLink(t *testing.T) {
	tests := []struct {
		name   string
		region *Region
		link   string
		want   string
	}{
		{
			name:   "localize link for Belarus",
			region: RegionBelarus,
			link:   "https://example.com",
			want:   "https://example.by",
		},
		{
			name:   "localize link for Russia",
			region: RegionRussia,
			link:   "https://example.com",
			want:   "https://example.ru",
		},
		{
			name:   "localize link for Kazakhstan",
			region: RegionKazakhstan,
			link:   "https://example.com",
			want:   "https://example.kz",
		},
		{
			name:   "localize link for Uzbekistan",
			region: RegionUzbekistan,
			link:   "https://example.com",
			want:   "https://example.uz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.region.LocalizeLink(tt.link)
			require.Equal(t, tt.want, result)
		})
	}
}
