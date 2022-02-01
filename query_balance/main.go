package main

import (
	"fmt"
	"os"

	lens "github.com/strangelove-ventures/lens/client"
	registry "github.com/strangelove-ventures/lens/client/chain_registry"
)

func main() {

	// For this example, we only need two pices of info:
	type walletInfo struct {
		// Lets use the Chain Registry to automatically get accurate chain info
		// This string must match the relevant directory name in the chain registry here:
		// https://github.com/cosmos/chain-registry
		chainRegName  string
		walletAddress string
	}

	wallet_1 := walletInfo{"<Chain-Registry-Name>", "<Wallet Address>"}

	// Fetches chain info from chain registry
	chainInfo, err := registry.DefaultChainRegistry().GetChain(wallet_1.chainRegName)
	if err != nil {
		fmt.Printf("Failed to get chain info. Err: %v \n", err)
	}

	rpc, err := chainInfo.GetRandomRPCEndpoint()
	if err != nil {
		fmt.Printf("Failed to get RPC endpoints on chain %s. Err: %v \n", chainInfo.ChainName, err)
	}

	// For this simple example, only two fields from the config are needed
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

	// Creates client object to pull chain info
	chainClient, err := lens.NewChainClient(&chainConfig, os.Getenv("HOME"), os.Stdin, os.Stdout)
	if err != nil {
		fmt.Printf("Failed to build new chain client for %s. Err: %v \n", chainInfo.ChainID, err)
	}

	balance, err := chainClient.QueryBalanceWithAddress(wallet_1.walletAddress)
	if err != nil {
		fmt.Printf("Failed to query balance. Err: %v", err)
	}

	fmt.Println(balance)
}
