package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        panic("Error loading .env file")
    }

    // Get the username and password from environment variables
    username := os.Getenv("FEEDS_MANAGER_USERNAME")
    password := os.Getenv("FEEDS_MANAGER_PASSWORD")

    // Generate the payload with the username and password
    payload := GenerateAuthPayload(username, password)

    // Marshal the payload to JSON
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        panic(err)
    }
    body := bytes.NewReader(payloadBytes)

    // Create and send the HTTP request
    req, err := http.NewRequest("POST", "xyz", body)
    if err != nil {
        panic(err)
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // Read and print the response body
    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(responseBody))
}
