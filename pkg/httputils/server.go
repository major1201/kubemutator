package httputils

import (
	"crypto/tls"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

// HTTPServer defines an HTTP server
type HTTPServer struct {
	Addr          string
	TLSConfig     *tls.Config
	RouterHandler func(router *gin.Engine)

	EnableCompress          bool
	EnableProxyHeaders      bool
	EnablePrometheusMetrics bool
	EnablePProf             bool
	EnableCORS              bool
}

func defaultRouterFunc(router *gin.Engine) {
	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(zap.L(), true))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})
}

func (hs *HTTPServer) getServer() (server *http.Server) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	serveMux := http.NewServeMux()
	// CORS
	if hs.EnableCORS {
		router.Use(cors.Default())
	}
	// compress
	if hs.EnableCompress {
		router.Use(gzip.Gzip(gzip.DefaultCompression))
	}
	// proxy headers
	if hs.EnableProxyHeaders {
		handlers.ProxyHeaders(serveMux)
	}
	if hs.RouterHandler != nil {
		hs.RouterHandler(router)
	} else {
		defaultRouterFunc(router)
	}

	// prometheus metrics
	if hs.EnablePrometheusMetrics {
		router.Any("/metrics", gin.WrapH(promhttp.Handler()))
	}

	// pprof
	if hs.EnablePProf {
		pprof.Register(router)
	}

	serveMux.Handle("/", router)

	server = &http.Server{
		Addr:      hs.Addr,
		Handler:   serveMux,
		TLSConfig: hs.TLSConfig,
	}
	return
}

// ListenAndServe serves the HTTP requests and listens a net interface
func (hs *HTTPServer) ListenAndServe() error {
	server := hs.getServer()
	if hs.TLSConfig == nil {
		return server.ListenAndServe()
	}
	return server.ListenAndServeTLS("", "")
}

// Serve just serves the HTTP requests
func (hs *HTTPServer) Serve(l net.Listener) error {
	server := hs.getServer()
	if hs.TLSConfig == nil {
		return server.Serve(l)
	}
	return server.ServeTLS(l, "", "")
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
