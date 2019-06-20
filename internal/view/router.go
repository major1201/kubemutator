package view

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/major1201/k8s-mutator/internal/view/mutate"
	"github.com/major1201/k8s-mutator/internal/view/reload"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

// SetRouter sets the main http route
func SetRouter(router *mux.Router) {
	// prometheus metrics
	router.Handle("/metrics", promhttp.Handler())

	// mutate
	mutateHandler := handlers.ContentTypeHandler(http.HandlerFunc(mutate.ServeMutate), "application/json")
	mutateHandler = handlers.LoggingHandler(os.Stdout, mutateHandler)
	router.Handle("/mutate", mutateHandler)

	// reload config
	router.HandleFunc("/reload", reload.ServeReload)
}
