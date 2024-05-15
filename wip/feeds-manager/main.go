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

func main() {

    flag.Parse()
    genConfig := loadGeneratorConfig()

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

    if a == "" || b == "" {
        return
    }

    // Get the lane ID as a concatenation of the two network IDs
    partsA := strings.Split(a, ",")
    partsB := strings.Split(b, ",")

    // Convert shortcut names to full names if needed
    chainAFullName := getFullChainName(genConfig, partsA[0])
    chainBFullName := getFullChainName(genConfig, partsB[0])

    laneID, _ := getLaneID(genConfig, chainAFullName, chainBFullName)

    chainAId, _ := getNetworkID(genConfig, chainAFullName)
    chainBId, _ := getNetworkID(genConfig, chainBFullName)
    
    // Fetch the chain details
    chainA, _   := FetchChainDetails(sessionToken, chainAId)
    chainB, _   := FetchChainDetails(sessionToken, chainBId)

    // CLO Has it as WEMIX Mainnet, but the test framework expects WeMix Mainnet`
    if chainA.Data.CCIP.Chain.Network.Name == "WEMIX Mainnet" {
        chainA.Data.CCIP.Chain.Network.Name = "WeMix Mainnet"
    }
    if chainB.Data.CCIP.Chain.Network.Name == "WEMIX Mainnet" {
        chainB.Data.CCIP.Chain.Network.Name = "WeMix Mainnet"
    }

    lane, _     := FetchLaneDetails(sessionToken, laneID)

    // create a TokenStore to hold the token details
    tokenStore := TokenStore{}
    populateConfigFromResponse(*chainA, &tokenStore)
    populateConfigFromResponse(*chainB, &tokenStore)

    feeTokenA, _  := tokenStore.GetTokenDetails(genConfig, chainAFullName, partsA[1])
    feeTokenB, _  := tokenStore.GetTokenDetails(genConfig, chainBFullName, partsB[1])
    txTokenA, _   := tokenStore.GetTokenDetails(genConfig, chainAFullName, partsA[2])
    txTokenB, _   := tokenStore.GetTokenDetails(genConfig, chainBFullName, partsB[2])
    
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

    populateLaneConfigFromResponse(genConfig, lane, chainA, chainB, &testConfig)

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
        friendlyAName, found := getChainMapping(genConfig, networks[0])
        if !found {
            fmt.Println("Error getting chain mapping: ", err)
            return
        }

        friendlyBName, found := getChainMapping(genConfig, networks[1])
        if !found {
            fmt.Println("Error getting chain mapping: ", err)
            return
        }
        
        fmt.Print("Friendly A Name: ", networks[0], friendlyAName)
        fmt.Print("Friendly B Name: ", networks[1], friendlyBName)

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
