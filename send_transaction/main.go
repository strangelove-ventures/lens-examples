package main

import (
	"context"
	"fmt"
	"log"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	lens "github.com/strangelove-ventures/lens/client"
	registry "github.com/strangelove-ventures/lens/client/chain_registry"
)

// Needed vars for this example:

var (
	// We will be fetching chain info from chain registry.
	// This string must match the relevant directory name in the chain registry here:
	// https://github.com/cosmos/chain-registry
	chainRegName = "osmosis"

	srcWalletAddress   = "osmo1fjw9hvt7dulewrn4e65u6f39arexhyk55qplwj"
	srcWalletMnemonic  = os.Getenv("testKeyMn")
	destination_wallet = ""
	amount_to_send     = "100000uosmo"
)

func main() {

	//	Fetches chain info from chain registry
	chainInfo, err := registry.DefaultChainRegistry().GetChain(chainRegName)
	if err != nil {
		log.Fatalf("Failed to get chain info. Err: %v \n", err)
	}

	//	Use Chain info to select random endpoint
	rpc, err := chainInfo.GetRandomRPCEndpoint()
	if err != nil {
		log.Fatalf("Failed to get random RPC endpoint on chain %s. Err: %v \n", chainInfo.ChainID, err)
	}

	// For this example, lets place the key directory in your PWD.
	pwd, _ := os.Getwd()
	key_dir := pwd + "/keys"

	// Build chain config
	chainConfig_1 := lens.ChainClientConfig{
		Key:     "default",
		ChainID: chainInfo.ChainID,
		RPCAddr: rpc,
		// GRPCAddr       string,
		AccountPrefix:  chainInfo.Bech32Prefix,
		KeyringBackend: "test",
		GasAdjustment:  1.2,
		GasPrices:      "0.01uosmo",
		KeyDirectory:   key_dir,
		Debug:          true,
		Timeout:        "20s",
		OutputFormat:   "json",
		SignModeStr:    "direct",
		Modules:        lens.ModuleBasics,
	}

	// Creates client object to pull chain info
	chainClient, err := lens.NewChainClient(&chainConfig_1, key_dir, os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalf("Failed to build new chain client for %s. Err: %v \n", chainInfo.ChainID, err)
	}

	// Lets restore a key with funds and name it "source_key", this is the wallet we'll use to send tx.
	src_address, err := chainClient.RestoreKey("source_key", srcWalletMnemonic)
	if err != nil {
		log.Fatalf("Failed to restore key. Err: %v \n", err)
	}

	//	Now that we know our key name, we can set it in our chain config
	chainConfig_1.Key = "source_key"

	// Sanitize coin amount and make it readable by SDK
	coins, err := sdk.ParseCoinNormalized(amount_to_send)
	if err != nil {
		log.Fatalf("Error parsing coin string. Error: %s", err)
	}

	//	Build transaction message
	req := &banktypes.MsgSend{
		FromAddress: chainClient.MustEncodeAccAddr(sdk.AccAddress(src_address)),
		ToAddress:   chainClient.MustEncodeAccAddr(sdk.AccAddress(destination_wallet)),
		Amount:      coins,
	}

	// Send message and get response
	res, err := chainClient.SendMsg(context.Background(), req)
	if err != nil {
		if res != nil {
			log.Fatalf("failed to send coins: code(%d) msg(%s)", res.Code, res.Logs)
		}
		log.Fatalf("Failed to send coins.Err: %w", err)
	}
	fmt.Println(chainClient.PrintTxResponse(res))
}
