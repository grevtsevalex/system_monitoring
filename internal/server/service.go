package server

import (
	"log"
	"math"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/collector"
	serverpb "github.com/grevtsevalex/system_monitoring/internal/server/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Service модель сервиса.
type Service struct {
	serverpb.UnimplementedMonitoringServiceServer
	clc collector.Collector
}

const (
	minSecondsForPeriod = 2
	maxSecondsForRange  = 60
)

// NewService конструктор сервиса.
func NewService(clc collector.Collector) *Service {
	return &Service{clc: clc}
}

// GetMetrics получение метрик.
func (s *Service) GetMetrics(req *serverpb.GetMetricsRequest, srv serverpb.MonitoringService_GetMetricsServer) error {
	log.Printf("new listener")

	periodSeconds := math.Max(float64(req.PeriodSec), minSecondsForPeriod)
	rangeSeconds := math.Min(float64(req.RangeSec), maxSecondsForRange)
	period := time.Duration(periodSeconds) * time.Second
	rang := time.Duration(rangeSeconds) * time.Second

L:
	for {
		select {
		case <-srv.Context().Done():
			log.Printf("listener disconnected")
			break L

		case <-time.After(period):
			msg := &serverpb.Snapshot{
				Msg:  s.clc.GetSnapshot(rang),
				Time: timestamppb.Now(),
			}

			if err := srv.Send(msg); err != nil {
				log.Printf("unable to send message to stats listener: %v", err)
				break L
			}
		}
	}

	return nil
}
