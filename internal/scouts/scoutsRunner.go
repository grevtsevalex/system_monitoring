package scouts

import (
	"fmt"
	"sync"
)

// Scout тип скаута (демона, который собирает инфу о хостовой машине).
type Scout interface {
	Run() error
	Stop() error
	Status() StatusID
}

// StatusID идентификатор статуса скаута.
type StatusID int

const (
	StatusIDSleeping        = iota // statusIDSleeping ожидает запуска.
	StatusIDPending                // statusIDPending запускается.
	StatusIDRunning                // statusIDRunning запущен.
	StatusIDStopping               // statusIDStopping останавливается.
	StatusIDCrashedWitError        // statusIDCrashedWitError упал с ошибкой.
)

// Logger тип логгера.
type Logger interface {
	Log(msg string)
	Debug(msg string)
	Error(msg string)
}

// ScoutsRunner тип раннера скаутов.
type ScoutsRunner interface {
	RunScouts() error
	StopScouts() error
	RegisterScout(name string, scout Scout)
	getScouts() map[string]Scout
}

// scoutsRunner модель раннера скаутов.
type scoutsRunner struct {
	mu     sync.RWMutex
	scouts map[string]Scout
	logger Logger
}

// NewScoutsRunner конструктор раннера скаутов.
func NewScoutsRunner(logger Logger) ScoutsRunner {
	return &scoutsRunner{logger: logger, scouts: make(map[string]Scout)}
}

// RunScouts запуск скаутов.
func (s *scoutsRunner) RunScouts() error {
	for name, scout := range s.scouts {
		if scout.Status() != StatusID(StatusIDSleeping) {
			continue
		}

		err := scout.Run()
		if err != nil {
			s.logger.Error(fmt.Sprintf("run scout %s : %s", name, err.Error()))
			return err
		}

		s.logger.Log(fmt.Sprintf("succesed run scout %s", name))
	}

	return nil
}

// RegisterScout добавить скаута в список.
func (s *scoutsRunner) RegisterScout(name string, scout Scout) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scouts[name] = scout
}

// StopScouts остановить скаутов.
func (s *scoutsRunner) StopScouts() error {
	for name, scout := range s.scouts {
		if scout.Status() == StatusIDCrashedWitError && scout.Status() == StatusIDStopping {
			continue
		}

		err := scout.Stop()
		if err != nil {
			s.logger.Error(fmt.Sprintf("stop scout %s : %s", name, err.Error()))
			return err
		}

		s.logger.Log(fmt.Sprintf("succesed stop scout %s", name))
	}

	return nil
}

// getScouts получить список скаутов.
func (s *scoutsRunner) getScouts() map[string]Scout {
	scoutsCopy := make(map[string]Scout, len(s.scouts))
	for k, v := range s.scouts {
		scoutsCopy[k] = v
	}
	return scoutsCopy
}
