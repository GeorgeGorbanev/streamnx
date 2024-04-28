package spotify

import (
	"fmt"
	"time"
)

type token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	fetchedAt   time.Time
}

func (t *token) authHeader() string {
	return fmt.Sprintf("Bearer %s", t.AccessToken)
}

func (t *token) isExpired() bool {
	return time.Since(t.fetchedAt) > time.Duration(t.ExpiresIn)*time.Second
}
