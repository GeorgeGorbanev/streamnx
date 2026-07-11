package compositekey

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestScheme() *Scheme {
	return NewScheme(
		"-",
		Part{Name: "storefront", Re: regexp.MustCompile(`^[a-z]{2}$`)},
		Part{Name: "id", Re: regexp.MustCompile(`^[0-9]+$`)},
	)
}

func TestSchemeParse(t *testing.T) {
	scheme := newTestScheme()

	key, err := scheme.Parse("us-123")

	require.NoError(t, err)

	storefront, err := scheme.Part(key, "storefront")
	require.NoError(t, err)
	require.Equal(t, "us", storefront)

	id, err := scheme.Part(key, "id")
	require.NoError(t, err)
	require.Equal(t, "123", id)

	missing, err := scheme.Part(key, "missing")
	require.Equal(t, "", missing)
	require.ErrorContains(t, err, `missing key part: "missing"`)
}

func TestSchemeParseRejectsInvalidPartCount(t *testing.T) {
	scheme := newTestScheme()

	_, err := scheme.Parse("us-123-extra")

	require.ErrorContains(t, err, "invalid key")
}

func TestSchemeParseRejectsInvalidPartValue(t *testing.T) {
	scheme := newTestScheme()

	_, err := scheme.Parse("us-abc")

	require.ErrorContains(t, err, "invalid key part")
}

func TestSchemeNewKey(t *testing.T) {
	scheme := newTestScheme()

	key, err := scheme.NewKey(map[string]string{
		"storefront": "us",
		"id":         "123",
	})

	require.NoError(t, err)

	storefront, err := scheme.Part(key, "storefront")
	require.NoError(t, err)
	require.Equal(t, "us", storefront)

	id, err := scheme.Part(key, "id")
	require.NoError(t, err)
	require.Equal(t, "123", id)
}

func TestSchemeNewKeyRejectsUnknownPart(t *testing.T) {
	scheme := newTestScheme()

	_, err := scheme.NewKey(map[string]string{
		"storefront": "us",
		"id":         "123",
		"extra":      "value",
	})

	require.ErrorContains(t, err, "unknown key part")
}

func TestSchemeNewKeyCopiesValues(t *testing.T) {
	scheme := newTestScheme()
	values := map[string]string{
		"storefront": "us",
		"id":         "123",
	}

	key, err := scheme.NewKey(values)
	require.NoError(t, err)
	values["id"] = "456"

	dump, err := scheme.Dump(key)

	require.NoError(t, err)
	require.Equal(t, "us-123", dump)
}

func TestSchemeDump(t *testing.T) {
	scheme := newTestScheme()
	key, err := scheme.NewKey(map[string]string{
		"storefront": "us",
		"id":         "123",
	})
	require.NoError(t, err)

	dump, err := scheme.Dump(key)

	require.NoError(t, err)
	require.Equal(t, "us-123", dump)
}

func TestSchemeDumpRejectsMissingPart(t *testing.T) {
	scheme := newTestScheme()

	_, err := scheme.NewKey(map[string]string{
		"storefront": "us",
	})

	require.ErrorContains(t, err, "missing key part")
}

func TestSchemeDumpRejectsKeyFromDifferentScheme(t *testing.T) {
	scheme := newTestScheme()
	otherScheme := NewScheme(
		":",
		Part{Name: "artist", Re: regexp.MustCompile(`^[a-z]+$`)},
		Part{Name: "slug", Re: regexp.MustCompile(`^[a-z]+$`)},
	)
	key, err := otherScheme.NewKey(map[string]string{
		"artist": "autechre",
		"slug":   "amber",
	})
	require.NoError(t, err)

	_, err = scheme.Dump(key)

	require.ErrorContains(t, err, "key belongs to different scheme")
}

func TestSchemeDumpRejectsZeroValueKey(t *testing.T) {
	scheme := newTestScheme()
	var key Key

	_, err := scheme.Dump(key)

	require.ErrorContains(t, err, "missing key part")
}
