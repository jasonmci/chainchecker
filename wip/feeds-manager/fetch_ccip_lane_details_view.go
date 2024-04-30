package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	//"fmt"

	//"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type CCIPLaneResponse struct {
    Data struct {
        CCIP struct {
            Lane CCIPLane `json:"lane"` // check on this
        } `json:"ccip"`
    } `json:"data"`
}
type CCIPLane struct {
    ID                        string     `json:"id"`
    DisplayName               string     `json:"displayName"`
    Status                    string     `json:"status"`
    DeployedTemplate          LaneDeployedTemplate     `json:"deployedTemplate"` // Assuming string; adjust if it's more complex
	DeployedProvisionalTemplate  map[string]interface{} `json:"deployedProvisionalTemplate"`
    LegA                      CCIPLaneLeg `json:"legA"`
    LegB                      CCIPLaneLeg `json:"legB"`
    LegAProvisional           CCIPLaneLeg `json:"legAProvisional"`
    LegBProvisional           CCIPLaneLeg `json:"legBProvisional"`
    ChainA                    CCIPChain   `json:"chainA"`
    ChainB                    CCIPChain   `json:"chainB"`
}

type CCIPLaneLeg struct {
    ID             string         `json:"id"`
    Tag            string         `json:"tag"`
    Status         string         `json:"status"`
    WorkflowRuns   []WorkflowRun  `json:"workflowRuns"`
    Source         CCIPEndpoint   `json:"source"`
    Destination    CCIPEndpoint   `json:"destination"`
    Dons           []DON          `json:"dons"`
}

type WorkflowRun struct {
    ID           string   `json:"id"`
    Name         string   `json:"name"`
    Status       string   `json:"status"`
    WorkflowType string   `json:"workflowType"`
    CreatedAt    string   `json:"createdAt"`
    Actions      []Action `json:"actions"`
}

type Action struct {
    ActionType string     `json:"actionType"`
    Name       string     `json:"name"`
    Run        ActionRun  `json:"run"`
    Tasks      []Task     `json:"tasks"`
}

type ActionRun struct {
    ID         string   `json:"id"`
    Status     string   `json:"status"`
    Network    Network  `json:"network"`
    CreatedAt  string   `json:"createdAt"`
}

type Task struct {
    Name string   `json:"name"`
    Run  TaskRun  `json:"run"`
}

type TaskRun struct {
    Error   string `json:"error"`
    ID      string `json:"id"`
    Input   string `json:"input"`
    Output  string `json:"output"`
    Status  string `json:"status"`
    TxHash  string `json:"txHash"`
}

type CCIPEndpoint struct {
    Chain     CCIPChain   `json:"chain"`
    Contracts []Contract  `json:"contracts"`
}

type CCIPChain struct {
    ID           string   `json:"id"`
    DisplayName  string   `json:"displayName"`
    Network      Network  `json:"network"`
    Contracts    []Contract `json:"contracts"`
    WorkflowRuns []WorkflowRun `json:"workflowRuns"`
}

type LaneContract struct {
    ID                      string   `json:"id"`
    Address                 string   `json:"address"`
    Tag                     string   `json:"tag"`
    TransferOwnershipStatus string   `json:"transferOwnershipStatus"`
    Name                    string   `json:"name"`
    Semver                  string   `json:"semver"`
    Metadata                string   `json:"metadata"` // Assuming string; adjust if it's more complex
    OwnerType               string   `json:"ownerType"`
    OwnerAddress            string   `json:"ownerAddress"`
    PendingOwnerAddress     string   `json:"pendingOwnerAddress"`
    PendingOwnerType        string   `json:"pendingOwnerType"`
    Network                 Network  `json:"network"`
}

type LaneNetwork struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    IconName    string `json:"iconName"`
    ExplorerURL string `json:"explorerURL"`
    ChainID     string `json:"chainID"`
    ChainType   string `json:"chainType"`
}

type DON struct {
    ID           string `json:"id"`
    ExecutionType string `json:"executionType"`
    Jobs         []Job  `json:"jobs"`
}

type Job struct {
    ID              string `json:"id"`
    Status          string `json:"status"`
    IsBootstrap     bool   `json:"isBootstrap"`
    CanPropose      bool   `json:"canPropose"`
    CanRevoke       bool   `json:"canRevoke"`
    ProposalChanged bool   `json:"proposalChanged"`
    NodeOperator    struct {
        ID   string `json:"id"`
        Name string `json:"name"`
    } `json:"nodeOperator"`
    Node struct {
        ID   string `json:"id"`
        Name string `json:"name"`
    } `json:"node"`
    AssignableNodes []struct {
        ID   string `json:"id"`
        Name string `json:"name"`
    } `json:"assignableNodes"`
}

type LaneDeployedTemplate struct {
    Deployments map[string]DeploymentDetails `json:"deployedTemplate"`
}

type DeploymentDetails struct {
    CommitDON           []string         `json:"commitDON"`
    CommitStore         CommitStore      `json:"commitStore"`
    CommitStoreAddress  string           `json:"commitStoreAddress"`
    ExecuteDON          []string         `json:"executeDON"`
    OffRamp             OffRamp          `json:"offRamp"`
    OffRampAddress      string           `json:"offRampAddress"`
    OnRamp              OnRamp           `json:"onRamp"`
    OnRampAddress       string           `json:"onRampAddress"`
    SupportedTokens     []string         `json:"supportedTokens"`
}

type CommitStore struct {
    DynamicConfig   DynamicConfig   `json:"dynamicConfig"`
    MinSeqNr        string          `json:"minSeqNr"`
    Ocr2Config      Ocr2Config      `json:"ocr2Config"`
    StaticConfig    StaticConfig    `json:"staticConfig"`
}

