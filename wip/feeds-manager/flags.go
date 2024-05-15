package main

import "flag"


var (
    a, b string
    printChains bool
    printLanes  bool
)

func init() {
    flag.StringVar(&a, "A", "", "Details for Lane A (network,paymentToken,transferToken)")
    flag.StringVar(&b, "B", "", "Details for Lane B (network,paymentToken,transferToken)")
    flag.BoolVar(&printChains, "c", false, "Print chains")
    flag.BoolVar(&printLanes, "l", false, "Print lanes")

}

// Define a struct to hold lane information
type ChainInfo struct {
    Network string
    PaymentTokens []string
    TransferTokens []string
}

type LaneInfo struct {
    Network string
    PaymentTokens []string
    TransferTokens []string
}