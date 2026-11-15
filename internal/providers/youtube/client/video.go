package client

import "errors"

type Video struct {
	Metadata
	Duration int
}

type Metadata struct {
	ID          string
	Title       string
	Creator     string
	Description string
	CoverURL    string
}

const SearchLimit = 10

var ErrNotFound = errors.New("not found")
