package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func indexHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "INFLUXDB")
	if session.Values["INFLUXDB_AUTH"] == nil {
		http.Redirect(response, request, "/login", http.StatusSeeOther)
		return
	}

	file, err := os.Open("../dist/index.html")
	if err != nil {
		log.Printf("error in request %q: %v\n", request.URL, err.Error())
		http.Error(response, "server error", 500)
		return
	}
	http.ServeContent(response, request, "index.html", time.Time{}, file)
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	file, err := os.Open("../dist/login.html")
	if err != nil {
		log.Printf("error in request %q: %v\n", request.URL, err.Error())
		http.Error(response, "server error", 500)
		return
	}
	http.ServeContent(response, request, "login.html", time.Time{}, file)
}
