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

func normalizePair(net1, net2 string) (string, string) {
    networks := []string{net1, net2}
    sort.Strings(networks)  // Sorts the slice in ascending order.
    return networks[0], networks[1]
}

var laneIDMap = map[string]string{
    "ARB,ETH": "43",
    "ETH,WEMIX": "51",
    // Add more mappings as needed
}

func getLaneID(net1, net2 string) (string, bool) {
    net1, net2 = normalizePair(net1, net2)
    id, found := laneIDMap[net1+","+net2]
    return id, found
}


// Network mappings from short name to full name
var networkMappings = map[string]string{
    "KROMA": "Kroma Mainnet",
    "WEMIX": "WEMIX Mainnet",
    "GNO":   "GnosisChain Mainnet",
    "POLYX": "Polygon zkEVM Mainnet",
    "OPT":   "Optimism Mainnet",
    "AVAX":  "Avalanche Mainnet",
    "POLY":  "Polygon Mainnet",
    "BSC":   "BSC Mainnet",
    "ETH":   "Ethereum Mainnet",
    "ARB":   "Arbitrum Mainnet",
    "BASE":  "Base Mainnet",
}

// PrintNetworkMappings prints all network mappings from the map
func PrintNetworkMappings() {
    fmt.Println("Network Mappings:")
    for shortName, fullName := range networkMappings {
        fmt.Printf("%-10s%s\n", shortName, fullName)
    }
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
    "POLYZ": "31",
}

func getNetworkID(networkName string) (string, bool) {
    id, found := networkIDMap[networkName]
    return id, found
}


// func fetchNetworkDetails(data *Data, networkName string) {
//     found := false
//     for _, chain := range data.CCIP.Chains {
//         if chain.Network.Name == networkName {
//             found = true
//             fmt.Printf("Network Name: %s\n", chain.Network.Name)
//             fmt.Printf("Explorer URL: %s\n", chain.Network.ExplorerURL)

//             for _, contract := range chain.Contracts {
//                 switch contract.Name {
//                 case "Router":
//                     fmt.Printf("Router Address: %s\n", contract.Address)
//                 case "PriceRegistry":
//                     fmt.Printf("PriceRegistry Address: %s\n", contract.Address)
//                 case "ARMContract":
//                     fmt.Printf("ARMContract Address: %s\n", contract.Address)
//                 }
//             }
//             break
//         }
//     }
//     if !found {
//         fmt.Println("Network not found.")
//     }
// }
