package main

// This file should only include adding handlers and calling listen and serve
// handler implementations should be seperated to their own files or put in files of related handlers

import (
	"context"
	"crypto/rand"
	b64 "encoding/base64"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	sessions "github.com/gorilla/sessions"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

var store *sessions.CookieStore

// stores all active clients
var clients map[string]influxdb2.Client

// creates the session key and sets the enviroment variable
func createSessionKey() {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Printf("Failed to Generate Key\n")
		return
	}
	err = os.Setenv("SESSION_KEY", b64.StdEncoding.EncodeToString(key))
	if err != nil {
		log.Printf("Session Key Failed to be Set\n")
		return
	}
	log.Printf("Session key: %v\n", os.Getenv("SESSION_KEY"))
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

func main() {
	clients = make(map[string]influxdb2.Client)

	createSessionKey()
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../dist/assets/"))))

	mux.HandleFunc("/attempt-login", login)
	mux.HandleFunc("/attempt-logout", logout)

	mux.HandleFunc("/get-uid", getUIDHandler) // for the UID for grafana

	mux.HandleFunc("/generate-query", queryGenerationHandler)
	mux.HandleFunc("/run-query", runQueryHandler)

	mux.HandleFunc("/get-buckets", getBucketsHandler)
	mux.HandleFunc("/get-fields", getFieldsHandler)
	mux.HandleFunc("/get-measurements", getMeasurementsHandler)
	mux.HandleFunc("/favicon.ico", func(http.ResponseWriter, *http.Request) {})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// listen and serve
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Server error: %v\n", err.Error())
			interrupt <- syscall.SIGTERM
		}

	}()

	log.Printf("Server started on port %s\n", server.Addr[1:])

	// wait for kill signal
	<-interrupt
	log.Println("Shutting down server")

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(context); err != nil {
		log.Fatalf("Server shutdown forced: %v", err)
	}

	log.Printf("Server shut down")
}
