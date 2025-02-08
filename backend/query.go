package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type Filter struct {
	Measurement string `json:"measurement"`
	Field       string `json:"field"`
	FilterType  string `json:"operator"`
	Value       string `json:"value"`
}

type TimeRange struct {
	Start string `json:"start"`
	Stop  string `json:"stop"`
}

type QueryBuilderOutput struct {
	Bucket       string    `json:"bucket"`
	Measurements []string  `json:"measurements"`
	Fields       []string  `json:"fields"`
	TimeRange    TimeRange `json:"timeRange"`
	Filters      []Filter  `json:"filters"`
}

func runQuery(client influxdb2.Client, query string) (*api.QueryTableResult, error) {
	org := os.Getenv("INFLUXDB_ORG")
	if org == "" {
		log.Println("failed to get org")
		return nil, errors.New("failed to get org")
	}
	queryAPI := client.QueryAPI(org)
	return queryAPI.Query(context.Background(), query)
}

func runQueryHandler(response http.ResponseWriter, request *http.Request) {
	var qbo QueryBuilderOutput
	var err error
	/*client, err := getClient(response, request)

	if err != nil {
		log.Printf("error in getting client: %v", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}*/

	err = json.NewDecoder(request.Body).Decode(&qbo)

	if err != nil {
		log.Printf("error in query json unmarshalling: %v", err)
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	query, err := generateQuery(qbo)
	if err != nil {
		log.Printf("error in query generation: %v", err)
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("generated query: %q", query)

	log.Printf("generated query %q", query)

	err = getAndUpdateDashboard(query)
	if err != nil {
		log.Printf("error in updating dashboard: %v", err)
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	/*
			data, err := runQuery(client, query)

			if err != nil {
				log.Printf("error in query execution: %v", err)
				http.Error(response, "Invalid request", http.StatusBadRequest)
			return
			}
			output := "["
		firstItem := true
		for data.Next() {
			record, err := json.Marshal(data.Record().Values())
				if err != nil {
					log.Printf("error in json marshalling: %v", err)
					http.Error(response, "Internal Server Error", http.StatusInternalServerError)
				}
			if !firstItem {
				output += ","
			} else {
				firstItem = false
			}
			output += string(record)
		}

		output += "]"
		log.Printf("marshalled json: %s", output)

		if err != nil {

		}

			_, _ = response.Write([]byte(output))

	*/
}

func queryGenerationHandler(response http.ResponseWriter, request *http.Request) {
	//client, _ := getClient(response, request)
	var qbo QueryBuilderOutput
	fmt.Print(request.Body)
	err := json.NewDecoder(request.Body).Decode(&qbo)

	if err != nil {
		log.Printf("error in query json unmarshalling: %v", err)
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	query, err := generateQuery(qbo)
	if err != nil {
		log.Printf("error in query generation: %v", err)
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err = response.Write([]byte(query))

	if err != nil {
		log.Printf("error in response writer: %v", err)
	}

}

func generateQuery(queryBuilderOutput QueryBuilderOutput) (string, error) {
	// bucket selection
	query := fmt.Sprintf("from(bucket: %q)\n", queryBuilderOutput.Bucket)
	const encoding_format = "2006-01-02T15:04-0700"
	// time range
	startTimeStamp, err := time.Parse(encoding_format, queryBuilderOutput.TimeRange.Start)

	if err != nil {
		return "", err
	}

	query += fmt.Sprintf("|> range(start: %d", startTimeStamp.Unix())
	if queryBuilderOutput.TimeRange.Stop != "" {
		stopTimeStamp, err := time.Parse(encoding_format, queryBuilderOutput.TimeRange.Stop)

		if err != nil {
			return "", err
		}
		query += fmt.Sprintf(", stop: %d", stopTimeStamp.Unix())
	}
	query += ")\n"

	// shaping
	query += "|> filter(fn: (r) => ("
	for idx, field := range queryBuilderOutput.Fields {
		if idx != 0 {
			query += " or "
		}
		query += fmt.Sprintf("r._field == %q", field)
	}
	query += ") and ("
	for idx, measurement := range queryBuilderOutput.Measurements {
		if idx != 0 {
			query += " or "
		}
		query += fmt.Sprintf("r._measurement == %q", measurement)
	}
	query += "))\n"

	// filtering
	for _, filter := range queryBuilderOutput.Filters {
		operator := ""
		switch filter.FilterType {
		case "==":
			operator = "=="
		case "!=":
			operator = "!="
		case ">":
			operator = ">"
		case "<":
			operator = "<"
		case ">=":
			operator = ">="
		case "<=":
			operator = "<="
		default:
			return "", errors.New("filter type is invalid")

		}
		query += fmt.Sprintf("|> filter(fn: (r) => r._field != %q or r._value %s %s)\n", filter.Field, operator, filter.Value)
	}
	return query, nil

}
