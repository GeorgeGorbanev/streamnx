package yandex

import (
	"slices"
	"strings"
)

const noRegionDomainZone = "com"

var Regions = []string{"by", "kz", "ru", "uz"}

func allDomainZonesRe() string {
	return strings.Join(allDomainZones(), "|")
}

func allDomainZones() []string {
	return append(slices.Clone(Regions), noRegionDomainZone)
}
