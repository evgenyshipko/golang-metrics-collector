package processData

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	sha256utils "github.com/evgenyshipko/golang-metrics-collector/internal/common/utils"
)

type Sha256Processor struct {
	HashKey string
}

func (p *Sha256Processor) Process(data []byte, headers map[string]string) ([]byte, map[string]string, error) {

	if p.HashKey == "" {
		return data, headers, nil
	}

	headers[consts.HashSha256Header] = sha256utils.GetHashedString(p.HashKey, data)

	return data, headers, nil
}
