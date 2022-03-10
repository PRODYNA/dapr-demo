package main

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
	muxlogrus "github.com/pytimer/mux-logrus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	listenAddress = "0.0.0.0:8000"
)

func main() {
	log.Info("Starting app")
	r := mux.NewRouter()
	r.HandleFunc("/schedule", ScheduleHandler).Methods("POST","OPTIONS")
	http.Handle("/", r)
	r.Use(muxlogrus.NewLogger().Middleware)

	srv := &http.Server{
		Handler: r,
		Addr:    listenAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.WithField("listenAddress", listenAddress ).Info("Starting listener")
	log.Fatal(srv.ListenAndServe())
	log.Info("Stopping listener")
}

func ScheduleHandler(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL.Path ).Info("Schedule triggered")
	if r.Method == "POST" {
		err := sendMessage()
		if err != nil {
			log.WithError(err).Warn("Unable to send message")
			w.WriteHeader(500)
			w.Write([]byte("error"))
			return
		}
	}
	w.WriteHeader( 200 )
	w.Write( []byte( "ok" ) )
}

func sendMessage() error {
	client,err := dapr.NewClient()
	if err != nil {
		log.WithError(err).Error("Unable to create client")
		return err
	}
	defer client.Close()
	ctx := context.Background()
	err = client.PublishEvent(ctx, "broker", "demo", "Hello")
	if err != nil {
		log.WithError(err).Error("Unable to send message")
		return err
	}
	return nil
}
