package main

import (
	//"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Define a struct to hold lane information
type ChainInfo struct {
    Network string
    PaymentTokens []string
    TransferTokens []string
}

type LaneInfo struct {
    Network string
    PaymentTokens []string
    TransferTokens []string
}


func parseChainInfo(arg string) (ChainInfo, error) {
    parts := strings.Split(arg, ",")
    if len(parts) < 3 {
        return ChainInfo{}, fmt.Errorf("incorrect format for lane argument. expecting 'network,paymentToken,transferToken'")
    }

    networkID, found := getNetworkID(parts[0])
    if !found {
        return ChainInfo{}, fmt.Errorf("network name %s not recognized", parts[0])
    }

    // laneID, found := getLaneID(parts[0], parts[1])
    // if !found {
    //     return ChainInfo{}, fmt.Errorf("no lane found for networks %s and %s", parts[0], parts[1])
    // }

    return ChainInfo{
        Network: networkID, // Store the ID, not the name
        PaymentTokens: parts[1:2],
        TransferTokens: parts[2:],
    }, nil
}

// func parseLaneInfo(arg string, laneB string) (LaneInfo, error) {
//     parts := strings.Split(arg, ",")
//     if len(parts) < 3 {
//         return LaneInfo{}, fmt.Errorf("incorrect format for lane argument. expecting 'network1,network2,paymentToken,transferToken'")
//     }

//     laneID, found := getLaneID(parts[0], parts[1])
//     if !found {
//         return LaneInfo{}, fmt.Errorf("no lane found for networks %s and %s", parts[0], parts[1])
//     }

//     return LaneInfo{
//         Network: laneID,
//         PaymentTokens: parts[2:3],
//         TransferTokens: parts[3:],
//     }, nil
// }


func fetchDataForLane(chain ChainInfo) {
    fmt.Printf("Network: %s\n", chain.Network)
    fmt.Printf("Payment Tokens: %v\n", chain.PaymentTokens)
    fmt.Printf("Transfer Tokens: %v\n", chain.TransferTokens)

    // Network configs
    // fetch arm address
    fmt.Printf("Fetching ARM address for network: %s\n", chain.Network)

    // fetch router address
    fmt.Printf("Fetching Router address for network: %s\n", chain.Network)

    // fetch Price Registry
    fmt.Printf("Fetching Price Registry for network: %s\n", chain.Network)


    // now we are going to fetch onramp offramp and commit stores
    // fetch onramp
    fmt.Printf("Fetching Onramp for network: %s\n", chain.Network)

    // fetch offramp
    fmt.Printf("Fetching Offramp for network: %s\n", chain.Network)

    // fetch commit store
    fmt.Printf("Fetching Commit Store for network: %s\n", chain.Network)

    
    // and now we need to fetch the fee tokens for fees and tokens to transfer

    // fetch fee tokens
    fmt.Printf("Fetching Fee Tokens for network: %s\n", chain.Network)

    // fetch tokens to transfer
    fmt.Printf("Fetching Tokens to Transfer for network: %s\n", chain.Network)

    // fetch token pool
    fmt.Printf("Fetching Token Pool for network: %s\n", chain.Network)

        // Implement actual fetching logic here
}


func main() {

    var a, b string
    flag.StringVar(&a, "A", "", "Details for Lane A (network,paymentToken,transferToken)")
    flag.StringVar(&b, "B", "", "Details for Lane B (network,paymentToken,transferToken)")
    flag.Parse()

    laneA, err := parseChainInfo(a)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    laneB, err := parseChainInfo(b)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Assuming you have functions to handle the API calls
    if laneA.Network != "" {
        fmt.Println("Fetching data for Lane A...")
        fetchDataForLane(laneA)
    }

    if laneB.Network != "" {
        fmt.Println("Fetching data for Lane B...")
        fetchDataForLane(laneB)
    }

    partsA := strings.Split(a, ",")
    partsB := strings.Split(b, ",")
    laneID, _ := getLaneID(partsA[0], partsB[0])

    fmt.Printf("chain ID: %s, other chain ID: %s\n", partsA[0], partsB[0])
    fmt.Printf("Lane ID: %s\n", laneID)

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
    
    //FetchSession(token)
    //FetchProfileHook(token)
    //FetchChainDetails(token, "21")
    // viewResponseBody := FetchCCIPView(token)



    //chainAResponseBody := FetchChainDetails(token, laneID)
    //chainBResponseBody := FetchChainDetails(token, "21")
    laneResponseBody := FetchLaneDetails(token, laneID)

    //fmt.Println("Chain A Response:", chainAResponseBody)
    //fmt.Println("Chain B Response:", chainBResponseBody)
    fmt.Println("Lane Response:", laneResponseBody)

    // var viewResponse CCIPViewResponse
    // if err := json.Unmarshal(viewResponseBody, &viewResponse); err != nil {
    //     log.Fatalf("Error parsing JSON: %v", err)
    // }
    
    // // Example of using the data
    // for _, chain := range viewResponse.Data.CCIP.Chains {
    //     fmt.Printf("Chain ID: %s, Display Name: %s, Network: %s\n",
    //         chain.ID, chain.DisplayName, chain.Network.Name)
    // }

    // var chainResponse CCIPChainResponse
    // if err := json.Unmarshal(chainResponseBody, &chainResponse); err != nil {
    //     log.Fatalf("Error parsing JSON: %v", err)
    // }
    
    // // Example of using the data
    // for _, token := range chainResponse.Data.CCIP.Chain.SupportedTokens {
    //     fmt.Printf("Chain ID: %s, Display Name: %s, Network: %s\n",
    //         token.Address, token.Token, token.PriceType)
    // }

    // var laneResponse CCIPLaneResponse
    // if err := json.Unmarshal(laneResponseBody, &laneResponse); err != nil {
    //     log.Fatalf("Error parsing JSON: %v", err)
    // }
    
    // // Example of using the data
    // for _, chain := range viewResponse.Data.CCIP.Chains {
    //     fmt.Printf("Chain ID: %s, Display Name: %s, Network: %s\n",
    //         chain.ID, chain.DisplayName, chain.Network.Name)
    // }

}
