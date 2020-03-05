package cronjob

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type allStatus int

const (
	statusStopped allStatus = iota
	statusStarted
)

var status allStatus

var cronList []*cron.Cron

// Register a cron to the cron job pool
func Register(cr *cron.Cron) {
	cronList = append(cronList, cr)
	if status == statusStarted {
		cr.Start()
	} else {
		cr.Stop()
	}
}

// Start cron jobs
func Start() {
	zap.L().Named("cronjob").Info("starting cron jobs")
	if status == statusStarted {
		return
	}

	for _, cr := range cronList {
		cr.Start()
	}
}

// Stop cron jobs
func Stop() {
	zap.L().Named("cronjob").Info("stopping cron jobs")
	if status == statusStopped {
		return
	}

	for _, cr := range cronList {
		cr.Stop()
	}
}
