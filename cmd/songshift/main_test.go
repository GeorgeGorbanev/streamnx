package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMainOutput(t *testing.T) {
	output := captureOutput(main)
	require.Equal(t, "Hello, world!\n", output)
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	if err := w.Close(); err != nil {
		log.Fatal(err)
	}
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
