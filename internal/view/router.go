package view

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	"github.com/major1201/kubemutator/internal/view/mutate"
	"github.com/major1201/kubemutator/pkg/httputils"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Router routes the API URL
func Router(router *gin.Engine) {
	router.Use(
		httputils.RequestIDMiddlewareFunc,
		httputils.TimeoutMiddleware(10*time.Second),
		httputils.LogMiddleware(zap.L().Named("http.request")),
		httputils.RecoveryMiddleware(zap.L().Named("http.recovery"), true),
	)
	mutateHandler := handlers.ContentTypeHandler(http.HandlerFunc(mutate.ServeMutate), "application/json")
	router.Any("/mutate", gin.WrapH(mutateHandler))
}
