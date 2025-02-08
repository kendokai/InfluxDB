package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Filter struct {
	measurement string `json:"measurement"`
	field       string `json:"field"`
	filterType  string `json:"type"`
	value       string `json:"value"`
}

type TimeRange struct {
	start string `json:"start"`
	stop  string `json:"stop"`
}

type QueryBuilderOutput struct {
	bucket       string    `json:"bucket"`
	measurements []string  `json:"measurements"`
	fields       []string  `json:"fields"`
	timeRange    TimeRange `json:"range"`
	filters      []Filter  `json:"filters"`
}

func generateQuery(queryBuilderOutput QueryBuilderOutput) string {

	// bucket selection
	query := fmt.Sprintf("from(bucket: %q)\n", queryBuilderOutput.bucket)

	// time range
	query += fmt.Sprintf("|> range(start: %s", queryBuilderOutput.timeRange.start)
	if queryBuilderOutput.timeRange.stop != "" {
		query += fmt.Sprintf(", stop: %s", queryBuilderOutput.timeRange.stop)
	}
	query += ")\n"

	// shaping
	query += "|> filter(fn: (r) => ("
	for idx, field := range queryBuilderOutput.fields {
		if idx != 0 {
			query += " or "
		}
		query += fmt.Sprintf("r._field == %q", field)
	}
	query += ") and ("
	for idx, measurement := range queryBuilderOutput.measurements {
		if idx != 0 {
			query += " or "
		}
		query += fmt.Sprintf("r._measurement == %q", measurement)
	}
	query += "))\n"

	// filtering
	for _, filter := range queryBuilderOutput.filters {
		query += fmt.Sprintf("|> filter(fn: (r) => r._field != %q or r._value %s %s)\n", filter.field, filter.filterType, filter.value)
	}

	return query
}

func runQueryWithAdminPrivilege(q string) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	org := os.Getenv("INFLUXDB_ORG")
	if org == "" {
		log.Fatalln("failed to get org")
	}
	token := os.Getenv("INFLUXDB_API_TOKEN")
	if token == "" {
		log.Fatalln("failed to get token")
	}
	url := os.Getenv("INFLUXDB_URL")
	if url == "" {
		log.Fatalln("failed to get url")
	}

	client := influxdb2.NewClient(url, token)
	defer client.Close()
	queryAPI := client.QueryAPI(org)

	result, err := queryAPI.Query(
		context.Background(),
		q,
	)

	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	count := 0

	for result.Next() {
		fmt.Printf("Query Result: %v\n", result.Record().Values())
		count++
	}

	log.Printf("Count: %v\n", count)

	// Check for any query error
	if result.Err() != nil {
		log.Fatalf("Query error: %v", result.Err())
	}
}

func main() {
	qbo := QueryBuilderOutput{"airSensor", []string{"airSensors"}, []string{"temperature"}, TimeRange{"-1000000", ""}, []Filter{{"airSensors", "humidity", "<", "100"}}}
	query := generateQuery(qbo)
	fmt.Print(query)
	runQueryWithAdminPrivilege(query)
}
