package discscout

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/scouts"
)

// discScout тип скаута сбора метрик загрузки дисков.
type discScout struct {
	status   scouts.StatusID
	ctx      context.Context
	cancelFn context.CancelFunc
	logger   scouts.Logger
	storage  *Storage
}

// MetricName название метрики.
const MetricName = "DISC"

// NewDiscScout конструктор скаута.
func NewDiscScout(ctx context.Context, logg scouts.Logger, st *Storage) *discScout {
	newCtx, cancelfn := context.WithCancel(ctx)
	return &discScout{ctx: newCtx, cancelFn: cancelfn, status: scouts.StatusIDSleeping, logger: logg, storage: st}
}

// Run запуск скаута.
func (l *discScout) Run() error {
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

			cmdLine := "iostat -dk | grep ',' | awk '{print $1, $2, $3, $4}'"
			cmd := exec.Command("bash", "-c", cmdLine)
			result, err := cmd.Output()
			if err != nil {
				l.logger.Error(fmt.Sprintf("calling iostat -dk: %s", err.Error()))
			}

			// cpuValues := strings.Trim(string(result), "\n")
			l.write(string(result))

			time.Sleep(time.Second * 1)
		}
	}()
	return nil
}

// Stop остановка скаута.
func (l *discScout) Stop() error {
	l.cancelFn()
	return nil
}

// Status получение статуса скаута.
func (l *discScout) Status() scouts.StatusID {
	return l.status
}

// write записать данные в хранилище.
func (l *discScout) write(data string) {
	discData := DiscData{}
	reader := strings.NewReader(data)

	sc := bufio.NewScanner(reader)
	for sc.Scan() {
		discValues := sc.Text()
		str := strings.ReplaceAll(discValues, ",", ".")
		values := strings.Split(str, " ")

		discName := values[0]

		tps, err := strconv.ParseFloat(values[1], 32)
		if err != nil {
			l.logger.Error(fmt.Sprintf("parse float: %s", err.Error()))
		}

		rps, err := strconv.ParseFloat(values[2], 32)
		if err != nil {
			l.logger.Error(fmt.Sprintf("parse float: %s", err.Error()))
		}

		wps, err := strconv.ParseFloat(values[3], 32)
		if err != nil {
			l.logger.Error(fmt.Sprintf("parse float: %s", err.Error()))
		}
		device := Device{
			Name: discName,
			Tps:  float32(tps),
			Rps:  float32(rps),
			Wps:  float32(wps)}

		discData.Devices = append(discData.Devices, device)
	}

	if err := sc.Err(); err != nil {
		l.logger.Error(fmt.Sprintf("scanning response: %s", err.Error()))
	}

	date := time.Now().UTC()
	discData.Date = date
	discData.Name = MetricName

	l.storage.Save(scouts.MertricRow{Name: MetricName, Date: date, Body: discData})
}
