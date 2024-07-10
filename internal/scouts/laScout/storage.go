package lascout

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

// NewLAStorage конструктор хранилища.
func NewLAStorage() *Storage {
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

	minute := make([]float32, 0, len(metrics))
	minuteX5 := make([]float32, 0, len(metrics))
	minuteX15 := make([]float32, 0, len(metrics))
	var lastMetricDate time.Time
	var metricName = metrics[0].Name

	for _, metric := range metrics {
		str := strings.ReplaceAll(metric.Body, ",", ".")
		values := strings.Split(str, ". ")
		byMinute, err := strconv.ParseFloat(values[0], 32)
		if err != nil {
			return scouts.MertricRow{}, fmt.Errorf("parse float: %w", err)
		}
		minute = append(minute, float32(byMinute))

		by5Minute, err := strconv.ParseFloat(values[1], 32)
		if err != nil {
			return scouts.MertricRow{}, fmt.Errorf("parse float: %w", err)
		}
		minuteX5 = append(minuteX5, float32(by5Minute))

		by15Minute, err := strconv.ParseFloat(values[2], 32)
		if err != nil {
			return scouts.MertricRow{}, fmt.Errorf("parse float: %w", err)
		}
		minuteX15 = append(minuteX15, float32(by15Minute))
		lastMetricDate = metric.Date
	}

	avgValues := calcAvg(minute, minuteX5, minuteX15)

	data := fmt.Sprintf("%.2f, %.2f, %.2f", avgValues[0], avgValues[1], avgValues[2])
	return scouts.MertricRow{Date: lastMetricDate, Body: data, Name: metricName}, nil
}

// calcAvg подсчет средних по 3 слайсам значений.
func calcAvg(min, min5, min15 []float32) []float32 {
	var minSum float32
	var min5Sum float32
	var min15Sum float32

	for _, v := range min {
		minSum += v
	}

	for _, v := range min5 {
		min5Sum += v
	}

	for _, v := range min15 {
		min15Sum += v
	}

	return []float32{minSum / float32(len(min)), min5Sum / float32(len(min)), min15Sum / float32(len(min))}
}
