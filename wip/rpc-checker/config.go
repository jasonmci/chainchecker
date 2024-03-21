// config.go

package main

import (
    "github.com/BurntSushi/toml"
)

type Config struct {
    Blockchains 			map[string]BlockchainRPCs	`toml:"blockchains"`
	Accounts map[string]	AccountDetails `toml:"accounts"`
}

type BlockchainRPCs struct {
	HTTP 				string 						`toml:"http"`
	WS   				string 						`toml:"ws"`
}

type Account struct {
	Name 				string 						`toml:"name"`
	Blockchains 		[]string 					`toml:"blockchains"`
	BlockchainDetails  	map[string]Blockchain 		`toml:"blockchainDetails"`
}

type AccountDetails struct {
    Blockchains []string                      `toml:"blockchains"`
    Tokens      map[string][]string           `toml:"tokens"`
}

type Blockchain struct {
	Tokens 				[]string 					`toml:"tokens"`
}

func loadConfig(filename string) (Config, error) {
    var config Config
    if _, err := toml.DecodeFile(filename, &config); err != nil {
        return config, err
    }
    return config, nil
}
