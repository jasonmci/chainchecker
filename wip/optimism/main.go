package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"

	"github.com/BurntSushi/toml"
)

type TokenConfig struct {
	Token []struct {
		Name            string
		ContractAddress string
		Decimals        int
	}
}

// Response structure to match the JSON response
type ApiResponse struct {
	Status  string
	Message string
	Result  string
}

func getTokenBalance(apiKey, contractAddress, address string, decimals int) *big.Float {
	url := fmt.Sprintf("https://api-optimistic.etherscan.io/api?module=account&action=tokenbalance&contractaddress=%s&address=%s&tag=latest&apikey=%s", contractAddress, address, apiKey)
	responseBody := makeAPICall(url)

	var apiResponse ApiResponse
	if err := json.Unmarshal([]byte(responseBody), &apiResponse); err != nil {
		log.Fatal("Error parsing JSON response:", err)
	}

	balanceWei, ok := big.NewInt(0).SetString(apiResponse.Result, 10)
	if !ok {
		log.Fatal("Error converting balance to big.Int")
	}

	// Convert balance based on token decimals
	divisor := new(big.Float).SetInt(big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	balance := new(big.Float).Quo(new(big.Float).SetInt(balanceWei), divisor)
	return balance
}

// Helper function to make the API call
func makeAPICall(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error making HTTP request:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	return string(body)
}

// Function to get Ether balance
func getEtherBalance(apiKey, address string) *big.Float {
	url := fmt.Sprintf("https://api-optimistic.etherscan.io/api?module=account&action=balance&address=%s&tag=latest&apikey=%s", address, apiKey)
	responseBody := makeAPICall(url)

	var apiResponse ApiResponse
	if err := json.Unmarshal([]byte(responseBody), &apiResponse); err != nil {
		log.Fatal("Error parsing JSON response:", err)
	}

	balanceWei, ok := big.NewInt(0).SetString(apiResponse.Result, 10)
	if !ok {
		log.Fatal("Error converting balance to big.Int")
	}

	// Convert Wei to Ether
	balanceEth := new(big.Float).Quo(new(big.Float).SetInt(balanceWei), big.NewFloat(1e18))
	return balanceEth

}

func main() {
	apiKey := "KGRN86XG15CUE5HWWZGU45SAN9GH9AZ3PD" // Replace with your actual API key
	address := "0x1A2A69e3eB1382FE34Bc579AdD5Bae39e31d4A2c" // Address to check

	var config TokenConfig
	if _, err := toml.DecodeFile("tokens.toml", &config); err != nil {
		log.Fatal(err)
	}

	for _, token := range config.Token {
		balance := getTokenBalance(apiKey, token.ContractAddress, address, token.Decimals)
		fmt.Printf("%s Balance: %s\n", token.Name, balance.Text('f', token.Decimals))
	}
	etherBalance := getEtherBalance(apiKey, address)
	fmt.Printf("Ether Balance: %s\n", etherBalance.Text('f', 18))
}
