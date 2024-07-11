package scouts

import (
	"time"
)

// MetricRow модель сообщения.
type MertricRow struct {
	Date time.Time
	Body interface{}
	Name string
}

// ScoutStorage тип хранилища данных, которые собрал скаут.
type ScoutStorage interface {
	Save(row MertricRow) bool
	GetAvgByRange(r time.Duration) (MertricRow, error)
}
