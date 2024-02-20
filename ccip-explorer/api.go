package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

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

func GenerateQueryString(	sender string, receiver string, sourceNetworkName string,
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

func QueryAPI(url string) (*ApiResponse, error) {
	
    apiResponse := ApiResponse{}
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request to API: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return &apiResponse, nil
}
