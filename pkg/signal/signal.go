package signal

import (
	"github.com/major1201/goutils"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Handler defines a signal handler struct
type Handler struct {
	Name string
	Func func(sig os.Signal) bool
}

type regItem struct {
	Signal  os.Signal
	Handler Handler
}

var (
	registry []regItem
	sigCh    chan os.Signal
)

// Register a signal handler to the registry
func Register(sig os.Signal, handler Handler) {
	registry = append(registry, regItem{
		Signal:  sig,
		Handler: handler,
	})
}

/*
Serve the signal notifier
	Should be run as: go signal.Serve()

	Currently support signals: SIGHUP, SIGINT, SIGQUIT, SIGABRT, SIGALRM, SIGTERM, SIGUSR1, SIGUSR2.

	The Handler function should returns a bool value, which means if the function accept the termination.
	If all got by the handler function returns "true" and the signal is one of SIGINT, SIGQUIT, SIGTERM,
	the program would exit with code 0.
*/
func Serve() {
	sigCh = make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGALRM, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

	for {
		time.Sleep(time.Second)
		sig := <-sigCh

		logger := zap.L().Named("signal").With(zap.String("signal", sig.String()))
		if s, ok := sig.(syscall.Signal); ok {
			logger = logger.With(zap.Int("no", int(s)))
		}
		logger.Info("signal received")

		allTrue := true
		for _, reg := range registry {
			if sig == reg.Signal {
				logger.Warn("running callback func", zap.String("name", reg.Handler.Name))
				allTrue = allTrue && reg.Handler.Func(sig)
			}
		}

		if allTrue && isTermSignal(sig) {
			zap.L().Named("signal").Warn("exiting")
			os.Exit(0)
		}
	}
}

func isTermSignal(sig os.Signal) bool {
	return goutils.Contains(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL)
}
