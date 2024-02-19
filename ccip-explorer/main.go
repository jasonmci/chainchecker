package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// Define a struct to match the TOML configuration
type Config struct {
	Networks map[string]struct {
		URL    string            `toml:"url"`
		Tokens map[string]string `toml:"tokens"`
	} `toml:"networks"`
}

type ApiResponse struct {
	Data struct {
		AllCcipTransactionsFlats struct {
			Nodes []struct {
				TransactionHash     string `json:"transactionHash"`
				FeeToken            string `json:"feeToken"`
				DestTransactionHash string `json:"destTransactionHash"`
				SourceNetworkName   string `json:"sourceNetworkName"`
				DestNetworkName     string `json:"destNetworkName"`
				MessageID           string `json:"messageId"`
			} `json:"nodes"`
		} `json:"allCcipTransactionsFlats"`
	} `json:"data"`
}

func findTokenNameByAddressInsensitive(tokens map[string]string, address string) (string, bool) {
	for tokenAddress, tokenName := range tokens {
		if strings.EqualFold(tokenAddress, address) {
			return tokenName, true
		}
	}
	return "", false
}

// Reverse lookup to find a token address by its name in a case-insensitive manner
func findTokenAddressByName(networkConfig *Config, networkName, tokenName string) (string, bool) {
	network, exists := networkConfig.Networks[networkName]
	if !exists {
		return "", false // Network not found
	}

	for address, name := range network.Tokens {
		if strings.EqualFold(name, tokenName) {
			return address, true
		}
	}

	return "", false // Token name not found
}


func resolveTokenAddress(networkName, feeToken string, config *Config) string {
	// Check if the feeToken is a known token name and a network name is provided
	if networkName != "" && feeToken != "" {
		// Convert token name to address if possible
		if address, found := findTokenAddressByName(config, networkName, feeToken); found {
			return address // Token name successfully resolved to an address
		}
	}
	// Return the feeToken as is if it's not a recognized token name or no network name is provided
	return feeToken
}

func generateQueryString(	sender string, receiver string, sourceNetworkName string,
							destNetworkName string, messageId string,feeToken string,
							first int, offset int) string {
	baseURL := "https://ccip.chain.link/api/query"
	queryParams := url.Values{}
	queryParams.Add("query", "LATEST_TRANSACTIONS_QUERY")

	// Construct the condition map dynamically based on provided arguments
	condition := map[string]interface{}{
		"sender": sender,
	}
	if receiver != "" {
		condition["receiver"] = receiver
	}
	if sourceNetworkName != "" {
		condition["sourceNetworkName"] = sourceNetworkName
	}
	if destNetworkName != "" {
		condition["destNetworkName"] = destNetworkName
	}
	if messageId != "" {
		condition["messageId"] = messageId
	}
	if feeToken != "" {
		condition["feeToken"] = feeToken
	}

	// Convert the condition map and other variables into a JSON string
	variables := map[string]interface{}{
		"first":     first,
		"offset":    offset,
		"condition": condition,
	}
	variablesJSON, err := json.Marshal(variables)
	if err != nil {
		fmt.Println("Error marshalling variables to JSON:", err)
		return ""
	}

	queryParams.Add("variables", string(variablesJSON))

	return baseURL + "?" + queryParams.Encode()
}

// Function to construct scan URL using the TOML configuration
func constructScanURL(networkName string, transactionHash string, config *Config) string {
	networkConfig, ok := config.Networks[networkName]
	if !ok {
		return "" // Return an empty string for unknown or unsupported networks
	}
	return networkConfig.URL + transactionHash
}

func constructMessageURL(messageID string) string {
	// Construct the message URL using the TOML configuration
	return "https://ccip.chain.link/msg/" + messageID
}

func main() {
	// Define command-line flags
	var (
		sender    = flag.String("sender", "", "Sender address")
		first     = flag.Int("first", 100, "First N results")
		offset    = flag.Int("offset", 0, "Offset for results")
		receiver  = flag.String("receiver", "", "Receiver address")
		source    = flag.String("source", "", "Source network name")
		dest      = flag.String("dest", "", "Destination network name")
		feeToken  = flag.String("feeToken", "", "Fee Token")
		messageId = flag.String("messageId", "", "Message ID")
	)

	// Parse the command-line flags
	flag.Parse()

	// Check if the sender flag was provided
	if *sender == "" {
		fmt.Println("The -sender flag is required.")
		flag.Usage() // Show flag usage
		os.Exit(1)   // Exit the program with an error code
	}

	// Load the TOML configuration
	var config Config
	if _, err := toml.DecodeFile("networks.toml", &config); err != nil {
		// Handle error
		fmt.Println("Error loading TOML configuration:", err)
		return
	}

	// sender := "0x1a2a69e3eb1382fe34bc579add5bae39e31d4a2c"
	// first := 200
	// offset := 0
	// receiver := ""
	// source := "polygon-mainnet"
	// dest := ""
	// feeToken := "LINK"
	// messageId := ""

	resolvedFeeToken := resolveTokenAddress(*source, *feeToken, &config)

	url := generateQueryString(*sender, *receiver, *source, *dest, *messageId, resolvedFeeToken, *first, *offset )
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}

	file, err := os.Create("transactions.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Updated header to include destination fields
	header := []string{	"Source Network Name",  "Dest Network Name", "Fee Token", "Message URL", "Source Scan URL", 
						"Dest Scan URL", "Message ID", "Source Transaction Hash", "Dest Transaction Hash"}
	if err := writer.Write(header); err != nil {
		fmt.Println("Error writing header to CSV file:", err)
		return
	}

	for _, node := range apiResponse.Data.AllCcipTransactionsFlats.Nodes {
		networkName := node.SourceNetworkName
		tokenAddress := node.FeeToken

		networkConfig, exists := config.Networks[networkName]
		if !exists {
			fmt.Println("Network not found in config")
			return
		}
		tokenName, exists := findTokenNameByAddressInsensitive(networkConfig.Tokens, tokenAddress)
		if !exists {
			// fmt.Println("Token not found in network config")
			// return
			tokenName = tokenAddress
		}
	
		sourceScanURL := constructScanURL(node.SourceNetworkName, node.TransactionHash, &config)
		destScanURL := constructScanURL(node.DestNetworkName, node.DestTransactionHash, &config)
		messageURL := constructMessageURL(node.MessageID)
		record := []string{
			node.SourceNetworkName,
			node.DestNetworkName,
			tokenName,
			sourceScanURL,
			destScanURL,
			messageURL,
			node.TransactionHash,
			node.DestTransactionHash,
			node.MessageID,
		}
		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing record to CSV file:", err)
			return
		}
	}

	fmt.Println("CSV file 'transactions.csv' created successfully with scan URLs for both source and destination.")
}
