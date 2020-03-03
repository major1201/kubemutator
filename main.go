package main

import (
	"github.com/major1201/kubemutator/internal/config"
	"github.com/major1201/kubemutator/internal/view"
	"github.com/major1201/kubemutator/pkg/log"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"os"
)

// Name inspects the project name
var Name = "kubemutator"

// Version means the project's version
var Version = "custom"

func init() {
	// start program
	zap.L().Info("starting up", zap.String("name", Name), zap.String("version", Version))
}

func runMain(c *cli.Context) error {
	// load config file
	config.SetPath(c.String("config"))
	if err := config.LoadConfig(); err != nil {
		zap.L().Named("config").Fatal("error loading config file", log.Error(err))
	}

	// configmap watcher
	//ctx, cancelContexts := context.WithCancel(context.Background())
	//go func() {
	//	configWatcher, err := watcher.New()
	//	if err != nil {
	//		zap.L().Named("watcher").Error("watcher start failed", zap.Error(err))
	//	}
	//
	//	sigChan := make(chan interface{}, 10)
	//
	//	// watcher restarter
	//	go func() {
	//		out:
	//		for {
	//			err := configWatcher.Watch(ctx, sigChan)
	//			switch err {
	//			case watcher.WatchChannelClosedError:
	//				zap.L().Named("watcher").Error("watcher got error, try to restart watcher", zap.Error(err))
	//			default:
	//				zap.L().Named("watcher").Error("unknown watcher error", zap.Error(err))
	//				break out
	//			}
	//		}
	//	}()
	//
	//	// reloader
	//	for {
	//		select {
	//		case <-sigChan:
	//			_ = config.LoadConfig()
	//		}
	//	}
	//}()

	// start serving
	view.ServeHTTP(c.String("listen"), view.ConfigTLS(c.String("tls-cert-file"), c.String("tls-private-key-file")))
	//cancelContexts()

	return nil
}

func main() {
	// parse flags
	if err := getCLIApp().Run(os.Args); err != nil {
		zap.L().Fatal("flag unexpected error", zap.Error(err))
	}
}
