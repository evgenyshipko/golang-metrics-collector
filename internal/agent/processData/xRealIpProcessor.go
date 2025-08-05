package processData

import (
	"net"
)

type XRealIpProcessor struct {
	OutboundIP net.IP
}

func (p *XRealIpProcessor) Process(data []byte, headers map[string]string) ([]byte, map[string]string, error) {

	if p.OutboundIP != nil {
		headers["X-real-ip"] = p.OutboundIP.String()
	}

	return data, headers, nil
}
