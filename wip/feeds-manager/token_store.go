package main

import "fmt"

// Token struct holds information about each token.
type Token struct {
	Name     string // Token name
	Address  string // Token blockchain address
	PoolAddress string // Token pool
	IsFeeToken bool // Is this token used as a fee token
}

// Config holds all chains and provides methods to fetch token details.
type TokenStore struct {
	Chains []Chain
}

// Chain struct represents each blockchain chain and holds its tokens.
type Chain struct {
	Name   string  // Name of the chain
	Tokens []Token // List of tokens associated with this chain
}

// AddToken adds a new token to a specific chain. If the chain does not exist, it adds the chain as well.
func (tokenStore *TokenStore) AddToken(chainName, tokenName, tokenAddress string, poolAddress string, isFeeToken bool) {
	// Check if chain already exists
	for i, chain := range tokenStore.Chains {
		if chain.Name == chainName {
			// Chain exists, add token to this chain
			tokenStore.Chains[i].Tokens = append(tokenStore.Chains[i].Tokens, Token{
				Name:       tokenName,
				Address:    tokenAddress,
				PoolAddress:  poolAddress,
				IsFeeToken: isFeeToken,
			})
			return
		}
	}

	tokenStore.Chains = append(tokenStore.Chains, Chain{
		Name: chainName,
		Tokens: []Token{{
			Name:       tokenName,
			Address:    tokenAddress,
			PoolAddress:  poolAddress,
			IsFeeToken: isFeeToken,
		}},
	})
}

func populateConfigFromResponse(response CCIPChainResponse, store *TokenStore) {
	for tokenName, token := range response.Data.CCIP.Chain.DeployedTemplate.Tokens {
		store.AddToken(response.Data.CCIP.Chain.Network.Name, tokenName, token.TokenAddress, token.PoolAddress, false)
	}
}

// GetTokenDetails retrieves token details from a specific chain.
func (tokenStore *TokenStore) GetTokenDetails(shortChainName, tokenName string) (Token, bool) {
	fullChainName, exists := networkMappings[shortChainName]
	if !exists {
        fmt.Println("Network mapping not found for:", shortChainName)
        return Token{}, false
    }

	for _, chain := range tokenStore.Chains {
		if chain.Name == fullChainName {
			for _, token := range chain.Tokens {
				if token.Name == tokenName {
					return token, true
				}
			}
		}
	}
	return Token{}, false
}

// for _, chain := range tokenStore.Chains {
// 	fmt.Println("Chain:", chain.Name)
// 	for _, token := range chain.Tokens {
// 		fmt.Printf("  Token: %s, Address: %s, PoolAddress: %s, Is Fee Token: %v\n", token.Name, token.Address, token.PoolAddress, token.IsFeeToken)
// 	}
// }