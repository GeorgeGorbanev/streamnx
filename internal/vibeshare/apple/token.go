package apple

import (
	"fmt"
	"regexp"
)

var (
	tokenBundleRe = regexp.MustCompile(`src="(/assets/index-[a-zA-Z0-9]+\.js)"`)
	tokenVarRe    = regexp.MustCompile(`headers\.Authorization\s*=\s*` + "`Bearer \\${([a-zA-Z0-9_]+)}`")
)

func parseBundleName(html []byte) string {
	matches := tokenBundleRe.FindSubmatch(html)
	if matches == nil || len(matches) < 2 {
		return ""
	}
	return string(matches[1])
}

func parseToken(jsBundle []byte) (string, error) {
	tokenVar := parseTokenVar(jsBundle)
	if tokenVar == "" {
		return "", fmt.Errorf("token variable not found")
	}

	tokenVarValue, err := parseVariableValue(jsBundle, tokenVar)
	if err != nil {
		return "", fmt.Errorf("failed to find token variable value: %s", err)
	}

	return tokenVarValue, err
}

func parseTokenVar(jsBundle []byte) string {
	matches := tokenVarRe.FindSubmatch(jsBundle)
	if matches == nil || len(matches) < 2 {
		return ""
	}
	return string(matches[1])
}

func parseVariableValue(jsBundle []byte, variable string) (string, error) {
	re := regexp.MustCompile(fmt.Sprintf(`\b%s\s*=\s*"([^"]+)"`, regexp.QuoteMeta(variable)))
	matches := re.FindSubmatch(jsBundle)
	if matches == nil || len(matches) < 2 {
		return "", fmt.Errorf("value of variable %s not found", variable)
	}
	return string(matches[1]), nil
}
