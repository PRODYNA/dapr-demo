package main

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
	muxlogrus "github.com/pytimer/mux-logrus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

const (
	listenAddress = "0.0.0.0:8000" // listen address
)

var (
	stateStoreName = "order"
	stateName      = "orderNumber"
)

func main() {
	log.Info("Starting app")
	r := mux.NewRouter()
	r.HandleFunc("/health", HealthHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/number", NumberHandler).Methods("POST", "OPTIONS")
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
	log.WithField("url", r.URL.Path).Info("Health triggered")
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

// Generates a new number by incrementing value in state
func NumberHandler(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL.Path).Info("Schedule triggered")
	if r.Method == "GET" {
	}
	w.WriteHeader(200)
	w.Write([]byte("ok"))
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
	ctx := context.Background()
	result, err := client.GetState(ctx, stateStoreName, stateName, nil)
	if err != nil {
		log.WithError(err).Warn("Unable to read state")
		number = 1
	} else {
		number, err = strconv.Atoi(string(result.Value))
	}
	number++
	log.Infof("New order number %i", number)

	err = client.SaveState(ctx, stateStoreName, stateName, []byte(strconv.Itoa(number)), nil)
	if err != nil {
		log.WithError(err).Error("Unable to save state")
		return 0, err
	}

	return number, nil
}
