package view

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/major1201/kubemutator/internal/view/mutate"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// SetRouter sets the main http route
func SetRouter(router *mux.Router) {
	// prometheus metrics
	router.Handle("/metrics", promhttp.Handler())

	// mutate
	mutateRoute := router.Path("/mutate").Subrouter()
	mutateRoute.Use(RequestIDMiddleware, LogMiddleware)
	mutateHandler := handlers.ContentTypeHandler(http.HandlerFunc(mutate.ServeMutate), "application/json")
	mutateRoute.Handle("", mutateHandler)
}
