package main

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
	muxlogrus "github.com/pytimer/mux-logrus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
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
	r.HandleFunc("/health", HealthHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/checkout", CheckoutHandler).Methods("POST", "OPTIONS")
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

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL.Path).Info("Health triggered")
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

// Handler for /checkout, sends a message
func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL.Path).Info("Schedule triggered")
	if r.Method == "POST" {

		// get order number
		number, err := getOrderNumber()
		if err != nil {
			log.WithError(err).Error("Unable to get order number")
			w.WriteHeader(500)
			w.Write([]byte("error"))
		}

		// send message
		err = sendMessage(number)
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

// Get order number by calling the service
func getOrderNumber() (number int, err error) {
	client, err := dapr.NewClient()
	if err != nil {
		log.WithError(err).Error("Unable to create DAPR client")
		return 0, err
	}
	defer client.Close()

	ctx := context.Background()
	defer ctx.Done()

	response, err := client.InvokeMethod(ctx, "number", "number", "GET")
	if err != nil {
		log.WithError(err).Warn("Unable to get order number")
		return 0, err
	}

	numberString := string(response)
	number, err = strconv.Atoi(numberString)
	if err != nil {
		log.WithError(err).WithField("numberString", numberString).Warn("Unable to convert atoi")
		return 0, err
	}
	log.WithField("number", number).Info("Received order number")

	return number, nil
}

// Send a message to pubsub
func sendMessage(number int) error {
	client, err := dapr.NewClient()
	if err != nil {
		log.WithError(err).Error("Unable to create DAPR client")
		return err
	}
	defer client.Close()

	ctx := context.Background()
	defer ctx.Done()

	data := []byte(strconv.Itoa(number))
	err = client.PublishEvent(ctx, pubsubName, topicName, data)
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