type DynamicConfig struct {
    PriceRegistry string `json:"priceRegistry"`
}

type Ocr2Config struct {
    F                int                  `json:"f"`
    OffchainConfig   OffchainConfig       `json:"offchainConfig"`
    OffchainConfigVersion string          `json:"offchainConfigVersion"`
    Oracles          []string             `json:"oracles"`
}

type OffchainConfig struct {
    DeltaGraceNanoseconds                       string `json:"deltaGraceNanoseconds"`
    DeltaProgressNanoseconds                    string `json:"deltaProgressNanoseconds"`
    DeltaResendNanoseconds                      string `json:"deltaResendNanoseconds"`
    DeltaRoundNanoseconds                       string `json:"deltaRoundNanoseconds"`
    DeltaStageNanoseconds                       string `json:"deltaStageNanoseconds"`
    MaxDurationObservationNanoseconds           string `json:"maxDurationObservationNanoseconds"`
    MaxDurationQueryNanoseconds                 string `json:"maxDurationQueryNanoseconds"`
    MaxDurationReportNanoseconds                string `json:"maxDurationReportNanoseconds"`
    MaxDurationShouldAcceptFinalizedReportNanoseconds string `json:"maxDurationShouldAcceptFinalizedReportNanoseconds"`
    MaxDurationShouldTransmitAcceptedReportNanoseconds string `json:"maxDurationShouldTransmitAcceptedReportNanoseconds"`
    RMax             int                        `json:"rMax"`
    ReportingPluginConfig ReportingPluginConfig `json:"reportingPluginConfig"`
    S                []int                      `json:"s"`
}

type ReportingPluginConfig struct {
    DestFinalityDepth int    `json:"DestFinalityDepth"`
    FeeUpdateDeviationPPB uint64 `json:"FeeUpdateDeviationPPB"`
    FeeUpdateHeartBeat string `json:"FeeUpdateHeartBeat"`
    InflightCacheExpiry string `json:"InflightCacheExpiry"`
    MaxGasPrice string `json:"MaxGasPrice"`
    SourceFinalityDepth int `json:"SourceFinalityDepth"`
}

type StaticConfig struct {
    ArmProxy string `json:"armProxy"`
    ChainSelector string `json:"chainSelector"`
    OnRamp string `json:"onRamp"`
    SourceChainSelector string `json:"sourceChainSelector"`
}

type OffRamp struct {
    DynamicConfig DynamicConfig `json:"dynamicConfig"`
    Ocr2Config    Ocr2Config    `json:"ocr2Config"`
    RateLimiterConfig RateLimiterConfig `json:"rateLimiterConfig"`
    SourceTokens  map[string]string `json:"sourceTokens"`
    StaticConfig  StaticConfig  `json:"staticConfig"`
    Tokens        []string      `json:"tokens"`
}

type OnRamp struct {
    AllowList             []string             `json:"allowList"`
    AllowListEnabled      bool                 `json:"allowListEnabled"`
    DynamicConfig         DynamicConfig        `json:"dynamicConfig"`
    FeeTokenConfigArgs    map[string]FeeTokenConfigArgs `json:"feeTokenConfigArgs"`
    NopsAndWeights        map[string]int       `json:"nopsAndWeights"`
    RateLimiterConfig     RateLimiterConfig    `json:"rateLimiterConfig"`
    StaticConfig          StaticConfig         `json:"staticConfig"`
    TokenTransferFeeConfigArgs map[string]TokenTransferFeeConfigArgs `json:"tokenTransferFeeConfigArgs"`
    Tokens                []string             `json:"tokens"`
}

type LaneRateLimiterConfig struct {
    Admin  string `json:"admin"`
    Config struct {
        Capacity  string `json:"capacity"`
        IsEnabled bool   `json:"isEnabled"`
        Rate      string `json:"rate"`
    } `json:"config"`
}

type FeeTokenConfigArgs struct {
    DestGasOverhead          int    `json:"destGasOverhead"`
    DestGasPerPayloadByte    int    `json:"destGasPerPayloadByte"`
    Enabled                  bool   `json:"enabled"`
    GasMultiplier            string `json:"gasMultiplier"`
    NetworkFeeAmountUSD      string `json:"networkFeeAmountUSD"`
}

type TokenTransferFeeConfigArgs struct {
    MaxFee  uint32 `json:"maxFee"`
    MinFee  uint32 `json:"minFee"`
    Ratio   int    `json:"ratio"`
}

func FetchLaneDetails(sessionToken, laneID string ) []byte {

	queryBytes, err := os.ReadFile("FetchCCIPLaneDetailsView.graphql")
    if err != nil {
        log.Fatalf("Failed to read GraphQL file: %v", err)
    }
    query := string(queryBytes)

    payload := map[string]interface{}{
        "operationName": "FetchCCIPLaneDetailsView",
        "variables": map[string]interface{}{
            "id": laneID,
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
    var response CCIPLaneResponse
    err = json.Unmarshal(responseBody, &response)
    if err != nil {
        log.Fatalf("Error parsing JSON response: %v", err)
    }
    // var resp2 LaneDetailsResponse
    // if err := json.Unmarshal(responseBody, &resp2); err != nil {
    //     log.Fatalf("Error parsing JSON: %v", err)
    // }

    lane := response.Data.CCIP.Lane
    fmt.Printf("Lane ID: %s, Status: %s\n", lane.ID, lane.Status)

    // Example of printing details for legA
    fmt.Printf("Leg A ID: %s, Status: %s\n", lane.LegA.ID, lane.LegA.Status)
    for _, workflow := range lane.LegA.WorkflowRuns {
        fmt.Printf("Workflow Run ID: %s, Name: %s, Status: %s\n", workflow.ID, workflow.Name, workflow.Status)
    }

    return responseBody
}
