package reload

import (
	"github.com/major1201/k8s-mutator/internal/config"
	"github.com/major1201/k8s-mutator/pkg/httputils"
	"github.com/major1201/k8s-mutator/pkg/log"
	"go.uber.org/zap"
	"net/http"
)

// ServeReload serves the /reload path
func ServeReload(w http.ResponseWriter, r *http.Request) {
	err := config.LoadConfig()
	if err != nil {
		zap.L().Named("config").Error("error loading config file", log.Error(err))
		httputils.WriteJSONWithCode(w, r, map[string]string{"message": err.Error()}, http.StatusInternalServerError)
	} else {
		httputils.WriteJSON(w, r, map[string]string{})
	}
}
