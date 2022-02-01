package main

import (
	"fmt"
	"os"

	lens "github.com/strangelove-ventures/lens/client"
	registry "github.com/strangelove-ventures/lens/client/chain_registry"
)

func main() {

	// Only need these two pieces of info to implement this examples:

	// Let us the chain Registry to query info about the chain.
	// To utalize this feature this string matches one of the directories here: https://github.com/cosmos/chain-registry:
	var chainRegName = "cosmoshub"
	var walletAddress = "cosmos1fjw9hvt7dulewrn4e65u6f39arexhyk5umj0cq"

	// Fetches chain info from chain registry
	chainInfo, err := registry.DefaultChainRegistry().GetChain(chainRegName)
	if err != nil {
		fmt.Printf("Failed to get chain info. Err: %v \n", err)
	}

	rpc, err := chainInfo.GetRandomRPCEndpoint()
	if err != nil {
		fmt.Printf("Failed to get RPC endpoints on chain %s. Err: %v \n", chainInfo.ChainName, err)
	}

	// For this simple example, only two fields are needed
	chainConfig := lens.ChainClientConfig{
		// Key            string,
		// ChainID        string,
		RPCAddr: rpc,
		// GRPCAddr       string,
		// AccountPrefix  string,
		KeyringBackend: "test",
		// GasAdjustment  float64,
		// GasPrices      string,
		// KeyDirectory   string,
		// Debug          bool,
		// Timeout        string,
		// OutputFormat   string,
		// SignModeStr    string,
		// Modules        []module.AppModuleBasic,
	}

	chainClient, err := lens.NewChainClient(&chainConfig, os.Getenv("HOME"), os.Stdin, os.Stdout)
	if err != nil {
		fmt.Printf("Failed to build new chain client for %s. Err: %v \n", chainInfo.ChainID, err)
	}

	balance, err := chainClient.QueryBalanceWithAddress(walletAddress)
	if err != nil {
		fmt.Printf("Failed to query balance. Err: %v", err)
	}

	fmt.Println(balance)
}
