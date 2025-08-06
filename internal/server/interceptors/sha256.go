package interceptors

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	sha256utils "github.com/evgenyshipko/golang-metrics-collector/internal/common/utils"
	setup2 "github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func Sha256ServerInterceptor(cfg setup2.ServerStartupValues) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		token := ""
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := md.Get(consts.HashSha256Header)
			if len(values) > 0 {
				// ключ содержит слайс строк, получаем первую строку
				token = values[0]
			}
		}

		if token != "" && cfg.HashKey != "" {
			var rawData []byte
			var err error
			if msg, ok := req.(proto.Message); ok {
				rawData, err = proto.Marshal(msg)
				if err != nil {
					logger.Instance.Warnf("error in SHA256 interceptor: %v", err)
					return nil, status.Error(codes.Internal, "error in SHA256 interceptor")
				}
			}

			tokenFromRequestData := sha256utils.GetHashedString(cfg.HashKey, rawData)

			if tokenFromRequestData != token {
				return nil, status.Error(codes.PermissionDenied, "invalid data")
			}

		}

		return handler(ctx, req)
	}
}
