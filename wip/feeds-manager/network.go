package main

import (
    "fmt"
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

// type Chain struct {
//     ID        string    `json:"id"`
//     DisplayName string  `json:"displayName"`
//     Network   Network   `json:"network"`
//     Contracts []Contract `json:"contracts"`
//     TypeName  string    `json:"__typename"`
// }

// type CCIP struct {
//     Chains []Chain `json:"chains"`
// }

// type Data struct {
//     CCIP CCIP `json:"ccip"`
// }

// type Response struct {
//     Data Data `json:"data"`
// }

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
