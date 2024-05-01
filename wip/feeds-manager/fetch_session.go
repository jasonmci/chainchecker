package main

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// FetchSessionInfo fetches the session information using the GraphQL query
func FetchSession(sessionToken string) string {
    // Read the GraphQL query from the file
    queryBytes, err := os.ReadFile("queries/FetchSession.graphql")
    if err != nil {
        log.Fatalf("Failed to read GraphQL file: %v", err)
    }
    query := string(queryBytes)

    // Create the JSON payload
    payload := map[string]interface{}{
        "operationName": "FetchSession",
        "variables":    map[string]interface{}{},
        "query":        query,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        log.Fatalf("Error marshaling payload: %v", err)
    }

    // Create a new HTTP request
    req, err := http.NewRequest("POST", "https://gql.feeds-manager.main.prod.cldev.sh/query", bytes.NewReader(payloadBytes))
    if err != nil {
        log.Fatalf("Error creating request: %v", err)
    }

    setCommonHeaders(req, sessionToken)

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Error executing request: %v", err)
    }
    defer resp.Body.Close()

    // Read and print the response body
    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Error reading response body: %v", err)
    }

	return string(responseBody)
}