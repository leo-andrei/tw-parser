package parser

import (
	"context"
	"fmt"
	"testing"

	"github.com/leo-andrei/tw-parser/internal/ethclient"
	"github.com/leo-andrei/tw-parser/internal/storage"
)

type mockStorage struct {
	storage.Storage
	data map[string][]storage.Transaction
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		data: make(map[string][]storage.Transaction),
	}
}

func (m *mockStorage) Store(key string, v storage.Transaction) error {
	m.data[key] = append(m.data[key], v)
	return nil
}

func (m *mockStorage) Get(key string) ([]storage.Transaction, error) {
	if txs, ok := m.data[key]; ok {
		return txs, nil
	}
	return nil, fmt.Errorf("not found")
}

type mockEthClient struct {
	ethclient.Client
	currentBlockNumber int
	blocks             map[int]*ethclient.Block
}

func (m *mockEthClient) BlockNumber() (int, error) {
	return m.currentBlockNumber, nil
}

func newMockEthClient() *mockEthClient {
	return &mockEthClient{
		blocks: make(map[int]*ethclient.Block),
	}
}

func (m *mockEthClient) BlockByNumber(blockNumber int) (*ethclient.Block, error) {
	if block, ok := m.blocks[blockNumber]; ok {
		return block, nil
	}
	return nil, fmt.Errorf("not found")
}

func TestNewParser(t *testing.T) {
	mockClient := newMockEthClient()
	p := New(context.Background(), newMockStorage(), mockClient, 123)
	if p == nil {
		t.Errorf("NewParser() returned nil")
	}
}

func TestGetCurrentBlock(t *testing.T) {
	mockClient := newMockEthClient()
	mockClient.currentBlockNumber = 123
	p := &parser{c: mockClient, storage: newMockStorage()}

	got := p.GetCurrentBlock()
	want := 123
	if got != want {
		t.Errorf("GetCurrentBlock() = %d, want %d", got, want)
	}
}

func TestSubscribe(t *testing.T) {
	mockStorage := newMockStorage()
	p := &parser{c: newMockEthClient(), storage: mockStorage}

	address := "0x123"
	success := p.Subscribe(address)
	if !success {
		t.Errorf("Subscribe() failed")
	}
	if _, ok := mockStorage.data[address]; !ok {
		t.Errorf("Subscribe() did not store the address")
	}
}

func TestGetTransactions(t *testing.T) {
	mockStorage := newMockStorage()
	p := &parser{c: newMockEthClient(), storage: mockStorage}

	address := "0x123"
	tx := storage.Transaction{BlockNumber: "1", Hash: "hash1", From: "from1", To: "to1", Type: "type1"}
	mockStorage.Store(address, tx)

	txs := p.GetTransactions(address)
	if len(txs) != 1 {
		t.Errorf("GetTransactions() returned %d transactions, want 1", len(txs))
	}
	if txs[0].Hash != tx.Hash {
		t.Errorf("GetTransactions() returned transaction hash %s, want %s", txs[0].Hash, tx.Hash)
	}
}
