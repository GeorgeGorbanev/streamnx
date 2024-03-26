package vibeshare

import (
	"fmt"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
)

type convertParams struct {
	ID     string
	Source music.Provider
	Target music.Provider
}

func (p *convertParams) marshal() []string {
	return []string{
		string(p.Source),
		p.ID,
		string(p.Target),
	}
}

func (p *convertParams) unmarshal(s []string) error {
	if len(s) != 3 {
		return fmt.Errorf("invalid convert params: %s", s)
	}
	sourceProvider := music.Provider(s[0])
	if !music.IsValidProvider(sourceProvider) {
		return fmt.Errorf("invalid source provider: %s", s[0])
	}
	targetProvider := music.Provider(s[2])
	if !music.IsValidProvider(targetProvider) {
		return fmt.Errorf("invalid target provider: %s", s[2])
	}

	p.Source = sourceProvider
	p.ID = s[1]
	p.Target = targetProvider

	return nil
}
