package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type CCIPLaneResponse struct {
	Data struct {
		CCIP struct {
			Lane struct {
				ID              string `json:"id"`
				DisplayName     string `json:"displayName"`
				Status          string `json:"status"`
				DeployedTemplate map[string]struct {
					CommitStoreAddress string `json:"commitStoreAddress"`
					OffRampAddress     string `json:"offRampAddress"`
					OnRampAddress      string `json:"onRampAddress"`
				} `json:"deployedTemplate"`
			} `json:"lane"`
		} `json:"ccip"`
	} `json:"data"`
}

func FetchLaneDetails(sessionToken, laneID string ) (*CCIPLaneResponse, error) {

	queryBytes, err := os.ReadFile("queries/FetchCCIPLaneDetailsView.graphql")
    if err != nil {
        log.Fatalf("Failed to read GraphQL file: %v", err)
    }
    query := string(queryBytes)

    payload := map[string]interface{}{
        "operationName": "FetchCCIPLaneDetailsView",
        "variables": map[string]interface{}{
            "id": laneID,
        },
        "query": query,
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        log.Fatalf("Error marshaling payload: %v", err)
    }

    // Create and set up the HTTP request
    req, err := http.NewRequest("POST", "https://gql.feeds-manager.main.prod.cldev.sh/query", bytes.NewReader(payloadBytes))
    if err != nil {
        log.Fatalf("Error creating request: %v", err)
    }
    setCommonHeaders(req, sessionToken)

    // Execute the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Error executing request: %v", err)
    }
    defer resp.Body.Close()

    // Handle the response
    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Error reading response body: %v", err)
    }

    // Parse the JSON response
    var response CCIPLaneResponse
    err = json.Unmarshal(responseBody, &response)
    if err != nil {
        return nil, fmt.Errorf("error parsing JSON response: %w", err)
    }

    return &response, nil
}
