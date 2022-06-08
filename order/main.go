package main

import (
	"fmt"
	"github.com/gorilla/mux"
	muxlogrus "github.com/pytimer/mux-logrus"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	listenAddress = "0.0.0.0:8000" // listen address
)

func main() {
	log.Info("Starting app")
	r := mux.NewRouter()
	r.HandleFunc("/order", OrderHandler).Methods("GET", "POST", "OPTIONS")
	r.HandleFunc("/health", HealthHandler).Methods("GET", "OPTIONS")
	http.Handle("/", r)
	r.Use(muxlogrus.NewLogger().Middleware)

	srv := &http.Server{
		Handler:      r,
		Addr:         listenAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.WithField("listenAddress", listenAddress).Info("Starting listener")
	log.Fatal(srv.ListenAndServe())
	log.Info("Stopping listener")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL.Path).Trace("Health triggered")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Handler for /schedule, sends a message
func OrderHandler(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL.Path).WithField("method", r.Method).Info("Order triggered")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Unable to read body")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	log.Info("Order received %s", body)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Order received %s", body)))
}
