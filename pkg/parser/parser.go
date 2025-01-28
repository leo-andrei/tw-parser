package parser

import (
	"context"
	"log"
	"time"

	"github.com/leo-andrei/tw-parser/internal/ethclient"
	"github.com/leo-andrei/tw-parser/internal/storage"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
}

type Transaction struct {
	BlockNumber string `json:"blockNumber"`
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Type        string `json:"type"`
}

type parser struct {
	c       ethclient.Client
	storage storage.Storage
}

func New(ctx context.Context, db storage.Storage, ethClient ethclient.Client, startBlock int) Parser {
	bs := NewScanner(ctx, db, ethClient, startBlock)
	bs.Start(30 * time.Millisecond)

	return &parser{c: ethClient, storage: db}
}

// GetCurrentBlock returns the current block number from the eth client
func (s *parser) GetCurrentBlock() int {
	blockNumber, err := s.c.BlockNumber()
	if err != nil {
		log.Fatalf("failed to get block number: %v", err)
	}
	return blockNumber
}

// Subscribe adds an address to the storage to be monitored for transactions
func (s *parser) Subscribe(address string) bool {
	if err := s.storage.Store(address, storage.Transaction{}); err != nil {
		return false
	}
	return true
}

// GetTransactions returns all transactions for a given address from the storage
func (s *parser) GetTransactions(address string) []Transaction {
	txs, err := s.storage.Get(address)
	if err != nil {
		return []Transaction{}
	}

	res := []Transaction{}
	for _, tx := range txs {
		res = append(res, Transaction{
			BlockNumber: tx.BlockNumber,
			Hash:        tx.Hash,
			From:        tx.From,
			To:          tx.To,
			Type:        tx.Type,
		})
	}

	return res
}
