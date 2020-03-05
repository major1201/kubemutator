package main

import (
	"github.com/major1201/kubemutator/internal/config"
	"github.com/major1201/kubemutator/internal/view"
	_ "github.com/major1201/kubemutator/internal/view/cronjobs"
	"github.com/major1201/kubemutator/pkg/cronjob"
	"github.com/major1201/kubemutator/pkg/log"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"os"
)

var (
	// Name inspects the project name
	Name = "kubemutator"

	// Version means the project's version
	Version = "custom"
)

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

	// start cron job
	cronjob.Start()

	// start serving
	view.ServeHTTP(c.String("listen"), view.ConfigTLS(c.String("tls-cert-file"), c.String("tls-private-key-file")))

	return nil
}

func main() {
	// parse flags
	if err := getCLIApp().Run(os.Args); err != nil {
		zap.L().Fatal("flag unexpected error", zap.Error(err))
	}
}
