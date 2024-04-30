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
    "ARB,BASE": "37",
    "ETH,OPT": "38",
    "AVAX,ETH": "39",
    "ETH,POLY": "40",
    "AVAX,POLY": "41",
    "BSC,ETH": "42",
    "ARB,ETH": "43",
    "BASE,ETH": "44",
    "BASE,OPT": "46",
    "AVAX,BSC": "47",
    "BSC,POLY": "48",
    "OPT,POLY": "49",
    "BASE,BSC": "50",
    "ETH,WEMIX": "51",
    "KROMA,WEMIX": "52",
    "ARB,POLY": "53",
    "ARB,BSC": "54",
    "ARB,OPT": "55",
    "ARB,AVAX": "56",
    "AVAX,OPT": "57",
    "BASE,POLY": "58",
    "BSC,OPT": "59",
    "AVAX,BASE": "60",
    "BSC,WEMIX": "61",
    "AVAX,WEMIX": "62",
    "POLY,WEMIX": "63",
    "ARB,WEMIX": "64",
    "OPT,WEMIX": "65",
    "ETH,GNO": "68",
    "GNO,POLY": "69",
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
