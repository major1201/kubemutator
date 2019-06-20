package main

import (
	"github.com/major1201/k8s-mutator/internal/config"
	"github.com/major1201/k8s-mutator/internal/view"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// AppVer means the project's version
const AppVer = "0.1.0-r1"

func initLog(stdout, stderr string, level zapcore.Level) {
	zap.NewProductionConfig()
	logger, _ := zap.Config{
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{stdout},
		ErrorOutputPaths: []string{stderr},
	}.Build()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
}

func runApp(c *cli.Context) {
	// load config file
	config.SetPath(c.String("config"))
	if err := config.LoadConfig(); err != nil {
		zap.L().Named("config").Fatal("error loading config file", zap.Error(err))
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
}

func main() {
	initLog("stdout", "stderr", zapcore.DebugLevel)

	zap.L().Named("system").Info("starting k8s-mutator", zap.String("version", AppVer))

	// parse flags
	if err := getApp().Run(os.Args); err != nil {
		zap.L().Fatal("flag unexpected error", zap.Error(err))
	}
}
