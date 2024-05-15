package main

import (
	"log"
	"sort"

	"github.com/BurntSushi/toml"
)

type Contract struct {
    ID                    string `json:"id"`
    Name                  string `json:"name"`
    Address               string `json:"address"`
    Tag                   string `json:"tag"`
    TransferOwnershipStatus string `json:"transferOwnershipStatus"`
    TypeName              string `json:"__typename"`
}

type Network struct {
    ID            string `json:"id"`
    Name          string `json:"name"`
    IconName      string `json:"iconName"`
    ExplorerURL   string `json:"explorerURL"`
    TypeName      string `json:"__typename"`
    ChainID       string `json:"chainID"`
    ChainType     string `json:"chainType"`
}

type GenConfig struct {
    LaneIDMap map[string]string `toml:"laneIDMap"`
    ChainIDMap map[string]string `toml:"chainIDMap"`
    ChainMapping map[string]string `toml:"chainMappings"`
    ShortcutMappings map[string]string `toml:"shortcutMappings"`
    NativeFeeTokenMap map[string]string `toml:"nativeFeeTokenMap"`
}

// loadConfig loads the GenConfig from a TOML file.
func loadGeneratorConfig() GenConfig {
    var genConfig GenConfig
    if _, err := toml.DecodeFile("config.toml", &genConfig); err != nil {
        log.Fatalf("Error loading config: %v", err)
    }
    return genConfig
}

func getFullChainName (genConfig GenConfig, chainName string) string {
    // if the chainName is not found in the shortcutMappings, return the chainName itself
    if _, found := genConfig.ShortcutMappings[chainName]; !found {
        return chainName
    }
    return genConfig.ShortcutMappings[chainName]
}


func getNetworkID(genConfig GenConfig, networkName string) (string, bool) {
    id, found := genConfig.ChainIDMap[networkName]
    return id, found
}


func normalizePair(net1, net2 string) (string, string) {
    networks := []string{net1, net2}
    sort.Strings(networks)  // Sorts the slice in ascending order.
    return networks[0], networks[1]
}

func getLaneID(genConfig GenConfig, net1, net2 string) (string, bool) {
    laneIDMap := genConfig.LaneIDMap
    net1, net2 = normalizePair(net1, net2)
    id, found := laneIDMap[net1+","+net2]
    return id, found
}

func getChainMapping(genConfig GenConfig, chainName string) (string, bool) {
    chainMapping := genConfig.ChainMapping
    id, found := chainMapping[chainName]
    return id, found
}
