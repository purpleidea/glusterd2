// Package rest implements the REST server for GlusterD
package rest

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	config "github.com/spf13/viper"
	"gopkg.in/tylerb/graceful.v1"
)

// GDRest is the GlusterD Rest server
type GDRest struct {
	Routes *mux.Router

	addr string
	srv  *graceful.Server
}

// New returns a GDRest object which can listen on the configured address
func New() *GDRest {
	rest := &GDRest{}

	rest.addr = config.GetString("restaddress")

	rest.Routes = mux.NewRouter()

	n := negroni.New()
	n.UseHandler(rest.Routes)

	rest.srv = &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    rest.addr,
			Handler: n,
		},
	}

	return rest
}

// Listen begins the GlusterD Rest server
func (r *GDRest) Listen() error {
	log.WithField("ip:port", r.addr).Info("Starting GlusterD REST server")
	err := r.srv.ListenAndServe()
	if err != nil {
		log.WithField("error", err).Error("Failed to start the Rest Server")
		return err
	}
	return nil
}

// Stop stops the GlusterD Rest server
func (r *GDRest) Stop() {
	log.Debug("Stopping the GlusterD Rest server")
	schan := r.srv.StopChan()
	r.srv.Stop(10 * time.Second)
	<-schan
	log.Info("Stopped the GlusterD Rest Server")

	return
}
