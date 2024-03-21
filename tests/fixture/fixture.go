package fixture

import (
	"embed"
	"log"
)

//go:embed */**
var embeddedFiles embed.FS

func Read(path string) []byte {
	file, err := embeddedFiles.ReadFile(path)
	if err != nil {
		log.Fatal("failed to read fixture", err)
	}
	return file
}
