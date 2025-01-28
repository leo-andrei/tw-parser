package parser

import (
	"context"
	"testing"
	"time"

	"github.com/leo-andrei/tw-parser/internal/ethclient"
)

func TestScannerStartAndStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockClient := newMockEthClient()
	mockStorage := newMockStorage()

	// Setup mock data
	block := &ethclient.Block{
		Number: "0x1",
		Transactions: []ethclient.Transaction{
			{Hash: "tx1", BlockNumber: "0x1", From: "0xfrom1", To: "0xto1"},
		},
	}
	mockClient.blocks[1] = block

	scanner := NewScanner(ctx, mockStorage, mockClient, 1)
	scanner.Start(10 * time.Millisecond)

	// Allow some time for the scanner to process
	time.Sleep(50 * time.Millisecond)

	// Check if the transaction was stored correctly
	txs, err := mockStorage.Get("0xfrom1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(txs) != 1 {
		t.Fatalf("Expected 1 transaction, got %d", len(txs))
	}
	if txs[0].Hash != "tx1" {
		t.Errorf("Expected hash 'tx1', got '%s'", txs[0].Hash)
	}

	// Test stopping the scanner
	cancel()
	time.Sleep(10 * time.Millisecond) // Give some time for the goroutine to exit
}
