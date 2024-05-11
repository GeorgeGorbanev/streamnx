package yandex

import (
	"slices"
	"strings"
)

const defaultDomainZone = "com"

type Region struct {
	DomainZone string
	Label      string
}

var (
	RegionBelarus = &Region{
		DomainZone: "by",
		Label:      "🇧🇾 Belarus",
	}
	RegionRussia = &Region{
		DomainZone: "ru",
		Label:      "🇷🇺 Russia",
	}
	RegionKazakhstan = &Region{
		DomainZone: "kz",
		Label:      "🇰🇿 Kazakhstan",
	}
	RegionUzbekistan = &Region{
		DomainZone: "uz",
		Label:      "🇺🇿 Uzbekistan",
	}
	Regions = []*Region{
		RegionBelarus,
		RegionKazakhstan,
		RegionRussia,
		RegionUzbekistan,
	}
)

func FindRegionByDomainZone(region string) *Region {
	for _, r := range Regions {
		if r.DomainZone == region {
			return r
		}
	}
	return nil
}

func (r *Region) IsValid() bool {
	return slices.Contains(Regions, r)
}

func (r *Region) LocalizeLink(link string) string {
	return strings.Replace(link, defaultDomainZone, r.DomainZone, 1)
}

func allDomainZonesRe() string {
	return strings.Join(allDomainZones(), "|")
}

func allDomainZones() []string {
	result := make([]string, 0, len(Regions)+1)
	result = append(result, defaultDomainZone)
	for _, r := range Regions {
		result = append(result, r.DomainZone)
	}
	return result
}
