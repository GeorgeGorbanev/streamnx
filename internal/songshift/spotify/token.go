package spotify

import "fmt"

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (t *Token) authHeader() string {
	return fmt.Sprintf("Bearer %s", t.AccessToken)
}
