package main

import (
	//"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        panic("Error loading .env file")
    }

    // Get the username and password from environment variables
    username    := os.Getenv("FEEDS_MANAGER_USERNAME")
    password    := os.Getenv("FEEDS_MANAGER_PASSWORD")

    // Setup command line flags
    //listNetworks := flag.Bool("list", false, "List all networks")
    //shortName := flag.String("network", "", "Short network name to fetch details for")
    flag.Parse()

    token, err := LoginUser(username, password)
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    fmt.Println("Obtained token:", token)
    
    FetchSession(token)
    FetchProfileHook(token)
    //FetchCCIPView(token)
    //FetchChainDetails(token, "21")
    fmt.Println("=-=-=-=-=-=-=-")
    FetchLaneDetails(token, "7")
    // chainId := "21"
    // chainResponse, err := FetchChainDetails(token, chainId)
    // if err != nil {
    //     log.Fatalf("Failed to fetch chain details: %v", err)
    // }


    // //Parse JSON data into struct
    // var response Response
    // if err := json.Unmarshal(jsonData, &response); err != nil {
    //     log.Fatalf("Error parsing JSON data: %v", err)
    // }

    // //Convert jsonData which is a string to a byte slice
    // jsonDataBytes := []byte(jsonData)
    // fmt.Println()
    // //Parse JSON data into struct
    // var data Data
    // if err := json.Unmarshal(jsonDataBytes, &data); err != nil {
    //     log.Fatalf("Error parsing JSON data: %v", err)
    //     fmt.Println(data)
    // }

    // if *listNetworks {
    //     PrintNetworkMappings()

    // } else if *shortName != "" {
    //     if fullName, ok := networkMappings[*shortName]; ok {
    //         fmt.Printf("Fetching details for network: %s\n", fullName)
    //     fetchNetworkDetails(&response.Data, fullName)
    //     fmt.Printf("Chain Details: %+v\n", chainResponse)

    // } else {
    //         fmt.Println("Invalid network short name.")
    //     }
    // } else {
    //     fmt.Println("No command specified. Use -list to list all networks or -network to specify a network short name.")
    // }

}
