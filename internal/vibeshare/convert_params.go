package vibeshare

import (
	"fmt"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
)

type convertParams struct {
	ID     string
	Source *music.Provider
	Target *music.Provider
}

func (p *convertParams) marshal() []string {
	return []string{
		p.Source.Code,
		p.ID,
		p.Target.Code,
	}
}

func (p *convertParams) unmarshal(s []string) error {
	if len(s) != 3 {
		return fmt.Errorf("invalid convert params: %s", s)
	}
	sourceProvider := music.FindProviderByCode(s[0])
	if sourceProvider == nil {
		return fmt.Errorf("invalid source provider: %s", s[0])
	}
	targetProvider := music.FindProviderByCode(s[2])
	if targetProvider == nil {
		return fmt.Errorf("invalid target provider: %s", s[2])
	}

	p.Source = sourceProvider
	p.ID = s[1]
	p.Target = targetProvider

	return nil
}
