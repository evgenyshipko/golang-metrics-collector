package requests

import (
	"context"
	"errors"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/interceptors"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	pb "github.com/evgenyshipko/golang-metrics-collector/internal/common/grpc"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"log"
	"time"
)

type GrpcRequester struct {
	connection *grpc.ClientConn
	timeout    time.Duration
}

func (r *GrpcRequester) SendMetric(metric consts.MetricData) error {

	client := pb.NewMetricsServiceClient(r.connection)

	var metricDataPb *pb.Metric

	logger.Instance.Infow("Sending metric", "metric", metric)

	switch metric.MType {
	case consts.GAUGE:
		metricDataPb = &pb.Metric{
			Id:    metric.ID,
			Type:  pb.MetricType_GAUGE,
			Value: &pb.Metric_Val{Val: *metric.Value},
		}
	case consts.COUNTER:
		metricDataPb = &pb.Metric{
			Id:    metric.ID,
			Type:  pb.MetricType_COUNTER,
			Value: &pb.Metric_Delta{Delta: *metric.Delta},
		}
	default:
		return errors.New("invalid metric type, should be gauge or counter")
	}

	resp, err := client.UpdateMetric(context.Background(), metricDataPb)

	if err != nil {
		logger.Instance.Warnf("Ошибка grpc-метода: %s", err.Error())
		return err
	}

	logger.Instance.Infof("Ответ grpc-метода: %s", resp.String())
	return nil
}

func NewGrpcRequester(cfg setup.AgentStartupValues) *GrpcRequester {
	conn, err := grpc.NewClient(cfg.Host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		grpc.WithChainUnaryInterceptor(interceptors.RetryInterceptor(cfg), interceptors.Sha256Interceptor(cfg), interceptors.XRealIpInterceptor(cfg)),
	)
	if err != nil {
		log.Fatal(err)
	}

	return &GrpcRequester{
		connection: conn,
	}
}

func (r *GrpcRequester) Close() error {
	return r.connection.Close()
}
