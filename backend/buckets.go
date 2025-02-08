package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func getBucketsHandler(response http.ResponseWriter, request *http.Request) {

	client, _ := getClient(response, request)

	bucketsAPI := client.BucketsAPI()
	buckets, err := bucketsAPI.GetBuckets(context.Background())
	if err != nil {
		log.Printf("error in bucket query: %v", err)
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}
	if buckets == nil {
		log.Println("no buckets")
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}
	output, err := json.Marshal(*buckets)
	if err != nil {
		log.Printf("error in json encoding: %v", err)
		http.Error(response, "internal server error", http.StatusInternalServerError)
	}
	response.Write(output)
}

type responseType struct {
	//Fields       []string `json:"fields"`
	Measurements []string `json:"measurements"`
}
type responseType2 struct {
	Fields []string `json:"fields"`
	//Measurements []string `json:"measurements"`
}

/*
TODO:
  - get bucket from header/body of request
  - check authentication
    *
*/
func getMeasurementsHandler(response http.ResponseWriter, request *http.Request) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Error loading .env file")
	}
	org := os.Getenv("INFLUXDB_ORG")
	if org == "" {
		log.Println("failed to get org")
		return
	}

	client, _ := getClient(response, request)
	queryAPI := client.QueryAPI(org)
	bucket := request.URL.Query().Get("bucket")
	log.Printf("%v\n", bucket)
	measurementsQuery := fmt.Sprintf(`import "influxdata/influxdb/schema" schema.measurements(bucket: "%s")`, bucket)
	result, err := queryAPI.Query(context.Background(), measurementsQuery)

	if err != nil {
		log.Printf("error in labels query: %v", err)
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}
	if result == nil {
		log.Println("no buckets")
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}

	measurements := make(map[string]bool)

	for result.Next() {
		//log.Printf("%v\n", result.Record())
		value := result.Record().ValueByKey("_value")
		measurements[value.(string)] = true
	}
	outputStruct := &responseType{
		make([]string, len(measurements)),
	}
	index := 0
	for measurement := range measurements {
		outputStruct.Measurements[index] = measurement
		index++
	}

	output, err := json.Marshal(outputStruct)

	if err != nil {
		log.Printf("error in json encoding: %v", err)
		http.Error(response, "internal server error", http.StatusInternalServerError)
	}
	log.Printf("%v\n", string(output))
	response.Write(output)

}

func getFieldsHandler(response http.ResponseWriter, request *http.Request) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Error loading .env file")
	}
	org := os.Getenv("INFLUXDB_ORG")
	if org == "" {
		log.Println("failed to get org")
		return
	}

	client, _ := getClient(response, request)

	queryAPI := client.QueryAPI(org)
	bucket := request.URL.Query().Get("bucket")
	measurement := request.URL.Query().Get("measurement")
	log.Printf("%v, %v\n", bucket, measurement)
	fieldsQuery := fmt.Sprintf(`import "influxdata/influxdb/schema" schema.fieldKeys(bucket: "%s")`, bucket)
	result, err := queryAPI.Query(context.Background(), fieldsQuery)
	if err != nil {
		log.Printf("error in labels query: %v", err)
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}
	if result == nil {
		log.Println("no buckets")
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}

	fields := make(map[string]bool)

	for result.Next() {
		//log.Printf("%v\n", result.Record())
		value := result.Record().ValueByKey("_value")
		fields[value.(string)] = true
	}
	outputStruct := &responseType2{
		make([]string, len(fields)),
	}
	index := 0
	for field := range fields {
		outputStruct.Fields[index] = field
		index++
	}

	output, err := json.Marshal(outputStruct)

	if err != nil {
		log.Printf("error in json encoding: %v", err)
		http.Error(response, "internal server error", http.StatusInternalServerError)
	}

	response.Write(output)

}
