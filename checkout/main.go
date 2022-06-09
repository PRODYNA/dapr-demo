package main

import (
	"context"
	"fmt"
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
	listenAddress = getenv("LISTEN_ADDRESS", "0.0.0.0:8000") // the port we listen on, with default
	pubsubName    = getenv("PUB_SUB_NAME", "default")        // the pubsub we are using, with default
	topicName     = getenv("TOPIC_NAME", "topic")            // the topic we are sending to, with default
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

// Used for Kubernetes probes
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL.Path).Trace("Health triggered")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Handler for /checkout, get a order number and send a message
func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL.Path).Info("Schedule triggered")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// get order number
	number, err := getOrderNumber()
	if err != nil {
		log.WithError(err).Error("Unable to get order number")
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	// send message
	err = sendMessage(number)
	if err != nil {
		log.WithError(err).Warn("Unable to send message")
		w.WriteHeader(http.StatusPreconditionFailed)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Order %d", number)))
}

// Get order number by calling the service
func getOrderNumber() (number int, err error) {
	client, err := dapr.NewClient()
	if err != nil {
		log.WithError(err).Error("Unable to create DAPR client")
		return 0, err
	}
	// defer client.Close()

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
	// defer client.Close()

	ctx := context.Background()
	defer ctx.Done()

	data := []byte(strconv.Itoa(number))
	llog := log.WithFields(
		log.Fields{
			"pubsubName": pubsubName,
			"topicName":  topicName,
			"data":       data,
		})

	err = client.PublishEvent(ctx, pubsubName, topicName, data)
	if err != nil {
		llog.WithError(err).Error("Unable to send message")
		return err
	}
	llog.Info("Message sent")
	return nil
}
