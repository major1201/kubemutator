package view

import (
	"crypto/tls"
	"github.com/gorilla/mux"
	plog "github.com/major1201/kubemutator/pkg/log"
	"go.uber.org/zap"
	"net/http"
)

var _log *zap.Logger

func log() *zap.Logger {
	if _log == nil {
		_log = zap.L().Named("api")
	}
	return _log
}

// ConfigTLS configures the TCP TLS config with certificate and private key file
func ConfigTLS(certFile, keyFile string) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		zap.L().Named("http").Fatal("load x509 key pair error", zap.Error(err))
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
		// TODO: uses mutual tls after we agree on what cert the apiserver should use.
		// ClientAuth:   tls.RequireAndVerifyClientCert,
	}
}

// ServeHTTP just serve the HTTP requests
func ServeHTTP(listenAddress string, tlsConfig *tls.Config) {
	router := mux.NewRouter()

	SetRouter(router)

	http.Handle("/", router)
	server := &http.Server{
		Addr:      listenAddress,
		TLSConfig: tlsConfig,
	}
	log().Info("starting http server", zap.String("listen", listenAddress))
	log().Fatal("http server ends", plog.Error(server.ListenAndServeTLS("", "")))
}
