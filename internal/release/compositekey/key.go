package compositekey

import (
	"fmt"
	"regexp"
	"strings"
)

type Key struct {
	scheme *Scheme
	values map[string]string
}

type Scheme struct {
	delimiter string
	parts     []Part
}

type Part struct {
	Name string
	Re   *regexp.Regexp
}

func NewScheme(delimiter string, parts ...Part) *Scheme {
	return &Scheme{
		delimiter: delimiter,
		parts:     append([]Part(nil), parts...),
	}
}

func (s *Scheme) NewKey(parts map[string]string) (Key, error) {
	for name := range parts {
		if !s.hasPart(name) {
			return Key{}, fmt.Errorf("unknown key part %q", name)
		}
	}

	partsCopy := make(map[string]string, len(s.parts))
	for _, part := range s.parts {
		value, ok := parts[part.Name]
		switch {
		case !ok:
			return Key{}, fmt.Errorf("missing key part %q", part.Name)
		case part.Re == nil:
			return Key{}, fmt.Errorf("invalid key part %q: missing regexp", part.Name)
		case !part.Re.MatchString(value):
			return Key{}, fmt.Errorf("invalid key part %q: value %q does not match %s",
				part.Name, value, part.Re.String())
		default:
			partsCopy[part.Name] = value
		}
	}

	return Key{
		scheme: s,
		values: partsCopy,
	}, nil
}

func (s *Scheme) Parse(dump string) (Key, error) {
	parts := strings.Split(dump, s.delimiter)
	if len(parts) != len(s.parts) {
		return Key{}, fmt.Errorf("invalid key: expected %d parts, got %d", len(s.parts), len(parts))
	}

	partsMap := make(map[string]string, len(s.parts))
	for i, part := range s.parts {
		partsMap[part.Name] = parts[i]
	}

	return s.NewKey(partsMap)
}

func (s *Scheme) Dump(key Key) (string, error) {
	if key.scheme != nil && key.scheme != s {
		return "", fmt.Errorf("key belongs to different scheme")
	}

	parts := make([]string, 0, len(s.parts))
	for _, schemePart := range s.parts {
		part, err := s.Part(key, schemePart.Name)
		if err != nil {
			return "", err
		}
		parts = append(parts, part)
	}

	return strings.Join(parts, s.delimiter), nil
}

func (s *Scheme) Part(key Key, name string) (string, error) {
	value, ok := key.values[name]
	if !ok {
		return "", fmt.Errorf("missing key part: %q", name)
	}
	return value, nil
}

func (s *Scheme) hasPart(name string) bool {
	for _, part := range s.parts {
		if part.Name == name {
			return true
		}
	}
	return false
}
