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

const (
	listenAddress = "0.0.0.0:8000" // listen address
)

var (
	stateStoreName = getenv("STATE_STORE_NAME", "order")
	stateName      = getenv("STATE_NAME", "orderNumber")
)

func main() {
	log.Info("Starting app")
	r := mux.NewRouter()
	r.HandleFunc("/health", HealthHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/number", NumberHandler).Methods("GET")
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
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

// Generates a new number by incrementing value in state
func NumberHandler(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL.Path).Info("Schedule triggered")
	number, err := getNumber()
	if err != nil {
		log.WithError(err).Error("Unable to read the number")
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(200)
	numberString := strconv.Itoa(number)

	w.Write([]byte(numberString))
}

// Get the number from the state, increment it and store it back to the state
// If there is no number in the state, initialize it with 1
func getNumber() (number int, err error) {
	number = 0
	client, err := dapr.NewClient()
	if err != nil {
		log.Error("Unable to create DAPR client")
		return 0, err
	}
	// defer client.Close()

	ctx := context.Background()

	result, err := client.GetState(ctx, stateStoreName, stateName, nil)
	if err != nil {
		log.WithError(err).WithField("stateStoreName", stateStoreName).WithField("stateName", stateName).Warn("Unable to read state")
		// don't go out here, this might be the first call to the state, so init the number and save a new state
	} else {
		number, err = strconv.Atoi(string(result.Value))
	}
	number++
	log.Infof("New order number %d", number)

	err = client.SaveState(ctx, stateStoreName, stateName, []byte(strconv.Itoa(number)), nil)
	if err != nil {
		log.WithError(err).WithField("stateStoreName", stateStoreName).WithField("stateName", stateName).Error("Unable to save state")
		return 0, err
	}

	return number, nil
}

// Read environment value with default
func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
