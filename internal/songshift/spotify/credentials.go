package spotify

import (
	"encoding/base64"
	"fmt"
)

type Credentials struct {
	ClientID     string
	ClientSecret string
}

func (c *Credentials) authHeader() string {
	credentials := fmt.Sprintf("%s:%s", c.ClientID, c.ClientSecret)
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	return fmt.Sprintf("Basic %s", encodedCredentials)
}
