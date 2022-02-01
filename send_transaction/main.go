package main

import (
	"fmt"
	"os"

	lens "github.com/strangelove-ventures/lens/client"
	registry "github.com/strangelove-ventures/lens/client/chain_registry"
)

func main() {
	var (
		keyName = "default"
	)

	// pull in some chain data
	osmosisInfo, err := registry.DefaultChainRegistry().GetChain("osmosis")
	if err != nil {
		fmt.Printf("Failed to get chain info. Err: %v \n", err)
	}

	rpc, err := osmosisInfo.GetRandomRPCEndpoint()
	if err != nil {
		fmt.Printf("Failed to get random RPC endpoint on chain %s. Err: %v \n", osmosisInfo.ChainID, err)
	}

	osmosisCfg1 := lens.ChainClientConfig{
		Key:            keyName,
		ChainID:        osmosisInfo.ChainID,
		RPCAddr:        rpc,
		AccountPrefix:  osmosisInfo.Bech32Prefix,
		KeyringBackend: "test",
		Debug:          true,
		Timeout:        "10s",
		OutputFormat:   "json",
		SignModeStr:    "direct",
		Modules:        lens.ModuleBasics,
	}

	osmosisClient, err := lens.NewChainClient(&osmosisCfg1, os.Getenv("HOME"), os.Stdin, os.Stdout)
	if err != nil {
		fmt.Printf("Failed to build new chain client for %s. Err: %v \n", osmosisInfo.ChainID, err)
	}

	// setup keys for an account
	keyOutput, err := osmosisClient.AddKey(keyName)
	if err != nil {
		fmt.Printf("Failed to create a new key on chain %s with name %s. Err: %v \n", osmosisClient.ChainId(), keyName, err)
	}

	fmt.Println(osmosisCfg1.KeyDirectory)
	fmt.Println(keyOutput.Address)
	fmt.Println(keyOutput.Mnemonic)

	// build a tx from sender acc
	// send tx to receiver
}
