package collector

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/scouts"
	"github.com/grevtsevalex/system_monitoring/internal/scouts/cpuScout"
	lascout "github.com/grevtsevalex/system_monitoring/internal/scouts/laScout"
)

// Logger тип логгера.
type Logger interface {
	Log(msg string)
	Error(msg string)
}

// Collector модель сборщика метрик.
type Collector struct {
	storages []scouts.ScoutStorage
	logger   Logger
}

// NewCollector конструктор коллектора.
func NewCollector(storages []scouts.ScoutStorage, logger Logger) Collector {
	return Collector{storages: storages, logger: logger}
}

var (
	ErrTypeAssertionLA  = errors.New("failed type assertion LA")
	ErrTypeAssertionCPU = errors.New("failed type assertion CPU")
)

// Snapshot модель снэпшота.
type Snapshot struct {
	Cpu Cpu
	LA  LA
}

// LA модель данных LoadAverage.
type LA struct {
	PerMinute   float32
	Per5Minute  float32
	Per15Minute float32
	Filled      bool
}

// Cpu модель данных Cpu.
type Cpu struct {
	Usr    float32
	Sys    float32
	Idle   float32
	Filled bool
}

// GetSnapshot получение снэпшота системы за период.
func (c *Collector) GetSnapshot(r time.Duration) *Snapshot {
	sn := &Snapshot{}
	var wg sync.WaitGroup
	for _, st := range c.storages {
		wg.Add(1)
		go func(st scouts.ScoutStorage) {
			defer wg.Done()
			metric, err := st.GetAvgByRange(r)
			if err != nil {
				c.logger.Error(fmt.Sprintf("получение усредненной метрики за период: %s", err.Error()))
			} else {
				err := sn.fill(metric)
				if err != nil {
					c.logger.Error(fmt.Sprintf("наполнение модели снэпшота: %s", err.Error()))
				}
			}
		}(st)
	}

	wg.Wait()

	return sn
}

// fill наполнить модель данными.
func (s *Snapshot) fill(metric scouts.MertricRow) error {
	if metric.Name == cpuScout.MetricName {
		data, ok := metric.Body.(cpuScout.CpuData)
		if !ok {
			return ErrTypeAssertionCPU
		}
		s.Cpu = Cpu{
			Usr:    data.Usr,
			Sys:    data.Sys,
			Idle:   data.Idle,
			Filled: true,
		}
	}

	if metric.Name == lascout.MetricName {
		data, ok := metric.Body.(lascout.LAData)
		if !ok {
			return ErrTypeAssertionLA
		}
		s.LA = LA{
			PerMinute:   data.PerMinute,
			Per5Minute:  data.Per5Minute,
			Per15Minute: data.Perf15Minute,
			Filled:      true,
		}
	}

	return nil
}
