package main

import (
	//"bytes"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"os"
)

type CCIPChainResponse struct {
    Data struct {
        CCIP struct {
            Chain struct {
                ID               string `json:"id"`
                DisplayName      string `json:"displayName"`
                Network               struct {
                    ID          string `json:"id"`
                    Name        string `json:"name"`
                    IconName    string `json:"iconName"`
                    ExplorerURL string `json:"explorerURL"`
                    ChainID     string `json:"chainID"`
                    ChainType   string `json:"chainType"`
                    Typename    string `json:"__typename"`
                } `json:"network"`
                DeployedTemplate DeployedTemplate `json:"deployedTemplate"`
                SupportedTokens  []struct {
                    Token         string `json:"token"`
                    Address       string `json:"address"`
                    PriceType     string `json:"priceType"`
                    TokenPoolType string `json:"tokenPoolType"`
                    Typename      string `json:"__typename"`
                } `json:"supportedTokens"`
                Contracts []struct {
                    ID                    string `json:"id"`
                    Address               string `json:"address"`
                    Tag                   string `json:"tag"`
                    TransferOwnershipStatus string `json:"transferOwnershipStatus"`
                    Name                  string `json:"name"`
                    Semver                string `json:"semver"`
                    Metadata              Metadata `json:"metadata"`
                    OwnerType             string `json:"ownerType"`
                    OwnerAddress          string `json:"ownerAddress"`
                    PendingOwnerAddress   string `json:"pendingOwnerAddress"`
                    PendingOwnerType      string `json:"pendingOwnerType"`

                    Typename string `json:"__typename"`
                } `json:"contracts"`
                WorkflowRuns []struct {
                    ID           string `json:"id"`
                    Name         string `json:"name"`
                    Status       string `json:"status"`
                    WorkflowType string `json:"workflowType"`
                    CreatedAt    string `json:"createdAt"`
                    Actions      []struct {
                        ActionType string `json:"actionType"`
                        Name       string `json:"name"`
                        Run        struct {
                            ID       string `json:"id"`
                            Status   string `json:"status"`
                            Network  struct {
                                ID          string `json:"id"`
                                Name        string `json:"name"`
                                IconName    string `json:"iconName"`
                                ExplorerURL string `json:"explorerURL"`
                                Typename    string `json:"__typename"`
                            } `json:"network"`
                            CreatedAt string `json:"createdAt"`
                            Typename  string `json:"__typename"`
                        } `json:"run"`
                        Tasks []struct {
                            Name string `json:"name"`
                            Run  struct {
                                Error   string `json:"error"`
                                ID      string `json:"id"`
                                Input   string `json:"input"`
                                Output  string `json:"output"`
                                Status  string `json:"status"`
                                TxHash  string `json:"txHash"`
                                Typename string `json:"__typename"`
                            } `json:"run"`
                            Typename string `json:"__typename"`
                        } `json:"tasks"`
                        Typename string `json:"__typename"`
                    } `json:"actions"`
                    Typename string `json:"__typename"`
                } `json:"workflowRuns"`
                Typename string `json:"__typename"`
            } `json:"chain"`
            Typename string `json:"__typename"`
        } `json:"ccip"`
    } `json:"data"`
}

type DeployedTemplate struct {
    ArmProxies map[string]ArmProxy `json:"armProxies"`
    Routers map[string]struct {
        Status         string `json:"status"`
        TypeAndVersion string `json:"typeAndVersion"`
    } `json:"routers"`
    PriceRegistries map[string]struct {
        FeeTokens         []string `json:"feeTokens"`
        PriceUpdaters     []string `json:"priceUpdaters"`
        StalenessThreshold int     `json:"stalenessThreshold"`
        Status             string  `json:"status"`
        TypeAndVersion     string  `json:"typeAndVersion"`
    } `json:"priceRegistries"` 
    Arms map[string]struct {
        IsCursed      bool   `json:"isCursed"`
        Status        string `json:"status"`
        TypeAndVersion string `json:"typeAndVersion"`
        Config struct {
            BlessWeightThreshold int `json:"blessWeightThreshold"`
            CurseWeightThreshold int `json:"curseWeightThreshold"`
        } `json:"config"`
    } `json:"arms"`
    Tokens map[string]struct {
        PoolAddress    string `json:"poolAddress"`
        TokenAddress   string `json:"tokenAddress"`
        TypeAndVersion string `json:"typeAndVersion"`
    } `json:"tokens"`
}

type ArmProxy struct {
    Arm           string `json:"arm"`
    Status        string `json:"status"`
    TypeAndVersion string `json:"typeAndVersion"`
}

type Arm struct {
    Config        Config `json:"config"`
    IsCursed      bool   `json:"isCursed"`
    Status        string `json:"status"`
    TypeAndVersion string `json:"typeAndVersion"`
}

type Config struct {
    //BlessWeightThreshold int       `json:"blessWeightThreshold"`
    CurseWeightThreshold int       `json:"curseWeightThreshold"`
    Voters               []Voter   `json:"voters"`
}

