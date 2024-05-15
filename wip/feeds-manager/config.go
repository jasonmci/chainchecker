package main

import (
	"os"

	"github.com/joho/godotenv"
)

// GqlEndpoint is the endpoint for GraphQL, shared across multiple files.
var GqlEndpoint string
var Username string
var Password string

func init() {
	if err := godotenv.Load(); err != nil {
        panic("Error loading .env file")
    }
    GqlEndpoint = os.Getenv("CLO_ENDPOINT")
	Username = os.Getenv("CLO_USERNAME")
	Password = os.Getenv("CLO_PASSWORD")
}
