package cpuScout

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/scouts"
)

// Storage хранилище.
type Storage struct {
	mu   sync.RWMutex
	rows map[int64]scouts.MertricRow
}

// NewCPUStorage конструктор хранилища.
func NewCPUStorage() *Storage {
	return &Storage{rows: make(map[int64]scouts.MertricRow, 0)}
}

// Save сохранить в хранилище.
func (s *Storage) Save(r scouts.MertricRow) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rows[r.Date.Unix()] = r
	return true
}

// GetByRange получить значения за период.
func (s *Storage) GetByRange(r time.Duration) []scouts.MertricRow {
	result := make([]scouts.MertricRow, 0)
	tsmpNow := time.Now().Unix()
	tsmpLowBorder := time.Now().Unix() - int64(r/time.Second)

	for timestamp, row := range s.rows {
		if timestamp <= tsmpNow && timestamp >= tsmpLowBorder {
			result = append(result, row)
		}
	}

	return result
}

// GetAvgByRange получить среднее значение за период.
func (s *Storage) GetAvgByRange(r time.Duration) (scouts.MertricRow, error) {
	metrics := s.GetByRange(r)
	if len(metrics) == 0 {
		return scouts.MertricRow{}, nil
	}

	if len(metrics) == 1 {
		return metrics[0], nil
	}

	user := make([]float32, 0, len(metrics))
	system := make([]float32, 0, len(metrics))
	idle := make([]float32, 0, len(metrics))
	var lastMetricDate time.Time
	var metricName = metrics[0].Name

	for _, metric := range metrics {
		str := strings.ReplaceAll(metric.Body, ",", ".")
		values := strings.Split(str, " ")
		usr, err := strconv.ParseFloat(values[0], 32)
		if err != nil {
			return scouts.MertricRow{}, fmt.Errorf("parse float: %w", err)
		}
		user = append(user, float32(usr))

		sys, err := strconv.ParseFloat(values[1], 32)
		if err != nil {
			return scouts.MertricRow{}, fmt.Errorf("parse float: %w", err)
		}
		system = append(system, float32(sys))

		idl, err := strconv.ParseFloat(values[2], 32)
		if err != nil {
			return scouts.MertricRow{}, fmt.Errorf("parse float: %w", err)
		}
		idle = append(idle, float32(idl))
		lastMetricDate = metric.Date
	}

	avgValues := calcAvg(user, system, idle)

	data := fmt.Sprintf("%.2f, %.2f, %.2f", avgValues[0], avgValues[1], avgValues[2])
	return scouts.MertricRow{Date: lastMetricDate, Body: data, Name: metricName}, nil
}

// calcAvg подсчет средних по 3 слайсам значений.
func calcAvg(usr, sys, idl []float32) []float32 {
	var usrSum float32
	var sysSum float32
	var idlSum float32

	for _, v := range usr {
		usrSum += v
	}

	for _, v := range sys {
		sysSum += v
	}

	for _, v := range idl {
		idlSum += v
	}

	return []float32{usrSum / float32(len(usr)), sysSum / float32(len(usr)), idlSum / float32(len(usr))}
}
