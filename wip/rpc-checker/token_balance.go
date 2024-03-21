package main

import (
	"bytes"
	//"fmt"
	"io"
	"net/http"
)

func sendRPCRequest(rpcURL string, payload string) ([]byte, error) {
	response, err := http.Post(rpcURL, "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBytes, nil
}

// func getTokenBalance(rpcURL string, address string, tokenAddress string) (string, error) {
// 	payloadBytes, err := prepareRPCPayload("eth_call")
// 	if err != nil {
// 		return "", err
// 	}

// 	payloadString := fmt.Sprintf(
// 		string(payloadBytes),
// 		tokenAddress,
// 		address,
// 		"latest",
// 	)

// 	responseBytes, err := sendRPCRequest(rpcURL, payloadString)
// 	if err != nil {
// 		return "", err
// 	}

// 	balance, err := parseRPCResponse(responseBytes)
// 	if err != nil {
// 		return "", err
// 	}

// 	return balance, nil
// }