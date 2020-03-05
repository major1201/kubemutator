package main

import (
	"github.com/major1201/kubemutator/internal/config"
	"github.com/major1201/kubemutator/internal/view"
	_ "github.com/major1201/kubemutator/internal/view/cronjobs"
	"github.com/major1201/kubemutator/pkg/cronjob"
	"github.com/major1201/kubemutator/pkg/httputils"
	"github.com/major1201/kubemutator/pkg/log"
	"github.com/major1201/kubemutator/pkg/signal"
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
	httpServer := httputils.HTTPServer{
		Addr:                    c.String("listen"),
		TLSConfig:               httputils.ConfigTLS(c.String("tls-cert-file"), c.String("tls-private-key-file")),
		RouterHandler:           view.Router,
		EnablePrometheusMetrics: true,
		EnablePProf:             true,
	}
	zap.L().Info("starting http server", zap.String("listen", httpServer.Addr))
	if err := httpServer.ListenAndServe(); err != nil {
		zap.L().Error("serve http server error", log.Error(err))
	}

	return nil
}

func main() {
	go signal.Serve()

	if err := getCLIApp().Run(os.Args); err != nil {
		zap.L().Fatal("flag unexpected error", zap.Error(err))
	}
}
