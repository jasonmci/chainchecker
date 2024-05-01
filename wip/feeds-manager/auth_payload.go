package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)


// LoginUser performs the login and returns the token
func LoginUser(username, password string) (string, error) {
    // Construct the payload
    query, err := os.ReadFile("queries/LoginMutation.graphql")
    if err != nil {
        log.Fatalf("Failed to read query file: %v", err)
    }
    payload := map[string]interface{}{
        "operationName": "Login",
        "variables": map[string]interface{}{
            "input": map[string]interface{}{
                "email":    username,
                "password": password,
            },
        },
        "query": string(query),
    }
    
    gqlEndpoint := os.Getenv("FEEDS_MANAGER_ENDPOINT")
    payloadBytes, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", gqlEndpoint, bytes.NewReader(payloadBytes))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Read the response body
    responseBody, _ := io.ReadAll(resp.Body)

    // Here you should add proper JSON unmarshaling based on your GraphQL server's response structure
    var result map[string]interface{}
    json.Unmarshal(responseBody, &result)

    // Extract the token from the response. This assumes the token is returned directly under session->token.
    // Adjust the path according to your actual response structure
    token := result["data"].(map[string]interface{})["login"].(map[string]interface{})["session"].(map[string]interface{})["token"].(string)

    return token, nil
}

