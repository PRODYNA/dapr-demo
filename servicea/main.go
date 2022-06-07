package main

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
	muxlogrus "github.com/pytimer/mux-logrus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

var (
	listenAddress = getenv("LISTEN_ADDRESS", "0.0.0.0:8000")
	pubsubName    = getenv("PUB_SUB_NAME", "default")
	topicName     = getenv("TOPIC_NAME", "topic")
)

const (
	data = "The order is in"
)

func main() {
	log.Info("Starting app")
	r := mux.NewRouter()
	r.HandleFunc("/checkout", ScheduleHandler).Methods("POST", "OPTIONS")
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

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// Handler for /schedule, sends a message
func ScheduleHandler(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL.Path).Info("Schedule triggered")
	if r.Method == "POST" {
		err := sendMessage()
		if err != nil {
			log.WithError(err).Warn("Unable to send message")
			w.WriteHeader(500)
			w.Write([]byte("error"))
			return
		}
	}
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL.Path).Info("Health triggered")
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

// Send a message to pubsub
func sendMessage() error {
	client, err := dapr.NewClient()
	if err != nil {
		log.WithError(err).Error("Unable to create DAPR client")
		return err
	}
	defer client.Close()

	ctx := context.Background()
	err = client.PublishEvent(ctx, pubsubName, topicName, []byte(data))
	if err != nil {
		log.WithError(err).WithFields(
			log.Fields{
				"pubsubName": pubsubName,
				"topicName":  topicName,
				"data":       data,
			}).Error("Unable to send message")
		return err
	}
	return nil
}
