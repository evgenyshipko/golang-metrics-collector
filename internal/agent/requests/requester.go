package requests

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
)

type Protocol string

const (
	HTTP = "http"
	GRPC = "grpc"
)

type Requester interface {
	SendMetric(metric consts.MetricData) error
	Close() error
}

func NewRequester(cfg setup.AgentStartupValues) Requester {
	if cfg.Protocol == GRPC {
		return NewGrpcRequester(cfg)
	} else {
		return NewHttpRequester(cfg)
	}
}
