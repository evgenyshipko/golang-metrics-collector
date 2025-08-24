package grpcServer

import (
	"context"
	"errors"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	pb "github.com/evgenyshipko/golang-metrics-collector/internal/common/grpc"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/interceptors"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/services"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"net"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServiceServer
	service services.Service
}

type GrpcServer struct {
	Address string
	Server  *grpc.Server
}

func (s *MetricsServer) UpdateMetric(ctx context.Context, in *pb.Metric) (*pb.MetricResponse, error) {

	var metricData *consts.MetricData
	switch in.Type {
	case pb.MetricType_GAUGE:
		value := in.GetVal()
		metricData = &consts.MetricData{
			ID:    in.Id,
			MType: consts.GAUGE,
			Value: &value,
		}
	case pb.MetricType_COUNTER:
		delta := in.GetDelta()
		metricData = &consts.MetricData{
			ID:    in.Id,
			MType: consts.COUNTER,
			Delta: &delta,
		}
	default:
		return nil, errors.New("metric type not should be gauge or counter")
	}

	_, err := s.service.ProcessMetric(ctx, *metricData)
	if err != nil {
		logger.Instance.Warnf("Ошибка записи метрик через grpc-метод: %v", err)
		return &pb.MetricResponse{}, err
	}

	logger.Instance.Infow("Метрики записаны успешно", "metricData", *metricData)

	return &pb.MetricResponse{
		Success: true,
		Message: "Metric processed",
	}, nil
}

func CreateGrpcServer(metricService services.Service, cfg setup.ServerStartupValues) *GrpcServer {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.Sha256(cfg), interceptors.TrustedIp(cfg)))
	// регистрируем сервис
	pb.RegisterMetricsServiceServer(server, &MetricsServer{service: metricService})
	return &GrpcServer{Address: cfg.Host, Server: server}
}

func (s *GrpcServer) Start() error {
	listen, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	return s.Server.Serve(listen)
}

func (s *GrpcServer) Stop() error {
	s.Server.GracefulStop()
	return nil
}
