package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/joho/godotenv"
)

func setupInflux(client influxdb2.Client, url string, envfile string) (string, string) {
	username := "admin"
	password := "password"
	org := "admin"
	bucket := "default"

	response, err := client.Setup(context.Background(), username, password, org, bucket, 0)
	if err != nil {
		log.Fatalf("Error in setup: %v\n", err)
	}
	log.Println("Succussfully Setup Server")
	auth := response.Auth
	file, err := os.Create(envfile)
	if err != nil {
		log.Fatalf("Error creating '.env': %v\n", err)
	}
	defer file.Close()
	file.WriteString(
		fmt.Sprintf(
			"INFLUXDB_API_TOKEN=%q\nINFLUXDB_URL=%q\nINFLUXDB_ORG=%q\n",
			*auth.Token,
			url,
			response.Org.Name,
		),
	)

	return *auth.Token, response.Org.Name
}

func loadSampleData(client influxdb2.Client, taskName string, sampleDataSet string, bucketName string, orgID string) {

	bucketsAPI := client.BucketsAPI()
	retentionType := domain.RetentionRuleTypeExpire
	bucketThing := &domain.Bucket{
		Name: bucketName,
		RetentionRules: []domain.RetentionRule{
			{
				Type:         &retentionType,
				EverySeconds: int64((24 * time.Hour).Seconds()), // Set retention policy to 24 hours
			},
		},
		OrgID: &orgID,
	}
	newBucketThing, err := bucketsAPI.CreateBucket(context.Background(), bucketThing)
	if err != nil {
		log.Fatalf("Error creating bucket: %v\n", err)
	}
	log.Printf("created bucket %q at %v\n", newBucketThing.Name, newBucketThing.CreatedAt)
	query := fmt.Sprintf(`
		   import "influxdata/influxdb/sample"
			option task = {
			name: %q,
			every: 15m,
			}
		   sample.data(set: %q)
		   	|> to(bucket: %q)
		   `,
		taskName, sampleDataSet, newBucketThing.Name)

	queryAPI := client.QueryAPI(orgID)

	result, err := queryAPI.Query(
		context.Background(),
		query,
	)

	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	if result.Err() != nil {
		log.Fatalf("Query error: %v", result.Err())
	}

	count := 0

	for result.Next() {
		// fmt.Printf("Query Result: %v\n", result.Record().Values())
		count++
	}

	log.Printf("Count: %v\n", count)

	// Check for any query error

}

/*
Script to setup the influxdb database
*/
func main() {
	envfile := "../../.env"

	err := godotenv.Load(envfile)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("INFLUXDB_API_TOKEN")
	orgName := os.Getenv("INFLUXDB_ORG")
	url := "http://influxdb:8086"

	if token == "" {
		client := influxdb2.NewClient(url, token)
		defer client.Close()
		token, orgName = setupInflux(client, url, envfile)
		log.Printf("API token: %v\n", token)
	}
	client := influxdb2.NewClient(url, token)
	defer client.Close()
	orgAPI := client.OrganizationsAPI()

	// List organizations
	orgs, err := orgAPI.GetOrganizations(context.Background())
	if err != nil {
		log.Fatalf("Error retrieving organizations: %v", err)
	}
	orgID := ""
	for _, Org := range *orgs {
		if Org.Name == orgName {
			orgID = *Org.Id
		}
	}
	if orgID == "" {
		log.Fatalf("couldn't find org %q\n", orgName)
	}
	loadSampleData(client, "Collect Bitcoin sample data", "bitcoin", "bitcoin", orgID)
	loadSampleData(client, "Collect air sensor sample data", "airSensor", "airSensor", orgID)
	loadSampleData(client, "Collect NOAA NDBC sample data", "noaa", "NOAA NDBC", orgID)
	//loadSampleData(client, "Collect Bitcoin sample data", "bitcoin", "default", orgID)
	//loadSampleData(client, "Collect air sensor sample data", "airSensor", "default", orgID)
	//loadSampleData(client, "Collect NOAA NDBC sample data", "noaa", "default", orgID)

}
