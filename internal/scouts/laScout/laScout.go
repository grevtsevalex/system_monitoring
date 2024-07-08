package lascout

import (
	"context"
	"fmt"

	"github.com/grevtsevalex/system_monitoring/internal/scouts"
)

type laScout struct {
	status scouts.StatusID
	ctx    context.Context
	logger scouts.Logger
}

func NewLoadAverageScout(ctx context.Context, logg scouts.Logger) *laScout {
	return &laScout{ctx: ctx, status: scouts.StatusIDSleeping, logger: logg}
}

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
		}
	}()
	return nil
}

func (l *laScout) Stop() error {

}

func (l *laScout) Status() scouts.StatusID {
	return l.status
}
