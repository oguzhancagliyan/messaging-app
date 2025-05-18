package scheduler

import (
	"context"
	"messaging-app/internal/service"
	"time"

	"go.uber.org/zap"
)

type Dispatcher interface {
	Start()
	Stop()
}

type dispatcher struct {
	ticker    *time.Ticker
	stopChan  chan struct{}
	service   service.MessageService
	isRunning bool
}

func NewDispatcher(svc service.MessageService) *dispatcher {
	return &dispatcher{
		service:  svc,
		ticker:   nil,
		stopChan: make(chan struct{}),
	}
}

func (d *dispatcher) Start() {
	if d.isRunning {
		zap.L().Warn("Dispatcher already running")
		return
	}

	d.ticker = time.NewTicker(2 * time.Minute)
	d.isRunning = true

	go func() {
		zap.L().Info("Dispatcher started. Will run every 2 minutes.")
		for {
			select {
			case <-d.ticker.C:
				zap.L().Info("Dispatcher tick triggered")
				err := d.service.SendUnsentMessages(context.Background())
				if err != nil {
					zap.L().Error("Dispatcher failed to send messages", zap.Error(err))
				}
			case <-d.stopChan:
				zap.L().Info("Dispatcher stopping...")
				d.ticker.Stop()
				d.isRunning = false
				return
			}
		}
	}()
}

func (d *dispatcher) Stop() {
	if !d.isRunning {
		zap.L().Warn("Dispatcher is not running")
		return
	}
	d.stopChan <- struct{}{}
}
