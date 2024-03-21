package main

import (
	"fmt"
	"log"
)
	
func checkRPCHealth() () {

    // Load config
    config, err := loadConfig("../secrets.toml")
    if err != nil {
        log.Fatal("Error loading config:", err)
    }
    
    for name, rpcs := range config.Blockchains {
        fmt.Printf("Processing blockchain: %s\n", name)

        // Get current block number using HTTP
        blockNumber, err := getCurrentBlockNumber(rpcs.HTTP)
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }
        fmt.Printf("Current block number (HTTPS) for %s: %s\n", name, blockNumber)

        // Get current block number using WebSocket
        blockNumber, err = getCurrentBlockNumberWS(rpcs.WS)
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }
        fmt.Printf("Current block number (WS) for %s: %s\n", name, blockNumber)
    }


}

