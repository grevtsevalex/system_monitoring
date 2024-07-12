package cpuScout

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
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

// MetricName название метрики.
const MetricName = "CPU"

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
	str := strings.ReplaceAll(data, ",", ".")
	values := strings.Split(str, " ")

	usr, err := strconv.ParseFloat(values[0], 32)
	if err != nil {
		l.logger.Error(fmt.Sprintf("parse float: %s", err.Error()))
	}

	sys, err := strconv.ParseFloat(values[1], 32)
	if err != nil {
		l.logger.Error(fmt.Sprintf("parse float: %s", err.Error()))
	}

	idl, err := strconv.ParseFloat(values[2], 32)
	if err != nil {
		l.logger.Error(fmt.Sprintf("parse float: %s", err.Error()))
	}

	date := time.Now().UTC()

	cpuData := CpuData{
		Date: date,
		Usr:  float32(usr),
		Sys:  float32(sys),
		Idle: float32(idl),
		Name: MetricName}

	l.storage.Save(scouts.MertricRow{Name: MetricName, Date: date, Body: cpuData})
}
