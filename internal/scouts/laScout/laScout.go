package lascout

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/scouts"
)

// laScout тип скаута сбора LoadAverage.
type laScout struct {
	status   scouts.StatusID
	ctx      context.Context
	cancelFn context.CancelFunc
	logger   scouts.Logger
	storage  *Storage
}

// NewLoadAverageScout конструктор скаута.
func NewLoadAverageScout(ctx context.Context, logg scouts.Logger, st *Storage) *laScout {
	newCtx, cancelfn := context.WithCancel(ctx)
	return &laScout{ctx: newCtx, cancelFn: cancelfn, status: scouts.StatusIDSleeping, logger: logg, storage: st}
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

			cmd := exec.Command("uptime")
			result, err := cmd.Output()
			if err != nil {
				l.logger.Error(fmt.Sprintf("calling uptime: %s", err.Error()))
			}

			loadAveragesValues := strings.Trim(string(result[len(result)-17:]), "\n")
			l.logger.Log(loadAveragesValues)
			l.write(loadAveragesValues)

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

// write записать данные в хранилище.
func (l *laScout) write(data string) {
	l.storage.Save(scouts.MertricRow{Date: time.Now().UTC(), Body: data, Name: "Load Average"})
}
