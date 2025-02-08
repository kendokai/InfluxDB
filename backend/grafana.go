package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func getUIDHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("request for UID")
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Error loading .env file")
	}
	dashboardUID := os.Getenv("GRAFANA_DASHBOARD_UID")
	if dashboardUID == "" {
		log.Println("cant find UID")
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}
	type DashboardUID struct {
		Uid string `json:"uid"`
	}
	uid := DashboardUID{Uid: dashboardUID}
	output, err := json.Marshal(&uid)
	if err != nil {
		log.Println("error marshalling json")

		http.Error(response, "internal server error", http.StatusInternalServerError)
	}
	_, err = response.Write(output)

	if err != nil {
		log.Println("error writing response")

		http.Error(response, "internal server error", http.StatusInternalServerError)
	}
}

type UpdateDashboardRequestData struct {
	Dashboard Dashboard `json:"dashboard"`
	FolderId  int       `json:"folderId"`
	Overwrite bool      `json:"overwrite"`
}

type GetDashboardMeta struct {
}

type GetDashboardRequestData struct {
	Meta      GetDashboardMeta `json:"meta"`
	Dashboard Dashboard        `json:"dashboard"`
}

type Dashboard struct {
	Id            interface{} `json:"id"`
	Uid           string      `json:"uid"`
	Title         string      `json:"title"`
	Tags          []string    `json:"tags"`
	SchemaVersion int         `json:"schemaVersion"`
	Version       int         `json:"version"`
	Panels        []Panel     `json:"panels"`
}

type Panel struct {
	Id    int    `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`

	Targets    []Target   `json:"targets"`
	Datasource Datasource `json:"datasource"`
}

type Target struct {
	Query      string     `json:"query"`
	Datasource Datasource `json:"datasource"`
}

type Datasource struct {
	Type string `json:"type"`
	Uid  string `json:"uid"`
}

func getDashboard(dashboardUid string) (*Dashboard, error) {
	GrafanaURL := os.Getenv("GRAFANA_URL")
	GrafanaAPIToken := os.Getenv("GRAFANA_API_TOKEN")

	request, err := http.NewRequest("GET", GrafanaURL+"/api/dashboards/uid/"+dashboardUid, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+GrafanaAPIToken)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch dashboard: %s", response.Status)
	}

	var dashboardResponse GetDashboardRequestData
	//bodyBytes, err := io.ReadAll(response.Body)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Print(string(bodyBytes))
	err = json.NewDecoder(response.Body).Decode(&dashboardResponse)
	if err != nil {
		return nil, err
	}

	return &dashboardResponse.Dashboard, nil
}

func updateDashboard(dashboard *Dashboard) error {
	GrafanaURL := os.Getenv("GRAFANA_URL")
	GrafanaAPIToken := os.Getenv("GRAFANA_API_TOKEN")
	write := UpdateDashboardRequestData{
		Dashboard: *dashboard,
		FolderId:  0,
		Overwrite: true,
	}

	payload, err := json.Marshal(write)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", GrafanaURL+"/api/dashboards/db", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+GrafanaAPIToken)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated {
		fmt.Println("Dashboard updated successfully")
		return nil
	}

	return fmt.Errorf("failed to update dashboard: %s", response.Status)
}

func getAndUpdateDashboard(newQuery string) error {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dashboardUID := os.Getenv("GRAFANA_DASHBOARD_UID")
	dashboard, err := getDashboard(dashboardUID)
	if err != nil {
		return err
	}
	newDashboard := Dashboard{
		Id:    dashboard.Id,
		Uid:   dashboard.Uid,
		Title: dashboard.Title,
		Panels: []Panel{{
			Id:    dashboard.Panels[0].Id,
			Type:  dashboard.Panels[0].Type,
			Title: dashboard.Panels[0].Title,
			Targets: []Target{{
				Query:      newQuery,
				Datasource: dashboard.Panels[0].Datasource,
			}},
			Datasource: dashboard.Panels[0].Datasource,
		}},
		SchemaVersion: dashboard.SchemaVersion,
		Version:       dashboard.Version + 1,
	}

	err = updateDashboard(&newDashboard)
	if err != nil {
		return err
	}
	log.Printf("Successfully updated dashboard %q!!!\n", dashboardUID)

	return nil
}
