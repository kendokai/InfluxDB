package main

import (
	"context"
	"fmt"
	"log"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

/*
Script to setup the influxdb database
*/
func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bucket := "airSensor"
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

	query := fmt.Sprintf(`
		   from(bucket: %q) |> range(start: -10000000)
		   `, bucket)

	queryAPI := client.QueryAPI(org)

	result, err := queryAPI.Query(
		context.Background(),
		query,
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
