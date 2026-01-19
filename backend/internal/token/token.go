package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type HashToken struct {
	SecretKey []byte
}

type Tokens struct {
	access_token  string
	refresh_token string
}

func NewHashToken(secret_key string) *HashToken {
	return &HashToken{
		SecretKey: []byte(secret_key),
	}
}

func (t *HashToken) HashRefreshToken(token string) string {
	h := hmac.New(sha256.New, t.SecretKey)
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func (t *HashToken) CompareToken(storedHash, token string) bool {
	tokenHash := t.HashRefreshToken(token)
	return hmac.Equal([]byte(storedHash), []byte(tokenHash))
}
