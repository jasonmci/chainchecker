package main

import (
	"github.com/BurntSushi/toml"
)

// TokenMap type for mapping token addresses to names
type TokenMap map[string]string

// LoadTokenMapFromToml loads the token mappings from a TOML file.
func LoadTokenMapFromToml(path string) (TokenMap, error) {
	var tokenMap TokenMap
	if _, err := toml.DecodeFile(path, &tokenMap); err != nil {
		return nil, err
	}
	return tokenMap, nil
}

// GetTokenName looks up the friendly name for a given token address.
func GetTokenName(tokenMap TokenMap, address string) (string, bool) {
	name, exists := tokenMap[address]
	return name, exists
}

