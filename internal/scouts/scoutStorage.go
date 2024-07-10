package scouts

import (
	"fmt"
	"time"
)

// MetricRow модель сообщения.
type MertricRow struct {
	Date time.Time
	Body string
	Name string
}

// ScoutStorage тип хранилища данных, которые собрал скаут.
type ScoutStorage interface {
	Save(row MertricRow) bool
	GetAvgByRange(r time.Duration) (MertricRow, error)
}

// String приведение к строке.
func (r *MertricRow) String() string {
	t := r.Date
	formatted := fmt.Sprintf("%02d/%s/%d:%02d:%02d:%02d",
		t.Day(), t.Month().String(), t.Year(),
		t.Hour(), t.Minute(), t.Second())
	return fmt.Sprintf("[%s]: <%s> %s", formatted, r.Name, r.Body)
}
