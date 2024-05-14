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
)

var (
    a, b string
    printChains bool
    printLanes  bool
)

func init() {
    flag.StringVar(&a, "A", "", "Details for Lane A (network,paymentToken,transferToken)")
    flag.StringVar(&b, "B", "", "Details for Lane B (network,paymentToken,transferToken)")
    flag.BoolVar(&printChains, "c", false, "Print chains")
    flag.BoolVar(&printLanes, "l", false, "Print lanes")
    flag.Parse()
}

func main() {
    
    sessionToken, err := LoginUser(Username, Password)
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    
    response, err := FetchCCIPView(sessionToken)
    if err != nil {
        log.Fatalf("Failed to fetch CCIP view: %v", err)
    }
    
    // Check if we need to print chains
    if printChains {
        PrintChains(response.Data.CCIP.Chains)
    }

    // Check if we need to print lanes
    if printLanes {
        PrintLanes(response.Data.CCIP.Lanes)
    }

    // Get the lane ID as a concatenation of the two network IDs
    partsA := strings.Split(a, ",")
    partsB := strings.Split(b, ",")
    laneID, _ := getLaneID(partsA[0], partsB[0])

    //use networkId map to get the chain ID
    chainAId, _ := getNetworkID(partsA[0])
    chainBId, _ := getNetworkID(partsB[0])

    // Fetch the chain details
    chainA, _   := FetchChainDetails(sessionToken, chainAId)
    chainB, _   := FetchChainDetails(sessionToken, chainBId)

    lane, _     := FetchLaneDetails(sessionToken, laneID)

    // create a TokenStore to hold the token details
    tokenStore := TokenStore{}
    populateConfigFromResponse(*chainA, &tokenStore)
    populateConfigFromResponse(*chainB, &tokenStore)

    feeTokenA, _  := tokenStore.GetTokenDetails(partsA[0], partsA[1])
    feeTokenB, _  := tokenStore.GetTokenDetails(partsB[0], partsB[1])
    txTokenA, _   := tokenStore.GetTokenDetails(partsA[0], partsA[2])
    txTokenB, _   := tokenStore.GetTokenDetails(partsB[0], partsB[2])
    
    testConfig := TestConfig{
        LaneConfigs: make(map[string]LaneConfig),
    }

    DeployedTemplateA := chainA.Data.CCIP.Chain.DeployedTemplate
    DeployedTemplateB := chainA.Data.CCIP.Chain.DeployedTemplate
    
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

    populateLaneConfigFromResponse(lane, chainA, chainB, &testConfig)

    testConfig.UpdateFeeToken(chainA.Data.CCIP.Chain.Network.Name, feeTokenA.Address, true)
    testConfig.AddBridgeToken(chainA.Data.CCIP.Chain.Network.Name, txTokenA.Address, txTokenA.PoolAddress)
    testConfig.UpdateRouter(chainA.Data.CCIP.Chain.Network.Name, chainARouter)
    testConfig.UpdatePriceRegistry(chainA.Data.CCIP.Chain.Network.Name, chainAPriceRegistry)
    testConfig.UpdateArm(chainA.Data.CCIP.Chain.Network.Name, chainAArm)

    testConfig.UpdateFeeToken(chainB.Data.CCIP.Chain.Network.Name, feeTokenB.Address, true)
    testConfig.AddBridgeToken(chainB.Data.CCIP.Chain.Network.Name, txTokenB.Address, txTokenB.PoolAddress)
    testConfig.UpdateRouter(chainB.Data.CCIP.Chain.Network.Name, chainBRouter)
    testConfig.UpdatePriceRegistry(chainB.Data.CCIP.Chain.Network.Name, chainBPriceRegistry)
    testConfig.UpdateArm(chainB.Data.CCIP.Chain.Network.Name, chainBArm)
    
    for templateKey, template := range lane.Data.CCIP.Lane.DeployedTemplate {
		networks := strings.Split(templateKey, "__")

        // separator
        friendlyAName := chainMappings[networks[0]]
        friendlyBName := chainMappings[networks[1]]

        testConfig.AddSrcContract(friendlyAName, friendlyBName, template.OnRampAddress, 55555)
        testConfig.AddDestContract(friendlyBName, friendlyAName, template.OffRampAddress, template.CommitStoreAddress, "111111")
	}

    // Marshal the same data into TOML using BurntSushi/toml
    var tomlBuffer bytes.Buffer
    if tomlErr := toml.NewEncoder(&tomlBuffer).Encode(testConfig); tomlErr != nil {
            fmt.Println("Error marshaling TOML: ", tomlErr)
            return
        }
        //fmt.Println("\nTOML output:")
        //fmt.Println(tomlBuffer.String())

        // write the TOML data to a file
        tomlFile, tomlErr := os.Create("test_config.toml")
        if tomlErr != nil {
            fmt.Println("Error creating TOML file: ", tomlErr)
            return
        }
        tomlFile.WriteString(tomlBuffer.String())
        tomlFile.Close()
        
    
    
    // Marshal the data into JSON
    jsonData, jsonErr := json.MarshalIndent(testConfig, "", "    ")
    if jsonErr != nil {
        fmt.Println("Error marshaling JSON: ", jsonErr)
        return
    }
    fmt.Println("JSON output:")
    fmt.Println(string(jsonData))

    // write the JSON data to a file
    jsonFile, jsonErr := os.Create("test_config.json")
    if jsonErr != nil {
        fmt.Println("Error creating JSON file: ", jsonErr)
        return
    }
    jsonFile.WriteString(string(jsonData))
    jsonFile.Close()
    
}
