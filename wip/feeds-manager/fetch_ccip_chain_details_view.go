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
    queryBytes, err := os.ReadFile("queries/FetchCCIPChainDetailsView.graphql")
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
    var response CCIPChainResponse
    err = json.Unmarshal(responseBody, &response)
    if err != nil {
        return nil, fmt.Errorf("error parsing JSON response: %w", err)
    }

    return &response, nil
}
