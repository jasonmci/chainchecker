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

// FetchProfileHook fetches the profile and related information using a GraphQL query
func FetchProfileHook(sessionToken string) string {
    // Read the GraphQL query from the file
    queryBytes, err := os.ReadFile("FetchProfileHook.graphql")
    if err != nil {
        log.Fatalf("Failed to read GraphQL file: %v", err)
    }
    query := string(queryBytes)

    // Create the JSON payload
    payload := map[string]interface{}{
        "operationName": "FetchProfileHook",
        "variables":    map[string]interface{}{},
        "query":        query,
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
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "*/*")
    req.Header.Set("X-Session-Token", sessionToken)
    // Additional headers from the curl command
    req.Header.Set("Accept-Language", "en-US,en;q=0.9")
    req.Header.Set("Origin", "https://feeds-manager.main.prod.cldev.sh")
    req.Header.Set("Referer", "https://feeds-manager.main.prod.cldev.sh/")
    req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
    req.Header.Set("Sec-Fetch-Dest", "empty")
    req.Header.Set("Sec-Fetch-Mode", "cors")
    req.Header.Set("Sec-Fetch-Site", "same-site")

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
	return string(responseBody)
}
