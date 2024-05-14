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

type CCIPViewResponse struct {
    Data struct {
        CCIP struct {
            Chains []struct {
                ID          string `json:"id"`
                DisplayName string `json:"displayName"`
                Network     struct {
                    ID          string `json:"id"`
                    Name        string `json:"name"`
                    IconName    string `json:"iconName"`
                    ExplorerURL string `json:"explorerURL"`
                    ChainID     string `json:"chainID"`
                    ChainType   string `json:"chainType"`
                    Typename    string `json:"__typename"`
                } `json:"network"`
                Contracts []struct {
                    ID                    string `json:"id"`
                    Name                  string `json:"name"`
                    Address               string `json:"address"`
                    Tag                   string `json:"tag"`
                    TransferOwnershipStatus string `json:"transferOwnershipStatus"`
                    Typename              string `json:"__typename"`
                } `json:"contracts"`
                Typename string `json:"__typename"`
            } `json:"chains"`
            Lanes []struct {
                ID          string `json:"id"`
                DisplayName string `json:"displayName"`
                Status      string `json:"status"`
                ChainA      struct {
                    ID          string `json:"id"`
                    DisplayName string `json:"displayName"`
                    Network     struct {
                        ID          string `json:"id"`
                        Name        string `json:"name"`
                        IconName    string `json:"iconName"`
                        ExplorerURL string `json:"explorerURL"`
                        ChainID     string `json:"chainID"`
                        ChainType   string `json:"chainType"`
                        Typename    string `json:"__typename"`
                    } `json:"network"`
                    Typename string `json:"__typename"`
                } `json:"chainA"`
                ChainB struct {
                    ID          string `json:"id"`
                    DisplayName string `json:"displayName"`
                    Network     struct {
                        ID          string `json:"id"`
                        Name        string `json:"name"`
                        IconName    string `json:"iconName"`
                        ExplorerURL string `json:"explorerURL"`
                        ChainID     string `json:"chainID"`
                        ChainType   string `json:"chainType"`
                        Typename    string `json:"__typename"`
                    } `json:"network"`
                    Typename string `json:"__typename"`
                } `json:"chainB"`
                Typename string `json:"__typename"`
            } `json:"lanes"`
        } `json:"ccip"`
    } `json:"data"`
}

// FetchCCIPView fetches CCIP related data using a GraphQL query
func FetchCCIPView(sessionToken string) (*CCIPViewResponse, error) {
    // Read the GraphQL query from the file
    queryBytes, err := os.ReadFile("queries/FetchCCIPView.graphql")
    if err != nil {
        log.Fatalf("Failed to read GraphQL file: %v", err)
    }
    query := string(queryBytes)

    // Create the JSON payload
    payload := map[string]interface{}{
        "operationName": "FetchCCIPView",
        "variables":    map[string]interface{}{},
        "query":        query,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        log.Fatalf("Error marshaling payload: %v", err)
    }

    // Create and set up the HTTP request
    req, err := http.NewRequest("POST", GqlEndpoint, bytes.NewReader(payloadBytes))
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
    var response CCIPViewResponse
    err = json.Unmarshal(responseBody, &response)
    if err != nil {
        log.Fatalf("Error parsing JSON response: %v", err)
    }
    return &response, nil
}

func PrintChains(chains []struct {
    ID          string `json:"id"`
    DisplayName string `json:"displayName"`
    Network     struct {
        ID          string `json:"id"`
        Name        string `json:"name"`
        IconName    string `json:"iconName"`
        ExplorerURL string `json:"explorerURL"`
        ChainID     string `json:"chainID"`
        ChainType   string `json:"chainType"`
        Typename    string `json:"__typename"`
    } `json:"network"`
    Contracts []struct {
        ID                    string `json:"id"`
        Name                  string `json:"name"`
        Address               string `json:"address"`
        Tag                   string `json:"tag"`
        TransferOwnershipStatus string `json:"transferOwnershipStatus"`
        Typename              string `json:"__typename"`
    } `json:"contracts"`
    Typename string `json:"__typename"`
}) {
    fmt.Println("Chains:")
    for _, chain := range chains {
        fmt.Printf("Chain ID: %s, Network Name: %s, Chain Type: %s\n", chain.ID, chain.Network.Name, chain.Network.ChainType)
    }
}

// PrintLanes prints all lanes from the response
func PrintLanes(lanes []struct {
    ID          string `json:"id"`
    DisplayName string `json:"displayName"`
    Status      string `json:"status"`
    ChainA      struct {
        ID          string `json:"id"`
        DisplayName string `json:"displayName"`
        Network     struct {
            ID          string `json:"id"`
            Name        string `json:"name"`
            IconName    string `json:"iconName"`
            ExplorerURL string `json:"explorerURL"`
            ChainID     string `json:"chainID"`
            ChainType   string `json:"chainType"`
            Typename    string `json:"__typename"`
        } `json:"network"`
        Typename string `json:"__typename"`
    } `json:"chainA"`
    ChainB struct {
        ID          string `json:"id"`
        DisplayName string `json:"displayName"`
        Network     struct {
            ID          string `json:"id"`
            Name        string `json:"name"`
            IconName    string `json:"iconName"`
            ExplorerURL string `json:"explorerURL"`
            ChainID     string `json:"chainID"`
            ChainType   string `json:"chainType"`
            Typename    string `json:"__typename"`
        } `json:"network"`
        Typename string `json:"__typename"`
    } `json:"chainB"`
    Typename string `json:"__typename"`
}) {
    fmt.Println("Lanes:")
    for _, lane := range lanes {
        fmt.Printf("Lane ID: %s, Status: %s, Chain A: %s (%s), Chain B: %s (%s)\n",
            lane.ID, lane.Status, lane.ChainA.Network.Name, lane.ChainA.Network.ChainType, lane.ChainB.Network.Name, lane.ChainB.Network.ChainType)
    }
}