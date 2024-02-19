package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    client, err := ethclient.Dial("")
    if err != nil {
        log.Fatalf("Failed to connect to the Avalanche network: %v", err)
    }
    fmt.Println("Connected to the Avalanche network")

    // Replace 'YOUR_ACCOUNT_ADDRESS' with the account address you're interested in
    account := common.HexToAddress("")
	blockNumber, err := client.BlockNumber(context.Background())
    if err != nil {
        log.Fatalf("Failed to get the latest block number: %v", err)
    }

	    // Start searching for transactions backward from the latest block
		for i := blockNumber; i > 0; i-- {
			block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
			if err != nil {
				log.Fatalf("Failed to get block: %v", err)
			}
	
			for _, tx := range block.Transactions() {
				if tx.To() != nil && *tx.To() == account {
					// Print details of the transaction
					fmt.Printf("Found transaction: %s\n", tx.Hash().Hex()) // Transaction hash
					fmt.Printf("Block number: %d\n", block.Number().Uint64()) // Block number
					fmt.Printf("Timestamp: %d\n", block.Time()) // Block timestamp in seconds since the epoch
					return
				}
			}
		}
    // Get the nonce for the account, which is the number of transactions sent from the account
    nonce, err := client.PendingNonceAt(context.Background(), account)
    if err != nil {
        log.Fatalf("Failed to fetch the nonce for account %s: %v", account.Hex(), err)
    }

    if nonce == 0 {
        fmt.Println("No transactions found for this account.")
    } else {
        // Fetch the last transaction by subtracting 1 from the nonce
        // Note: This method assumes the account nonce represents the total number of transactions sent from the account. 
        // In practice, you might need to fetch the transaction list and find the latest one by other means.
        fmt.Printf("The last transaction nonce for account %s is %d\n", account.Hex(), nonce-1)
    }
}
