package commonUtils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GetHashedString(hashKey string, data []byte) string {
	h := hmac.New(sha256.New, []byte(hashKey))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
