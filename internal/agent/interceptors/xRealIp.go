package interceptors

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func XRealIpInterceptor(cfg setup.AgentStartupValues) func(ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		if cfg.OutboundIp != nil {
			ip := *cfg.OutboundIp
			md := metadata.New(map[string]string{consts.XRealIpHeader: ip.String()})
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		err := invoker(ctx, method, req, reply, cc, opts...)

		return err
	}
}
