package main

import (
	"fmt"
	"sort"
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


// Network mappings from short name to full name
var networkMappings = map[string]string{
    "KROMA": "Kroma Mainnet",
    "WEMIX": "WeMix Mainnet",
    "GNO":   "GnosisChain Mainnet",
    "OPT":   "Optimism Mainnet",
    "AVAX":  "Avalanche Mainnet",
    "POLY":  "Polygon Mainnet",
    "BSC":   "BSC Mainnet",
    "ETH":   "Ethereum Mainnet",
    "ARB":   "Arbitrum Mainnet",
    "BASE":  "Base Mainnet",
}

var chainMappings = map[string]string{
    "ethereum-mainnet-arbitrum-1":  "Arbitrum Mainnet",
    "ethereum-mainnet-base-1":      "Base Mainnet",
    "ethereum-mainnet-optimism-1":  "Optimism Mainnet",
    "matic-mainnet":                "Polygon Mainnet",
    "ethereum-mainnet":             "Ethereum Mainnet",
    "avalanche-mainnet":            "Avalanche Mainnet",
    "bsc-mainnet":                  "BSC Mainnet",
    "ethereum-mainnet-kroma-1":     "Kroma Mainnet",
    "wemix-mainnet":                "WeMix Mainnet",
    "xdai-mainnet":                 "GnosisChain Mainnet",
}

var networkIDMap = map[string]string{
    "OPT": "18",
    "AVAX": "19",
    "POLY": "20",
    "ETH": "21",
    "BSC": "22",
    "ARB": "23",
    "BASE": "24",
    "WEMIX": "27",
    "KROMA": "28",
    "GNO": "30",
}

// map friendly network name to its native fee token address
var nativeFeeTokenMap = map[string]string{
    "Optimism Mainnet": "0x4200000000000000000000000000000000000006",
    "Avalanche Mainnet": "0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7",
    "Polygon Mainnet": "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
    "Ethereum Mainnet": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
    "BSC Mainnet": "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c",
    "Arbitrum Mainnet": "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
    "Base Mainnet": "0x4200000000000000000000000000000000000006",
    "WeMix Mainnet": "0x7D72b22a74A216Af4a002a1095C8C707d6eC1C5f",
    "Kroma Mainnet": "0x4200000000000000000000000000000000000001",
    "GnosisChain Mainnet": "0xe91D153E0b41518A2Ce8Dd3D7944Fa863463a97d",
}

// PrintNetworkMappings prints all network mappings from the map
func PrintNetworkMappings() {
    fmt.Println("Network Mappings:")
    for shortName, fullName := range networkMappings {
        fmt.Printf("%-10s%s\n", shortName, fullName)
    }
}

func getNetworkID(networkName string) (string, bool) {
    id, found := networkIDMap[networkName]
    return id, found
}

func normalizePair(net1, net2 string) (string, string) {
    networks := []string{net1, net2}
    sort.Strings(networks)  // Sorts the slice in ascending order.
    return networks[0], networks[1]
}

func getLaneID(net1, net2 string) (string, bool) {
    net1, net2 = normalizePair(net1, net2)
    id, found := laneIDMap[net1+","+net2]
    return id, found
}