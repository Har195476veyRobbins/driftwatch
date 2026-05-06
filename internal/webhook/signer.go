package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Signer signs a raw byte payload and returns a hex-encoded signature.
type Signer interface {
	Sign(payload []byte) string
}

// HMACSigner signs payloads using HMAC-SHA256.
type HMACSigner struct {
	secret []byte
}

// NewHMACSigner returns an HMACSigner using the provided secret.
func NewHMACSigner(secret string) *HMACSigner {
	return &HMACSigner{secret: []byte(secret)}
}

// Sign returns the HMAC-SHA256 hex digest of payload.
func (h *HMACSigner) Sign(payload []byte) string {
	mac := hmac.New(sha256.New, h.secret)
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}
