package interceptors

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	setup2 "github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

func TrustedIp(cfg setup2.ServerStartupValues) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		ipString := ""
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := md.Get(consts.XRealIpHeader)
			if len(values) > 0 {
				// ключ содержит слайс строк, получаем первую строку
				ipString = values[0]
			}
		}

		if ipString != "" && cfg.TrustedSubnet != nil {

			ip := net.ParseIP(ipString)
			if ip == nil {
				return nil, status.Error(codes.PermissionDenied, "Invalid IP address in X-Real-IP header")
			}

			if !cfg.TrustedSubnet.Contains(ip) {
				return nil, status.Error(codes.PermissionDenied, "Access denied: IP not in trusted subnet")
			}
		}

		return handler(ctx, req)
	}
}
