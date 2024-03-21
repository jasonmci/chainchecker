package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
    apiKey       = os.Getenv("COVALENT_API_KEY") // Fetch API Key from environment variable
    accountAddress = "0x1" // Replace with the account address you want to check
)

// Define a struct to unmarshal the API response
type ApiResponse struct {
    Data struct {
        Address string `json:"address"`
        Items   []struct {
            ContractName string `json:"contract_name"`
            ContractTickerSymbol string `json:"contract_ticker_symbol"`
            Balance string `json:"balance"`
            Quote   float64 `json:"quote"`
        } `json:"items"`
    } `json:"data"`
}

func FetchBalances(chainId string) {
    url := fmt.Sprintf("https://api.covalenthq.com/v1/%s/address/%s/balances_v2/?key=%s", chainId, accountAddress, apiKey)
    response, err := http.Get(url)
    if err != nil {
        fmt.Printf("Failed to get response: %v\n", err)
        return
    }
    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil {
        fmt.Printf("Failed to read response body: %v\n", err)
        return
    }

    var apiResponse ApiResponse
    if err := json.Unmarshal(body, &apiResponse); err != nil {
        fmt.Printf("Failed to unmarshal response: %v\n", err)
        return
    }

    fmt.Printf("Token Balances for %s on Chain ID %s:\n", accountAddress, chainId)
    for _, item := range apiResponse.Data.Items {
        fmt.Printf("- %s (%s): %s (Approx. USD Value: $%.2f)\n", item.ContractName, item.ContractTickerSymbol, item.Balance, item.Quote)
    }
}

func main() {
    chainIds := []string{"1", "10", "8453","43114", "42161", "56", "137", "255", "1111"} // Example chain IDs

    for _, chainId := range chainIds {
        FetchBalances(chainId)
		fmt.Println("-----")
    }
}