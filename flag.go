package main

import (
	"errors"
	"github.com/major1201/goutils"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

func getApp() *cli.App {
	app := cli.NewApp()
	app.Name = Name
	app.HelpName = app.Name
	app.Usage = app.Name
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "show help",
		},
		cli.VersionFlag,
		cli.StringFlag{
			Name:  "config, c",
			Usage: "config file path, default(/etc/k8s-mutator/config.yml)",
			Value: "/etc/k8s-mutator/config.yml",
		},
		cli.StringFlag{
			Name:   "tls-cert-file",
			Usage:  "File containing the default x509 Certificate for HTTPS",
			EnvVar: "TLS_CERT_FILE",
		},
		cli.StringFlag{
			Name:   "tls-private-key-file",
			Usage:  "File containing the default x509 private key matching --tls-cert-file.",
			EnvVar: "TLS_PRIVATE_KEY_FILE",
		},
		cli.StringFlag{
			Name:  "listen",
			Usage: "listen address",
			Value: ":443",
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.Bool("help") {
			cli.ShowAppHelpAndExit(c, 0)
		}
		// check additional config
		if err := checkConfig(c); err != nil {
			zap.L().Named("Config").Fatal("config error", zap.Error(err))
		}
		runApp(c)
		return nil
	}
	app.HideHelp = true
	return app
}

func checkConfig(c *cli.Context) error {
	if goutils.IsEmpty(c.String("tls-cert-file")) {
		return errors.New("tls-cert-file is empty")
	}
	if goutils.IsEmpty(c.String("tls-private-key-file")) {
		return errors.New("tls-private-key-file is empty")
	}
	return nil
}
