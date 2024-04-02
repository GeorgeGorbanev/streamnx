package apple

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseBundleName(t *testing.T) {
	html := []byte(`
		<!DOCTYPE html>
		<html lang="en">
			<head><script type="module" crossorigin="" src="/assets/index-d23b7a84.js"></script></head>
			<body>sample body</body>
	 	</html>`)

	result := parseBundleName(html)
	require.Equal(t, "/assets/index-d23b7a84.js", result)
}

func Test_parseTokenVar(t *testing.T) {
	jsBundle := []byte("sampleJs(); headers.Authorization = `Bearer ${token}`; anotherSampleJs();")
	result := parseTokenVar(jsBundle)
	require.Equal(t, "token", result)
}

func Test_parseVariableValue(t *testing.T) {
	jsBundle := []byte(`sampleJs(); var var1 = "varValue1"; var var2 = "varValue2"; var var3 = "varValue3"`)
	result, err := parseVariableValue(jsBundle, "var2")
	require.NoError(t, err)
	require.Equal(t, "varValue2", result)
}

func Test_parseToken(t *testing.T) {
	jsBundle := []byte(`token="sampleToken"; headers.Authorization = ` + "`Bearer ${token}`" + `; anotherSampleJs();`)
	result, err := parseToken(jsBundle)
	require.NoError(t, err)
	require.Equal(t, "sampleToken", result)
}
