package interceptors

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	sha256utils "github.com/evgenyshipko/golang-metrics-collector/internal/common/commonUtils"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func Sha256Interceptor(cfg setup.AgentStartupValues) func(ctx context.Context,
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

		md := metadata.New(map[string]string{})

		if cfg.HashKey != "" {

			var rawData []byte
			var err error
			if msg, ok := req.(proto.Message); ok {
				rawData, err = proto.Marshal(msg)
				if err != nil {
					logger.Instance.Warnf("error in SHA256 interceptort: %v", err)
					return err
				}
			}

			hash := sha256utils.GetHashedString(cfg.HashKey, rawData)
			md.Set(consts.HashSha256Header, hash)
		}

		ctxWithMetadata := metadata.NewOutgoingContext(ctx, md)

		err := invoker(ctxWithMetadata, method, req, reply, cc, opts...)

		return err
	}
}
