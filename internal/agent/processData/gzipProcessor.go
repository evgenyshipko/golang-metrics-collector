package processData

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/gzip"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
)

type GZipProcessor struct{}

func (p *GZipProcessor) Process(data []byte, headers map[string]string) ([]byte, map[string]string, error) {

	headers["Content-Encoding"] = "gzip"

	compressedBody, err := gzip.Compress(data)
	if err != nil {
		logger.Instance.Warnw("GZipProcessor", "gzip.Compress err", err)
		return []byte{}, headers, err
	}

	return compressedBody, headers, nil
}
