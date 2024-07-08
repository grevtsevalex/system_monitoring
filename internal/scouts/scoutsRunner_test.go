package scouts

import (
	"testing"

	"github.com/grevtsevalex/system_monitoring/internal/logger"
	"github.com/stretchr/testify/require"
)

// TestScout тестовый скаут.
type TestScout struct {
	status StatusID
}

// Status получить статус скаута.
func (s *TestScout) Status() StatusID {
	return s.status
}

// Run запустить скаута.
func (s *TestScout) Run() error {
	s.status = StatusIDPending
	return nil
}

// Stop остановить работу скаута.
func (s *TestScout) Stop() error {
	s.status = StatusIDStopping
	return nil
}

// newTestScout конструктор тестового скаута.
func newTestScout() *TestScout {
	return &TestScout{}
}

func TestScoutsRunner(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		scout1 := newTestScout()
		scout2 := newTestScout()
		scout3 := newTestScout()

		logg := logger.New("info")

		sr := NewScoutsRunner(logg)

		sr.RegisterScout("first", scout1)
		sr.RegisterScout("second", scout2)
		sr.RegisterScout("third", scout3)

		err := sr.RunScouts()
		require.NoError(t, err)

		for name, sc := range sr.getScouts() {
			switch name {
			case "first":
				require.True(t, sc.Status() == StatusIDPending || sc.Status() == StatusIDRunning)
			case "second":
				require.True(t, sc.Status() == StatusIDPending || sc.Status() == StatusIDRunning)
			case "third":
				require.True(t, sc.Status() == StatusIDPending || sc.Status() == StatusIDRunning)
			default:
				t.Errorf("undefined scout :%s", name)
			}
		}

		err = sr.StopScouts()
		require.NoError(t, err)

		for name, sc := range sr.getScouts() {
			switch name {
			case "first":
				require.True(t, sc.Status() == StatusIDCrashedWitError || sc.Status() == StatusIDStopping)
			case "second":
				require.True(t, sc.Status() == StatusIDCrashedWitError || sc.Status() == StatusIDStopping)
			case "third":
				require.True(t, sc.Status() == StatusIDCrashedWitError || sc.Status() == StatusIDStopping)
			default:
				t.Errorf("undefined scout :%s", name)
			}
		}
	})
}
