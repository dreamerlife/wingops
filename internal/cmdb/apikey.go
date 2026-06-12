package cmdb

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type APIKey struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	KeyID  string `json:"key_id"`
	Secret string `json:"secret,omitempty"`
	Status string `json:"status"`
}

func NewDevelopmentAPIKey() APIKey {
	return APIKey{
		KeyID:  "dev-sync-key",
		Secret: "dev-sync-secret",
		Status: "active",
	}
}

func (k APIKey) Active() bool {
	return k.Status == "" || k.Status == "active"
}

func (k APIKey) VerifySignature(body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(k.Secret))
	_, _ = mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
