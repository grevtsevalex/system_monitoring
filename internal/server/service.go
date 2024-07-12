package server

import (
	"fmt"
	"math"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/collector"
	serverpb "github.com/grevtsevalex/system_monitoring/internal/server/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Service модель сервиса.
type Service struct {
	serverpb.UnimplementedMonitoringServiceServer
	clc    collector.Collector
	logger Logger
}

const (
	minSecondsForPeriod = 2
	maxSecondsForRange  = 60
)

// NewService конструктор сервиса.
func NewService(clc collector.Collector, logger Logger) *Service {
	return &Service{clc: clc, logger: logger}
}

// GetMetrics получение метрик.
func (s *Service) GetMetrics(req *serverpb.GetMetricsRequest, srv serverpb.MonitoringService_GetMetricsServer) error {
	s.logger.Log("new listener")

	periodSeconds := math.Max(float64(req.PeriodSec), minSecondsForPeriod)
	rangeSeconds := math.Min(float64(req.RangeSec), maxSecondsForRange)
	period := time.Duration(periodSeconds) * time.Second
	rang := time.Duration(rangeSeconds) * time.Second

L:
	for {
		select {
		case <-srv.Context().Done():
			s.logger.Log("listener disconnected")
			break L

		case <-time.After(period):
			sn := s.clc.GetSnapshot(rang)
			msg := &serverpb.Snapshot{}
			if sn.Cpu.Filled {
				msg.Cpu = &serverpb.CpuMessage{
					Usr:  sn.Cpu.Usr,
					Sys:  sn.Cpu.Sys,
					Idle: sn.Cpu.Idle,
				}
			}

			if sn.LA.Filled {
				msg.La = &serverpb.LAMessage{
					PerMinute:    sn.LA.PerMinute,
					Per5Minutes:  sn.LA.Per5Minute,
					Per15Minutes: sn.LA.Per15Minute,
				}
			}
			msg.Time = timestamppb.Now()

			if err := srv.Send(msg); err != nil {
				s.logger.Error(fmt.Sprintf("unable to send message to metrics listener: %v", err))
				break L
			}
		}
	}

	return nil
}
