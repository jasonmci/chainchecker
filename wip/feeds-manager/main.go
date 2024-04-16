package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Structs based on the JSON structure
type Network struct {
    Name string `json:"name"`
}

type Chain struct {
    Network Network `json:"network"`
}

type CCIP struct {
    Chains []Chain `json:"chains"`
}

type Data struct {
    CCIP CCIP `json:"ccip"`
}

type Response struct {
    Data Data `json:"data"`
}

// Network mappings from short name to full name
var networkMappings = map[string]string{
    "KROMA": "Kroma Mainnet",
    "WEMIX": "WEMIX Mainnet",
    "GNO":   "GnosisChain Mainnet",
    "POLYX": "Polygon zkEVM Mainnet",
    "OPT":   "Optimism Mainnet",
    "AVAX":  "Avalanche Mainnet",
    "POLY":  "Polygon Mainnet",
    "BSC":   "BSC Mainnet",
    "ETH":   "Ethereum Mainnet",
    "ARB":   "Arbitrum Mainnet",
    "BASE":  "Base Mainnet",
}

func fetchNetworkDetails(data *Data, networkName string) {
    fmt.Printf("Details for %s:\n", networkName)
    for _, chain := range data.CCIP.Chains {
        if chain.Network.Name == networkName {
            fmt.Printf(" - Network Name: %s\n", chain.Network.Name)
            return
        }
    }
    fmt.Println("Network not found.")
}

func main() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        panic("Error loading .env file")
    }

    // Get the username and password from environment variables
    username    := os.Getenv("FEEDS_MANAGER_USERNAME")
    password    := os.Getenv("FEEDS_MANAGER_PASSWORD")

    // Setup command line flags
    listNetworks := flag.Bool("list", false, "List all networks")
    shortName := flag.String("network", "", "Short network name to fetch details for")
    flag.Parse()
    flag.Parse()

    token, err := LoginUser(username, password)
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    fmt.Println("Obtained token:", token)
    
    FetchSession(token)
    FetchProfileHook(token)
    jsonData := FetchCCIPView(token)

    // Parse JSON data into struct
    var response Response
    if err := json.Unmarshal(jsonData, &response); err != nil {
        log.Fatalf("Error parsing JSON data: %v", err)
    }

    // Convert jsonData which is a string to a byte slice
    jsonDataBytes := []byte(jsonData)
    fmt.Println()
    // Parse JSON data into struct
    var data Data
    if err := json.Unmarshal(jsonDataBytes, &data); err != nil {
        log.Fatalf("Error parsing JSON data: %v", err)
        fmt.Println(data)
    }

    if *listNetworks {
        for _, chain := range response.Data.CCIP.Chains {
            fmt.Println(chain.Network.Name)
        }
    } else if *shortName != "" {
        if fullName, ok := networkMappings[*shortName]; ok {
            fetchNetworkDetails(&response.Data, fullName)
        } else {
            fmt.Println("Invalid network short name.")
        }
    } else {
        fmt.Println("No command specified. Use -list to list all networks or -network to specify a network short name.")
    }

}

// ListNetworkNames prints the names of all networks
func ListNetworkNames(data *Data) {
    fmt.Println("Listing all networks:")
    for _, chain := range data.CCIP.Chains {
        fmt.Printf(" - %s\n", chain.Network.Name)
    }
}