package extid

import (
	"regexp"
	"strings"
)

var isrcRegexp = regexp.MustCompile(`^[A-Z]{2}(?:[A-Z0-9]{3}[0-9]{7}|-[A-Z0-9]{3}-[0-9]{2}-[0-9]{5})$`)

func NormalizeISRC(rawISRC string) (string, bool) {
	normalizedISRC := strings.ToUpper(strings.TrimSpace(rawISRC))
	if !isrcRegexp.MatchString(normalizedISRC) {
		return "", false
	}
	return strings.ReplaceAll(normalizedISRC, "-", ""), true
}
