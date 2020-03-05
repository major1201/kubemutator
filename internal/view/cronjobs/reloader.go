package cronjobs

import (
	"github.com/major1201/kubemutator/internal/config"
	"github.com/major1201/kubemutator/pkg/cronjob"
	"github.com/major1201/kubemutator/pkg/log"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func init() {
	cr := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	_, err := cr.AddFunc("*/20 * * * * *", func() {
		ReloadConfig()
	})
	if err != nil {
		zap.L().Named("cronjob").Fatal("reload config failed")
	}

	cronjob.Register(cr)
}

// ReloadConfig reload the config every 20 seconds
func ReloadConfig() {
	if err := config.LoadConfig(); err != nil {
		zap.L().Error("reload config error", log.Error(err))
	}
}
