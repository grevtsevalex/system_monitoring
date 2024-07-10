package collector

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/scouts"
)

// Collector модель сборщика метрик.
type Collector struct {
	storages []scouts.ScoutStorage
}

// NewCollector конструктор коллектора.
func NewCollector(storages []scouts.ScoutStorage) Collector {
	return Collector{storages: storages}
}

// GetSnapshot получение снэпшота системы за период.
func (c *Collector) GetSnapshot(r time.Duration) string {
	var sb strings.Builder
	var wg sync.WaitGroup
	for _, st := range c.storages {
		wg.Add(1)
		go func(st scouts.ScoutStorage) {
			defer wg.Done()
			metric, err := st.GetAvgByRange(r)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				sb.WriteString(metric.String())
			}
		}(st)
	}

	wg.Wait()

	return sb.String()
}
