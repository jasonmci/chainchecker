package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
	
    "github.com/gorilla/websocket"
)

type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
    Method  string        `json:"method"`
    Params  []interface{} `json:"params"`
    ID      int           `json:"id"`
}

// RPCResponse represents a JSON-RPC response payload
type RPCResponse struct {
    ID      int             `json:"id"`
    JSONRPC string          `json:"jsonrpc"`
    Result  json.RawMessage `json:"result"`
    Error   interface{}     `json:"error"`
}

func prepareRPCPayload(method string) ([]byte, error) {
    requestPayload := RPCRequest{
        JSONRPC: "2.0",
        Method:  method,
        Params:  []interface{}{},
        ID:      1,
    }
    return json.Marshal(requestPayload)
}

func parseRPCResponse(responseBytes []byte) (string, error) {
    var rpcResponse RPCResponse
    if err := json.Unmarshal(responseBytes, &rpcResponse); err != nil {
        return "", err
    }

    if rpcResponse.Error != nil {
        return "", fmt.Errorf("RPC Error: %v", rpcResponse.Error)
    }

    var blockNumber string
    if err := json.Unmarshal(rpcResponse.Result, &blockNumber); err != nil {
        return "", err
    }

    return blockNumber, nil
}


func getCurrentBlockNumber(rpcURL string) (string, error) {
    payloadBytes, err := prepareRPCPayload("eth_blockNumber")
    if err != nil {
        return "", err
    }
    body := bytes.NewReader(payloadBytes)

    resp, err := http.Post(rpcURL, "application/json", body)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    responseBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return parseRPCResponse(responseBytes)
}

func getCurrentBlockNumberWS(wsURL string) (string, error) {
    c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        log.Fatal("dial:", err)
        return "", err
    }
    defer c.Close()

    payloadBytes, err := prepareRPCPayload("eth_blockNumber")
    if err != nil {
        return "", err
    }

    err = c.WriteMessage(websocket.TextMessage, payloadBytes)
    if err != nil {
        log.Fatal("write:", err)
        return "", err
    }

    _, responseBytes, err := c.ReadMessage()
    if err != nil {
        log.Fatal("read:", err)
        return "", err
    }

    return parseRPCResponse(responseBytes)
}

func main() {
    checkRPCHealth()
    //getTokenBalance()
}
