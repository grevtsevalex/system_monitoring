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
				s.writeCpu(msg, sn.Cpu)
			}

			if sn.LA.Filled {
				s.writeLa(msg, sn.LA)
			}

			if sn.Disc.Filled {
				s.writeDisc(msg, sn.Disc)
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

// writeLa записать в сообщение данные о LoadAverages.
func (s *Service) writeLa(msg *serverpb.Snapshot, data collector.LA) {
	msg.La = &serverpb.LAMessage{
		PerMinute:    data.PerMinute,
		Per5Minutes:  data.Per5Minute,
		Per15Minutes: data.Per15Minute,
	}
}

// writeCpu записать в сообщение данные о CPU.
func (s *Service) writeCpu(msg *serverpb.Snapshot, data collector.Cpu) {
	msg.Cpu = &serverpb.CpuMessage{
		Usr:  data.Usr,
		Sys:  data.Sys,
		Idle: data.Idle,
	}
}

// writeDisc записать в сообщение данные о нагрузке дисков.
func (s *Service) writeDisc(msg *serverpb.Snapshot, data collector.DiscData) {
	msg.Disc = &serverpb.DiscMessage{}
	for _, deviceInfo := range data.Devices {
		device := serverpb.Device{
			Name: deviceInfo.Name,
			Tps:  deviceInfo.Tps,
			Rps:  deviceInfo.Rps,
			Wps:  deviceInfo.Wps,
		}
		msg.Disc.Devices = append(msg.Disc.Devices, &device)
	}
}
