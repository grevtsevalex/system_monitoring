package lascout

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/scouts"
)

// laScout тип скаута сбора LoadAverage.
type laScout struct {
	status   scouts.StatusID
	ctx      context.Context
	cancelFn context.CancelFunc
	logger   scouts.Logger
}

// NewLoadAverageScout конструктор скаута.
func NewLoadAverageScout(ctx context.Context, logg scouts.Logger) *laScout {
	newCtx, cancelfn := context.WithCancel(ctx)
	return &laScout{ctx: newCtx, cancelFn: cancelfn, status: scouts.StatusIDSleeping, logger: logg}
}

// Run запуск скаута.
func (l *laScout) Run() error {
	l.status = scouts.StatusIDPending
	go func() {
		defer func() {
			l.status = scouts.StatusIDStopping
		}()
		l.status = scouts.StatusIDRunning
		for {
			select {
			case <-l.ctx.Done():
				l.logger.Error(fmt.Sprintf("la scout stopping by context: %s", l.ctx.Err().Error()))
				return
			default:
			}

			// call linux fn to get info
			// write to storage or aggregator
			cmd := exec.Command("uptime")
			result, err := cmd.Output()
			if err != nil {
				l.logger.Error(fmt.Sprintf("calling uptime: %s", err.Error()))
			}

			l.logger.Log(string(result))

			time.Sleep(time.Second * 1)
		}
	}()
	return nil
}

// Stop остановка скаута.
func (l *laScout) Stop() error {
	l.cancelFn()
	return nil
}

// Status получение статуса скаута.
func (l *laScout) Status() scouts.StatusID {
	return l.status
}
