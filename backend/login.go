package main

import (
	"context"
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type login_form struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Takes in the http request and adds the appropriate session stuff and creates a logged in client
// returns an error if the login was unsucessful
func login(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "INFLUXDB")

	if session.Values["INFLUXDB_AUTH"] != nil {
		http.Error(response, "Already logged in", http.StatusUnauthorized)
		return
	}

	var data login_form

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(response, "Unable to read body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(response, "Unable to parse body", http.StatusBadRequest)
		log.Printf("error parsing json: %v\n", err.Error())
		return
	}
	username := data.Username
	password := data.Password

	influxURL := os.Getenv("INFLUXDB_URL")
	client := influxdb2.NewClient(influxURL, "")

	err = client.UsersAPI().SignIn(context.Background(), username, password)
	if err != nil {
		log.Printf("Login Failed for user %v: %v\n", username, err.Error())
		http.Error(response, "Bad credentials provided", http.StatusUnauthorized)
		return
	}

	// creates a unique random 32 byte key for the hashtable
	key := make([]byte, 32)
	rand.Read(key)
	b64Key := b64.StdEncoding.EncodeToString(key)
	for clients[b64Key] != nil {
		rand.Read(key)
		b64Key = b64.StdEncoding.EncodeToString(key)
	}

	session.Values["INFLUXDB_AUTH"] = b64Key
	session.Save(request, response)
	clients[b64Key] = client

	response.WriteHeader(http.StatusOK)
}

// logout removes the auth session key and signs out of the client
// error if the logout fails
func logout(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "INFLUXDB")
	if session.Values["INFLUXDB_AUTH"] == nil {
		return
	}

	auth := session.Values["INFLUXDB_AUTH"].(string)
	session.Values["INFLUXDB_AUTH"] = nil
	session.Save(request, response)

	err := (clients[auth]).UsersAPI().SignOut(context.Background())
	if err != nil {
		log.Printf("Couldn't Log out properly %v\n", err.Error())
		http.Error(response, "Bad credentials provided", http.StatusUnauthorized)
		return
	}
	delete(clients, auth)
	response.WriteHeader(http.StatusOK)
}

// returns the client which is associated with the user
func getClient(response http.ResponseWriter, request *http.Request) (influxdb2.Client, error) {
	session, _ := store.Get(request, "INFLUXDB")
	if session.Values["INFLUXDB_AUTH"] == nil {
		http.Redirect(response, request, "/login", http.StatusSeeOther)
		return nil, errors.New("no session saved")
	}
	auth := session.Values["INFLUXDB_AUTH"].(string)
	if clients[auth] == nil {
		session.Values["INFLUXDB_AUTH"] = nil
		session.Save(request, response)
		http.Redirect(response, request, "/login", http.StatusSeeOther)
		return nil, errors.New("no client found")
	}
	return clients[auth], nil
}
