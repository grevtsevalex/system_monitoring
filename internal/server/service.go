package server

import (
	"log"
	"time"

	serverpb "github.com/grevtsevalex/system_monitoring/internal/server/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// defaultInterval интервал отправки данных клиенту.
const defaultInterval = 5 * time.Second

// Service модель сервиса.
type Service struct {
	serverpb.UnimplementedMonitoringServiceServer
	interval time.Duration
}

// NewService конструктор сервиса.
func NewService() *Service {
	return &Service{
		interval: defaultInterval,
	}
}

// GetMetrics получение метрик.
func (s *Service) GetMetrics(req *serverpb.GetMetricsRequest, srv serverpb.MonitoringService_GetMetricsServer) error {
	log.Printf("new listener")

L:
	for {
		select {
		case <-srv.Context().Done():
			log.Printf("listener disconnected")
			break L

		case <-time.After(s.interval):
			for i := range 10 {
				msg := &serverpb.Snapshot{
					Number: uint32(i),
					Time:   timestamppb.Now(),
				}

				if err := srv.Send(msg); err != nil {
					log.Printf("unable to send message to stats listener: %v", err)
					break L
				}
			}
		}
	}

	return nil
}
