package processData

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Sha256Processor struct {
	HashKey string
}

func (p *Sha256Processor) Process(data []byte, headers map[string]string) ([]byte, map[string]string, error) {

	h := hmac.New(sha256.New, []byte(p.HashKey))
	h.Write(data)
	headers["HashSHA256"] = hex.EncodeToString(h.Sum(nil))

	return data, headers, nil
}
