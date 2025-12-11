package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Sign(appID uint64, appSecret string) string {
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write([]byte(Uint64ToString(appID)))

	signatureBytes := h.Sum(nil)

	return hex.EncodeToString(signatureBytes)
}

func VerifySignature(appID uint64, appSecret, sign string) bool {
	expectedSignature := Sign(appID, appSecret)
	return expectedSignature == sign
}
