package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

func setCommonHeaders(req *http.Request, sessionToken string) {
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "*/*")
    req.Header.Set("X-Session-Token", sessionToken)
}

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
    
    token, err := LoginUser(username, password)
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    
    // Get the lane ID as a concatenation of the two network IDs
    partsA := strings.Split(a, ",")
    partsB := strings.Split(b, ",")
    laneID, _ := getLaneID(partsA[0], partsB[0])

    //use networkId map to get the chain ID
    chainAId, _ := getNetworkID(partsA[0])
    chainBId, _ := getNetworkID(partsB[0])

    // Fetch the chain details
    chainA, _ := FetchChainDetails(token, chainAId)
    chainB, _ := FetchChainDetails(token, chainBId)
    lane, _ := FetchLaneDetails(token, laneID)

    // create a TokenStore to hold the token details
    tokenStore := TokenStore{}
    populateConfigFromResponse(*chainA, &tokenStore)
    populateConfigFromResponse(*chainB, &tokenStore)
    
	// Displaying updated configuration
	for _, chain := range tokenStore.Chains {
		fmt.Println("Chain:", chain.Name)
		for _, token := range chain.Tokens {
			fmt.Printf("  Token: %s, Address: %s, PoolAddress: %s, Is Fee Token: %v\n", token.Name, token.Address, token.PoolAddress, token.IsFeeToken)
		}
	}

    feeTokenA, _  := tokenStore.GetTokenDetails(partsA[0], partsA[1])
    feeTokenB, _  := tokenStore.GetTokenDetails(partsB[0], partsB[1])
    txTokenA, _   := tokenStore.GetTokenDetails(partsA[0], partsA[2])
    txTokenB, _   := tokenStore.GetTokenDetails(partsB[0], partsB[2])


    fmt.Printf("Token Address: %s, Pool Address %s\n", feeTokenA.Address, feeTokenA.PoolAddress)
    fmt.Printf("Token Address: %s, Pool Address %s\n", feeTokenB.Address, feeTokenB.PoolAddress)
    fmt.Printf("Token Address: %s, Pool Address %s\n", txTokenA.Address, txTokenA.PoolAddress)
    fmt.Printf("Token Address: %s, Pool Address %s\n", txTokenB.Address, txTokenB.PoolAddress)

    fmt.Printf("New stuff Lane ID: %s, Display Name: %s, Status: %s\n", lane.Data.CCIP.Lane.ID, lane.Data.CCIP.Lane.DisplayName, lane.Data.CCIP.Lane.Status)

    //DeployedTemplateLane := lane.Data.CCIP.Lane.DeployedTemplate
    DeployedTemplateA := chainA.Data.CCIP.Chain.DeployedTemplate
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
                FeeToken: feeTokenA.Address,
                BridgeTokens: []string{txTokenA.Address},
                BridgeTokensPools: []string{txTokenA.PoolAddress},
                Arm: chainAArm,
                Router: chainARouter,
                PriceRegistry: chainAPriceRegistry,
                WrappedNative: feeTokenA.Address,
                SrcContracts: map[string]SrcContract{
                    chainB.Data.CCIP.Chain.Network.Name: {OnRamp: "", DeployedAt: 11111111},
                },
                DestContracts: map[string]DestContract{
                    chainB.Data.CCIP.Chain.Network.Name: {OffRamp: "", CommitStore: "", ReceiverDapp: ""},
                },
            },
            chainB.Data.CCIP.Chain.Network.Name : {
                IsNativeFeeToken: true,
                FeeToken: feeTokenB.Address,
                BridgeTokens: []string{txTokenB.Address},
                BridgeTokensPools: []string{txTokenB.PoolAddress},
                Arm: chainBArm,
                Router: chainBRouter,
                PriceRegistry: chainBPriceRegistry,
                WrappedNative: feeTokenB.Address,
                SrcContracts: map[string]SrcContract{
                    chainA.Data.CCIP.Chain.Network.Name: {OnRamp: "", DeployedAt: 11111111},
                },
                DestContracts: map[string]DestContract{
                    chainA.Data.CCIP.Chain.Network.Name: {OffRamp: "", CommitStore: "", ReceiverDapp: ""},
                },
            },
        },
    }

    // Update the OnRamp for ChainB under ChainA's SrcContracts
    laneConfig.LaneConfigs[chainA.Data.CCIP.Chain.Network.Name ].SrcContracts[chainB.Data.CCIP.Chain.Network.Name ] = SrcContract{
        OnRamp:     "newOnRampAddress",
        DeployedAt: 11111111,
    }
        // } else {
        //     fmt.Println("ChainA configuration not found.")
        // }

    // the following are some examples of changing the laneConfig data
    // change IsNativeFeeToken to false
    laneConfig.LaneConfigs[chainA.Data.CCIP.Chain.Network.Name] = LaneConfig{
        IsNativeFeeToken: false,
    }

    // add a third lane
    laneConfig.LaneConfigs["third"] = LaneConfig{
        IsNativeFeeToken: true,
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

}
