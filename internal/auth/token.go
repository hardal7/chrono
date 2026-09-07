package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/hardal7/chrono/internal/util/config"
)

const tokenBytes = 32

func generateToken() (string, error) {
	b := make([]byte, tokenBytes)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(token string) string {
	secret := []byte(config.App.HashSecret)

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(token))

	return hex.EncodeToString(mac.Sum(nil))
}
