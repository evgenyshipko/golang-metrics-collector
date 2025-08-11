package processData

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"net"
)

type XRealIpProcessor struct {
	OutboundIP net.IP
}

func (p *XRealIpProcessor) Process(data []byte, headers map[string]string) ([]byte, map[string]string, error) {

	if p.OutboundIP != nil {
		headers[consts.XRealIpHeader] = p.OutboundIP.String()
	}

	return data, headers, nil
}
