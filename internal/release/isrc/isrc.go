package isrc

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`^[A-Z]{2}(?:[A-Z0-9]{3}[0-9]{7}|-[A-Z0-9]{3}-[0-9]{2}-[0-9]{5})$`)

func Normalize(raw string) (string, bool) {
	normalized := strings.ToUpper(strings.TrimSpace(raw))
	if !re.MatchString(normalized) {
		return "", false
	}
	return strings.ReplaceAll(normalized, "-", ""), true
}
