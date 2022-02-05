package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go"
	lens "github.com/strangelove-ventures/lens/client"
	registry "github.com/strangelove-ventures/lens/client/chain_registry"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"golang.org/x/sync/errgroup"
)

var (
	RtyAttNum = uint(5)
	RtyAtt    = retry.Attempts(RtyAttNum)
	RtyDel    = retry.Delay(time.Millisecond * 400)
	RtyErr    = retry.LastErrorOnly(true)
)

func main() {

	// Fetches chain info from chain registry
	chainInfo, err := registry.DefaultChainRegistry().GetChain("osmosis")
	if err != nil {
		log.Fatalf("failed to get chain info. err: %v", err)
	}

	//	Use Chain info to select random endpoint
	rpc, err := chainInfo.GetRandomRPCEndpoint()
	if err != nil {
		log.Fatalf("failed to get RPC endpoints on chain %s. err: %v", chainInfo.ChainName, err)
	}

	// Creates client object to pull chain info
	chainClient, err := lens.NewChainClient(&lens.ChainClientConfig{RPCAddr: rpc, KeyringBackend: "test"}, os.Getenv("HOME"), os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalf("failed to build new chain client for %s. err: %v", chainInfo.ChainID, err)
	}

	// create the database connection
	db, err := ConnectToDatabase("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to database. err: %v", err)
	}

	// create the indexer object
	i := &Indexer{
		Client: chainClient,
		DB:     db,
	}

	// run the indexer
	if err := i.ForEachBlock([]int64{}, i.IndexIBCTransactions, 100); err != nil {
		log.Fatalf("failed to index blocks. err: %v", err)
	}
}

type Indexer struct {
	Client *lens.ChainClient
	DB     *sql.DB

	logger tmlog.Logger
}

func (i *Indexer) ForEachBlock(blocks []int64, cb func(height int64) error, concurrentBlocks int) error {
	fmt.Println("starting block queries for", i.Client.Config.ChainID)
	var (
		eg           errgroup.Group
		mutex        sync.Mutex
		failedBlocks = make([]int64, 0)
		sem          = make(chan struct{}, concurrentBlocks)
	)
	for _, h := range blocks {
		h := h
		sem <- struct{}{}

		eg.Go(func() error {
			if err := cb(h); err != nil {
				if strings.Contains(err.Error(), "wrong ID: no ID") {
					mutex.Lock()
					failedBlocks = append(failedBlocks, h)
					mutex.Unlock()
				} else {
					return fmt.Errorf("[height %d] - failed to get block. err: %s", h, err.Error())
				}
			}
			<-sem
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	if len(failedBlocks) > 0 {
		return i.ForEachBlock(failedBlocks, cb, concurrentBlocks)
	}
	return nil
}

func (i *Indexer) LogRetryGetBlock(n uint, err error, h int64) {
	i.logger.Error("retry", "attempt", n, "err", err, "height", h)
}

func ConnectToDatabase(driver, connString string) (*sql.DB, error) {
	db, err := sql.Open(driver, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open db, ensure db server is running & check conn string. err: %s", err.Error())
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to db, ensure db server is running & check conn string. err: %s", err.Error())
	}
	return db, nil
}
