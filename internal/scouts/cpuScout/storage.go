package cpuScout

import (
	"errors"
	"sync"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/scouts"
)

// Storage хранилище.
type Storage struct {
	mu   sync.RWMutex
	rows map[int64]scouts.MertricRow
}

// CpuData модель данных CPU
type CpuData struct {
	Usr  float32
	Sys  float32
	Idle float32
	Date time.Time
	Name string
}

var ErrFewAvgValues = errors.New("few avg values")

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
		return scouts.MertricRow{Date: metrics[0].Date, Body: metrics[0], Name: metrics[0].Name}, nil
	}

	user := make([]float32, 0, len(metrics))
	system := make([]float32, 0, len(metrics))
	idle := make([]float32, 0, len(metrics))

	var lastMetricDate time.Time
	var metricName = metrics[0].Name

	for _, metric := range metrics {
		data := metric.Body.(CpuData)
		user = append(user, data.Usr)
		system = append(system, data.Sys)
		idle = append(idle, data.Idle)
		lastMetricDate = metric.Date
	}

	avgValues := calcAvg(user, system, idle)

	if len(avgValues) != 3 {
		return scouts.MertricRow{}, ErrFewAvgValues
	}

	data := CpuData{
		Usr:  avgValues[0],
		Sys:  avgValues[1],
		Idle: avgValues[2],
		Name: metricName,
		Date: lastMetricDate,
	}

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
