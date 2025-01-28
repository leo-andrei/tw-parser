package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"time"

	"github.com/leo-andrei/tw-parser/internal/ethclient"
	"github.com/leo-andrei/tw-parser/internal/storage/inmemory"
	"github.com/leo-andrei/tw-parser/pkg/parser"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	db := inmemory.New()
	fmt.Println("============ Storage Info ============")
	dbInfo, err := db.Info()
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println("type", dbInfo.Name)
	fmt.Println("===================================================")

	ethClient := ethclient.New("https://ethereum-rpc.publicnode.com")
	fmt.Println("============ ETH Client Info ============")
	ethClientInfo, err := ethClient.Info()
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println("apiVersion", ethClientInfo.ApiVersion)
	fmt.Println("endpoint", ethClientInfo.Endpoint)
	fmt.Println("===================================================")

	startBlock := 19912329
	fmt.Println("============ Start at Block ============")
	fmt.Println("number", startBlock)
	fmt.Println("===================================================")

	parser := parser.New(ctx, db, ethClient, startBlock)

	// get the current block number
	fmt.Println("========== Current Block ============")
	blockNumber := parser.GetCurrentBlock()
	fmt.Println("number", blockNumber)
	fmt.Println("===================================================")

	fmt.Println("Waiting 3 seconds for storage to be populated")
	time.Sleep(3 * time.Second)

	// get a random address from the storage
	keys, err := db.Keys()
	if err != nil {
		fmt.Println("keys err", err)
	}
	k, err := rand.Int(rand.Reader, big.NewInt(int64(len(keys)-1)))
	if err != nil {
		fmt.Println("rand err", err)
	}
	key := keys[k.Int64()]

	fmt.Printf("========== Get Transactions (address: %s) ============\n", key)
	txs := parser.GetTransactions(key)
	fmt.Printf("[\n")
	for _, tx := range txs {
		fmt.Printf("{\ntype: %s,\nhash: %s,\nfrom: %s,\nto: %s,\nblock: %s\n}\n", tx.Type, tx.Hash, tx.From, tx.To, tx.BlockNumber)
	}
	fmt.Printf("]")
	fmt.Println("===================================================")

	fmt.Printf("========== Subscribe to Address (address: %s) ============\n", key)
	if ok := parser.Subscribe(key); !ok {
		fmt.Println("subscribe failed")
	} else {
		fmt.Println("Subscribed")
	}
	fmt.Println("===================================================")

	// auto shutdown after 3 seconds
	fmt.Println("Shutting down in 3 seconds...")
	time.Sleep(3 * time.Second)
	signalChan <- os.Interrupt

	<-signalChan
	fmt.Println("Shutting down in progress, cancelling context...")
	cancel()
}
