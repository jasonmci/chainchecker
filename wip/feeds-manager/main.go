package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
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

    var a, b string
    flag.StringVar(&a, "A", "", "Details for Lane A (network,paymentToken,transferToken)")
    flag.StringVar(&b, "B", "", "Details for Lane B (network,paymentToken,transferToken)")
    flag.Parse()

    // Setup command line flags
    //listNetworks := flag.Bool("list", false, "List all networks")
    //shortName := flag.String("network", "", "Short network name to fetch details for")
    //flag.Parse()

    // Get the lane ID as a concatenation of the two network IDs
    partsA := strings.Split(a, ",")
    partsB := strings.Split(b, ",")
    laneID, _ := getLaneID(partsA[0], partsB[0])

    //use networkId map to get the chain ID
    chainAId, _ := getNetworkID(partsA[0])
    chainBId, _ := getNetworkID(partsB[0])

    fmt.Printf("chain ID: %s, other chain ID: %s\n", partsA[0], partsB[0])
    fmt.Printf("Lane ID: %s\n", laneID)


    token, err := LoginUser(username, password)
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    
    chainA, _ := FetchChainDetails(token, chainAId)
    chainB, _ := FetchChainDetails(token, chainBId)

    //lane, _ := FetchLaneDetails(token, laneID)
    
    //DeployedTemplateLane := lane.Data.CCIP.Lane.DeployedTemplate
    DeployedTemplateA := chainA.Data.CCIP.Chain.DeployedTemplate
    
    SupportedTokensA := chainA.Data.CCIP.Chain.SupportedTokens
    SupportedTokensB := chainB.Data.CCIP.Chain.SupportedTokens

    DeployedTemplateB := chainB.Data.CCIP.Chain.DeployedTemplate

    var chainAArm string
    for address := range DeployedTemplateA.Arms {
        chainAArm = address
    }

    var chainARouter string
    for address := range DeployedTemplateA.Routers {
        chainARouter = address
    }

    var chainAPriceRegistry string
    for address := range DeployedTemplateA.PriceRegistries {
        chainAPriceRegistry = address
    }

    var tokenALINK string
    var tokenANative string

    for _, token := range SupportedTokensA {
        if token.Token == "LINK" {
            tokenALINK = token.Address
        }
        if token.Token == "WETH" {
            tokenANative = token.Address
        }
    }

    var tokenBLINK string
    var tokenBNative string

    for _, token := range SupportedTokensB {
        if token.Token == "LINK" {
            tokenBLINK = token.Address
        }
    }

    var chainBArm string
    for address := range DeployedTemplateB.Arms {
        chainBArm = address
    }

    var chainBRouter string
    for address := range DeployedTemplateB.Routers {
        chainBRouter = address
    }

    var chainBPriceRegistry string
    for address := range DeployedTemplateB.PriceRegistries {
        chainBPriceRegistry = address
    }

    laneConfig := TestConfig{
        LaneConfigs: map[string]LaneConfig{
            chainA.Data.CCIP.Chain.Network.Name : {
                IsNativeFeeToken: true,
                FeeToken: tokenALINK,
                BridgeTokens: []string{""},
                BridgeTokensPools: []string{""},
                Arm: chainAArm,
                Router: chainARouter,
                PriceRegistry: chainAPriceRegistry,
                WrappedNative: tokenANative,
                SrcContracts: map[string]SrcContract{
                    chainB.Data.CCIP.Chain.Network.Name: {OnRamp: "", DeployedAt: 11111111},
                },
                DestContracts: map[string]DestContract{
                    chainB.Data.CCIP.Chain.Network.Name: {OffRamp: "", CommitStore: "", ReceiverDapp: ""},
                },
            },
            chainB.Data.CCIP.Chain.Network.Name : {
                IsNativeFeeToken: true,
                FeeToken: tokenBLINK,
                BridgeTokens: []string{""},
                BridgeTokensPools: []string{""},
                Arm: chainBArm,
                Router: chainBRouter,
                PriceRegistry: chainBPriceRegistry,
                WrappedNative: tokenBNative,
                SrcContracts: map[string]SrcContract{
                    chainA.Data.CCIP.Chain.Network.Name: {OnRamp: "", DeployedAt: 11111111},
                },
                DestContracts: map[string]DestContract{
                    chainA.Data.CCIP.Chain.Network.Name: {OffRamp: "", CommitStore: "", ReceiverDapp: ""},
                },
            },
        },
    }

       
        // Marshal the data into JSON
        jsonData, err := json.MarshalIndent(laneConfig, "", "    ")
        if err != nil {
            fmt.Println("Error marshaling JSON: ", err)
            return
        }
    
        // Print the JSON string
        fmt.Println(string(jsonData))



        // Marshal the data into JSON
    jsonData, jsonErr := json.MarshalIndent(laneConfig, "", "    ")
    if jsonErr != nil {
        fmt.Println("Error marshaling JSON: ", jsonErr)
        return
    }
    fmt.Println("JSON output:")
    fmt.Println(string(jsonData))

    // Marshal the same data into TOML using BurntSushi/toml
    var tomlBuffer bytes.Buffer
    if tomlErr := toml.NewEncoder(&tomlBuffer).Encode(laneConfig); tomlErr != nil {
        fmt.Println("Error marshaling TOML: ", tomlErr)
        return
    }
    fmt.Println("\nTOML output:")
    fmt.Println(tomlBuffer.String())

    // fmt.Println("Chain A Response:", myresp.Data.CCIP.Chain.DisplayName)

    // fmt.Printf("Chain ID: %s, Display Name: %s, Network: %s, Chain Type: %s\n",
    // myresp.Data.CCIP.Chain.ID, myresp.Data.CCIP.Chain.DisplayName, myresp.Data.CCIP.Chain.Network.Name, myresp.Data.CCIP.Chain.Network.ChainType)

}
