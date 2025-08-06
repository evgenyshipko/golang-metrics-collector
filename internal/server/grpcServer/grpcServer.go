package grpcServer

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/services"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/evgenyshipko/golang-metrics-collector/internal/common/grpc"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServiceServer
	service services.Service
}

func (s *MetricsServer) UpdateMetric(ctx context.Context, in *pb.Metric) (*pb.MetricResponse, error) {

	var metricData *consts.MetricData
	if in.Type == pb.MetricType_GAUGE {
		value := in.GetVal()
		metricData = &consts.MetricData{
			ID:    in.Id,
			MType: consts.GAUGE,
			Value: &value,
		}
	} else if in.Type == pb.MetricType_COUNTER {
		delta := in.GetDelta()
		metricData = &consts.MetricData{
			ID:    in.Id,
			MType: consts.COUNTER,
			Delta: &delta,
		}
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

func StartGrpcServer(metricService services.Service) {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterMetricsServiceServer(s, &MetricsServer{service: metricService})

	logger.Instance.Info("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
