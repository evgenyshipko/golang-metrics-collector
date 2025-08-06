package requests

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	pb "github.com/evgenyshipko/golang-metrics-collector/internal/common/grpc"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type GrpcRequester struct {
	connection *grpc.ClientConn
}

func (r *GrpcRequester) SendMetric(metric consts.MetricData) error {

	client := pb.NewMetricsServiceClient(r.connection)

	ctx := context.Background()

	var metricDataPb *pb.Metric

	logger.Instance.Infow("Sending metric", "metric", metric)

	if metric.MType == consts.GAUGE {
		metricDataPb = &pb.Metric{
			Id:    metric.ID,
			Type:  pb.MetricType_GAUGE,
			Value: &pb.Metric_Val{Val: *metric.Value},
		}
	} else if metric.MType == consts.COUNTER {
		metricDataPb = &pb.Metric{
			Id:    metric.ID,
			Type:  pb.MetricType_COUNTER,
			Value: &pb.Metric_Delta{Delta: *metric.Delta},
		}
	}

	resp, err := client.UpdateMetric(ctx, metricDataPb)

	if err != nil {
		logger.Instance.Warnf("Ошибка grpc-метода: %s", err.Error())
		return err
	}

	logger.Instance.Infof("Ответ grpc-метода: %s", resp.String())
	return nil
}

func NewGrpcRequester(_ setup.AgentStartupValues) *GrpcRequester {
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
