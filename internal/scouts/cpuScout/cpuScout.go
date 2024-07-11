package cpuScout

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/scouts"
)

// cpuScout тип скаута сбора метрик нагрузки CPU.
type cpuScout struct {
	status   scouts.StatusID
	ctx      context.Context
	cancelFn context.CancelFunc
	logger   scouts.Logger
	storage  *Storage
}

// NewCPUScout конструктор скаута.
func NewCPUScout(ctx context.Context, logg scouts.Logger, st *Storage) *cpuScout {
	newCtx, cancelfn := context.WithCancel(ctx)
	return &cpuScout{ctx: newCtx, cancelFn: cancelfn, status: scouts.StatusIDSleeping, logger: logg, storage: st}
}

// Run запуск скаута.
func (l *cpuScout) Run() error {
	l.status = scouts.StatusIDPending
	go func() {
		defer func() {
			l.status = scouts.StatusIDStopping
		}()
		l.status = scouts.StatusIDRunning
		for {
			select {
			case <-l.ctx.Done():
				l.logger.Error(fmt.Sprintf("cpu scout stopping by context: %s", l.ctx.Err().Error()))
				return
			default:
			}

			cmdLine := "iostat -c | grep ',' | awk '{print $1, $3, $6}'"
			cmd := exec.Command("bash", "-c", cmdLine)
			result, err := cmd.Output()
			if err != nil {
				l.logger.Error(fmt.Sprintf("calling iostat: %s", err.Error()))
			}

			cpuValues := strings.Trim(string(result), "\n")
			l.logger.Log(cpuValues)
			l.write(cpuValues)

			time.Sleep(time.Second * 1)
		}
	}()
	return nil
}

// Stop остановка скаута.
func (l *cpuScout) Stop() error {
	l.cancelFn()
	return nil
}

// Status получение статуса скаута.
func (l *cpuScout) Status() scouts.StatusID {
	return l.status
}

// write записать данные в хранилище.
func (l *cpuScout) write(data string) {
	l.storage.Save(scouts.MertricRow{Date: time.Now().UTC(), Body: data, Name: "CPU"})
}
