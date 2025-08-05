package processData

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
)
import "github.com/evgenyshipko/golang-metrics-collector/internal/agent/encrypt"

type EncryptBodyProcessor struct {
	CryptoPublicKeyPath string
}

func (p *EncryptBodyProcessor) Process(data []byte, headers map[string]string) ([]byte, map[string]string, error) {

	if p.CryptoPublicKeyPath == "" {
		return data, headers, nil
	}

	publicKey, err := encrypt.LoadPublicKey(p.CryptoPublicKeyPath)
	if err != nil {
		logger.Instance.Warnf("Failed to load public key: %v", err)
		return nil, nil, err
	}

	encryptedMsg, err := encrypt.EncryptWithPublicKey(data, publicKey)
	if err != nil {
		logger.Instance.Warnf("Failed to encrypt message: %v", err)
		return nil, nil, err
	}

	headers["X-Encrypted"] = "true"

	return encryptedMsg, headers, nil
}
