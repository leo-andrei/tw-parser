package parser

import (
	"context"
	"fmt"
	"time"

	"github.com/leo-andrei/tw-parser/internal/ethclient"
	"github.com/leo-andrei/tw-parser/internal/storage"
)

type scanner struct {
	ctx              context.Context
	db               storage.Storage
	ethclient        ethclient.Client
	startBlock       int
	lastScannedBlock int
}

func NewScanner(ctx context.Context, db storage.Storage, ethclient ethclient.Client, startBlock int) *scanner {
	return &scanner{
		ctx:              ctx,
		db:               db,
		ethclient:        ethclient,
		startBlock:       startBlock,
		lastScannedBlock: startBlock - 1,
	}
}

func (bs *scanner) Start(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-bs.ctx.Done():
				fmt.Println("backgroundScan: context done")
				return
			case <-ticker.C:
				ticker.Stop()

				// block scan
				blockToScan := bs.lastScannedBlock + 1
				// get block by number
				block, err := bs.ethclient.BlockByNumber(blockToScan)
				if err != nil {
					fmt.Println("backgroundScan: error getting block by number", err)
					return
				}

				txsToStore := bs.scanBlockTxs(block)
				if err := bs.storeTxs(txsToStore); err != nil {
					fmt.Println("backgroundScan: error storing transaction", err)
					continue
				}
				bs.lastScannedBlock = blockToScan

				ticker.Reset(interval)
			}
		}
	}()
}

func (bs *scanner) scanBlockTxs(block *ethclient.Block) map[string][]storage.Transaction {
	txsToStore := map[string][]storage.Transaction{}
	for _, tx := range block.Transactions {
		// save transaction to db
		txToStore := storage.Transaction{
			Hash:        tx.Hash,
			BlockNumber: tx.BlockNumber,
			From:        tx.From,
			To:          tx.To,
		}

		txToStore.Type = "outbound"
		txsToStore[tx.From] = []storage.Transaction{txToStore}
		txToStore.Type = "inbound"
		txsToStore[tx.To] = []storage.Transaction{txToStore}
	}
	return txsToStore
}

func (bs *scanner) storeTxs(txs map[string][]storage.Transaction) error {
	for address, transactions := range txs {
		for _, tx := range transactions {
			if err := bs.db.Store(address, tx); err != nil {
				continue
			}
		}
	}
	return nil
}
