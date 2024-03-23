package fixture

import (
	"embed"
	"fmt"
)

//go:embed */**
var embeddedFiles embed.FS

func Read(path string) []byte {
	file, err := embeddedFiles.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("failed to read fixture: %v", err))
	}
	return file
}
