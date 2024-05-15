// network_test.go
package main

import (
	"bytes"
	"testing"
)

// contains checks if a string is contained within another string.
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

// TestGetNetworkID checks if getNetworkID returns the correct ID and found status for various network names.
func TestGetNetworkID(t *testing.T) {
	genConfig := loadGeneratorConfig()
	tests := []struct {
		networkName string
		expectedID  string
		expectedFound bool
	}{
		{"Optimism Mainnet", "18", true},
		{"Avalanche Mainnet", "19", true},
		{"UNKNOWN", "", false},
	}

	for _, test := range tests {
		id, found := getNetworkID(genConfig, test.networkName)
		if id != test.expectedID || found != test.expectedFound {
			t.Errorf("getNetworkID(%q) = %q, %v; want %q, %v", test.networkName, id, found, test.expectedID, test.expectedFound)
		}
	}
}

// TestNormalizePair checks if normalizePair returns the network names in sorted order.
func TestNormalizePair(t *testing.T) {
	tests := []struct {
		net1, net2 string
		expected1, expected2 string
	}{
		{"ETH", "BSC", "BSC", "ETH"},
		{"POLY", "AVAX", "AVAX", "POLY"},
		{"BASE", "OPT", "BASE", "OPT"},
	}

	for _, test := range tests {
		n1, n2 := normalizePair(test.net1, test.net2)
		if n1 != test.expected1 || n2 != test.expected2 {
			t.Errorf("normalizePair(%q, %q) = %q, %q; want %q, %q", test.net1, test.net2, n1, n2, test.expected1, test.expected2)
		}
	}
}

// TestGetLaneID checks if getLaneID returns the correct lane ID and found status for various network pairs.
func TestGetLaneID(t *testing.T) {
	genConfig := loadGeneratorConfig()
	tests := []struct {
		net1, net2 string
		expectedID string
		expectedFound bool
	}{
		{"Arbitrum Mainnet", "Base Mainnet", "37", true},
		{"Ethereum Mainnet", "Polygon Mainnet", "40", true},
		{"UNKNOWN1", "UNKNOWN2", "", false},
	}

	for _, test := range tests {
		id, found := getLaneID(genConfig, test.net1, test.net2)
		if id != test.expectedID || found != test.expectedFound {
			t.Errorf("getLaneID(%q, %q) = %q, %v; want %q, %v", test.net1, test.net2, id, found, test.expectedID, test.expectedFound)
		}
	}
}

func TestGetChainMapping(t *testing.T) {
	genConfig := loadGeneratorConfig()
	tests := []struct {
		networkName string
		expectedChain string
	}{
		{"ethereum-mainnet-arbitrum-1", "Arbitrum Mainnet"},
		{"ethereum-mainnet-base-1", "Base Mainnet"},
		{"ethereum-mainnet-optimism-1", "Optimism Mainnet"},
		{"matic-mainnet", "Polygon Mainnet"},
		{"ethereum-mainnet", "Ethereum Mainnet"},
		{"avalanche-mainnet", "Avalanche Mainnet"},
		{"bsc-mainnet", "BSC Mainnet"},
		{"ethereum-mainnet-kroma-1", "Kroma Mainnet"},
		{"wemix-mainnet", "WeMix Mainnet"},
		{"xdai-mainnet", "GnosisChain Mainnet"},
	}

	for _, test := range tests {
		chain, _ := getChainMapping(genConfig, test.networkName)
		if chain != test.expectedChain  {
			t.Errorf("getChainMapping(%q) = %q; want %q", test.networkName, chain, test.expectedChain)
		}
	}
}