package discscout

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

// DiscData модель данных загрузки дисков.
type DiscData struct {
	Devices []Device
	Date    time.Time
	Name    string
}

// Device модель устройства.
type Device struct {
	Name string
	Tps  float32
	Rps  float32
	Wps  float32
}

var ErrFewAvgValues = errors.New("few avg values")

// NewDiscStorage конструктор хранилища.
func NewDiscStorage() *Storage {
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

	var lastMetricDate time.Time
	var metricName = metrics[0].Name

	discsTps := make(map[string][]float32)
	discsRps := make(map[string][]float32)
	discsWps := make(map[string][]float32)

	for _, metric := range metrics {
		data := metric.Body.(DiscData)
		for _, deviceInfo := range data.Devices {
			discsTps[deviceInfo.Name] = append(discsTps[deviceInfo.Name], deviceInfo.Tps)
			discsRps[deviceInfo.Name] = append(discsRps[deviceInfo.Name], deviceInfo.Rps)
			discsWps[deviceInfo.Name] = append(discsWps[deviceInfo.Name], deviceInfo.Wps)
			lastMetricDate = metric.Date
		}
	}

	discsTpsAvg, discsRpsAvg, discsWpsAvg := calcAvg(discsTps, discsRps, discsWps)
	avgDiscData := DiscData{
		Name: metricName,
		Date: lastMetricDate,
	}

	for deviceName, avgTps := range discsTpsAvg {
		device := Device{
			Name: deviceName,
			Tps:  avgTps,
			Rps:  discsRpsAvg[deviceName],
			Wps:  discsWpsAvg[deviceName],
		}
		avgDiscData.Devices = append(avgDiscData.Devices, device)
	}

	return scouts.MertricRow{Date: lastMetricDate, Body: avgDiscData, Name: metricName}, nil
}

// calcAvg подсчет средних по 3 слайсам значений.
func calcAvg(discsTps, discsdRps, discsWps map[string][]float32) (map[string]float32, map[string]float32, map[string]float32) {
	tpsResult := make(map[string]float32)
	rpsResult := make(map[string]float32)
	wpsResult := make(map[string]float32)

	for deviceName, values := range discsTps {
		var sumValues float32
		for _, val := range values {
			sumValues += val
		}
		tpsResult[deviceName] = sumValues / float32(len(values))
	}

	for deviceName, values := range discsdRps {
		var sumValues float32
		for _, val := range values {
			sumValues += val
		}
		rpsResult[deviceName] = sumValues / float32(len(values))
	}

	for deviceName, values := range discsWps {
		var sumValues float32
		for _, val := range values {
			sumValues += val
		}
		wpsResult[deviceName] = sumValues / float32(len(values))
	}

	return tpsResult, rpsResult, wpsResult
}