type Voter struct {
    BlessVoteAddr    string `json:"blessVoteAddr"`
    BlessWeight      int    `json:"blessWeight"`
    CurseUnvoteAddr  string `json:"curseUnvoteAddr"`
    CurseVoteAddr    string `json:"curseVoteAddr"`
    CurseWeight      int    `json:"curseWeight"`
}

type Metadata struct {
    AllowList        []string                 `json:"allowList"`
    AllowListEnabled bool                     `json:"allowListEnabled"`
    OffRamps         map[string]Ramp          `json:"offRamps"`
    OnRamps          map[string]Ramp          `json:"onRamps"`
    Price            string                   `json:"price"`
    TokenAddress     string                   `json:"tokenAddress"`
    TokenSymbol      string                   `json:"tokenSymbol"`
    Type             string                   `json:"type"`
    TypeAndVersion   string                   `json:"typeAndVersion"`
}

type Ramp struct {
    Allowed          bool                     `json:"allowed"`
    RateLimiterConfig RateLimiterConfig       `json:"rateLimiterConfig"`
}

type RateLimiterConfig struct {
    Capacity    string `json:"capacity"`
    IsEnabled   bool   `json:"isEnabled"`
    Rate        string `json:"rate"`
}

func FetchChainDetails(sessionToken, chainID string) (*CCIPChainResponse, error) {
    queryBytes, err := os.ReadFile("FetchCCIPChainDetailsView.graphql")
    if err != nil {
        log.Fatalf("Failed to read GraphQL file: %v", err)
    }
    query := string(queryBytes)

    payload := map[string]interface{}{
        "operationName": "FetchCCIPChainDetailsView",
        "variables": map[string]interface{}{
            "id": chainID,
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
    var response CCIPChainResponse
    err = json.Unmarshal(responseBody, &response)
    if err != nil {
        return nil, fmt.Errorf("error parsing JSON response: %w", err)
    }

    // Print the chain contracts details
    for _, contract := range response.Data.CCIP.Chain.Contracts {
        fmt.Printf("Contract ID: %s, Name: %s, Token Symbol: %s, Address: %s, Tag %s, Semver %s\n", contract.ID, contract.Name, contract.Metadata.TokenSymbol, contract.Address, contract.Tag, contract.Semver)
    }
    fmt.Println("--------------------------------------------------")

    fmt.Println("ARMs:")
    for address, arm := range response.Data.CCIP.Chain.DeployedTemplate.Arms {
        fmt.Printf("ARM Address: %s\n", address)
        fmt.Printf("Is Cursed: %t\n", arm.IsCursed)
        fmt.Printf("Status: %s\n", arm.Status)
        fmt.Printf("Type and Version: %s\n", arm.TypeAndVersion)
        fmt.Printf("Bless Weight Threshold: %d\n", arm.Config.BlessWeightThreshold)
        fmt.Printf("Curse Weight Threshold: %d\n\n", arm.Config.CurseWeightThreshold)
    }

    // print some text as a separator
    fmt.Println("--------------------------------------------------")

    // print out arms details
    for _, armProxy := range response.Data.CCIP.Chain.DeployedTemplate.ArmProxies {
        fmt.Println("ARM Address:", armProxy.Arm)
        fmt.Println("Status:", armProxy.Status)
        fmt.Println("Type and Version:", armProxy.TypeAndVersion)
    }

    fmt.Println("--------------------------------------------------")

    // Iterate over the routers map
    for address, router := range response.Data.CCIP.Chain.DeployedTemplate.Routers {
        fmt.Printf("Router Address: %s\n", address)
        fmt.Printf("Status: %s\n", router.Status)
        fmt.Printf("Type and Version: %s\n\n", router.TypeAndVersion)
    }

    fmt.Println("--------------------------------------------------")

    // Iterate over the price registries map
    fmt.Println("Price Registries:")
    for address, registry := range response.Data.CCIP.Chain.DeployedTemplate.PriceRegistries {
        fmt.Printf("Registry Address: %s\n", address)
        fmt.Printf("Fee Tokens: %v\n", registry.FeeTokens)
        fmt.Printf("Price Updaters: %v\n", registry.PriceUpdaters)
        fmt.Printf("Staleness Threshold: %d\n", registry.StalenessThreshold)
        fmt.Printf("Status: %s\n", registry.Status)
        fmt.Printf("Type and Version: %s\n\n", registry.TypeAndVersion)
    }

    fmt.Println("--------------------------------------------------")
    // Iterate over the tokens map
    fmt.Println("Tokens:")
    for tokenName, token := range response.Data.CCIP.Chain.DeployedTemplate.Tokens {
        fmt.Printf("Token Name: %s\n", tokenName)
        fmt.Printf("Pool Address: %s\n", token.PoolAddress)
        fmt.Printf("Token Address: %s\n", token.TokenAddress)
        fmt.Printf("Type and Version: %s\n\n", token.TypeAndVersion)
    }

    // print the chain network details
    fmt.Printf("Chain ID: %s, Display Name: %s, Network: %s, Chain Type: %s\n",
        response.Data.CCIP.Chain.ID, response.Data.CCIP.Chain.DisplayName, response.Data.CCIP.Chain.Network.Name, response.Data.CCIP.Chain.Network.ChainType)

    return &response, nil
}
