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
	listenAddress = "0.0.0.0:8000" // listen address
	pubsubName = "demoqueue"
	topicName  = "demoqueue"
	data       = "Hello"
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

// Handler for /schedule, sends a message
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

// Send a message to pubsub
func sendMessage() error {
	client,err := dapr.NewClient()
	if err != nil {
		log.WithError(err).Error("Unable to create DAPR client")
		return err
	}
	defer client.Close()
	ctx := context.Background()

	in := &dapr.InvokeBindingRequest{
		Name:      "demoqueue",
		Operation: "create",
		Data: []byte(data),
		Metadata: map[string]string{"k1": "v1", "k2": "v2"},
	}

	_, err = client.InvokeBinding(ctx, in)
	if err != nil {
		log.WithError(err).WithFields(
			log.Fields{
				"pubsubName": pubsubName,
				"topicName": topicName,
				"data": data,
			}).Error("Unable to send message")
		return err
	}
	return nil
}