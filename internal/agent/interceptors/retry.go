package interceptors

import (
	"context"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func RetryInterceptor(cfg setup.AgentStartupValues) func(ctx context.Context,
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
		maxRetries := len(cfg.RetryIntervals)
		var err error

		ctx, cancel := context.WithTimeout(ctx, cfg.RequestWaitTimeout)

		defer cancel()

		for i := 0; i < maxRetries; i++ {
			err = invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}
			logger.Instance.Info(fmt.Sprintf("Попытка %d, ждем %v перед следующим запросом...\n", i, cfg.RetryIntervals[i]))

			// Проверяем, стоит ли повторять
			if status.Code(err) != codes.Unavailable {
				return err
			}

			time.Sleep(cfg.RetryIntervals[i])
		}
		return err
	}
}
