package inmemory

import (
	"testing"

	"github.com/leo-andrei/tw-parser/internal/storage"
)

func TestInfo(t *testing.T) {
	db := New()
	info, err := db.Info()
	if err != nil {
		t.Errorf("Info() error = %v, wantErr %v", err, nil)
	}
	if info.Name != "inmemory" {
		t.Errorf("Info() Name = %v, want %v", info.Name, "inmemory")
	}
}

func TestStoreAndGet(t *testing.T) {
	db := New()
	key := "testKey"
	tx := storage.Transaction{Hash: "tx1", From: "addr1", To: "addr2", BlockNumber: "0x4bad55"}

	// test Store
	err := db.Store(key, tx)
	if err != nil {
		t.Errorf("Store() error = %v, wantErr %v", err, nil)
	}

	// test Get
	txs, err := db.Get(key)
	if err != nil {
		t.Errorf("Get() error = %v, wantErr %v", err, nil)
	}
	if len(txs) != 1 {
		t.Errorf("Get() txs len = %d, want %d", len(txs), 1)
	}
	if txs[0] != tx {
		t.Errorf("Get() tx = %v, want %v", txs[0], tx)
	}

	// test Store by adding a new tx
	tx = storage.Transaction{Hash: "tx1", From: "addr1", To: "addr4", BlockNumber: "0x4bad56"}
	err = db.Store(key, tx)
	if err != nil {
		t.Errorf("Store() error = %v, wantErr %v", err, nil)
	}
}

func TestGetNotFound(t *testing.T) {
	db := New()
	_, err := db.Get("nonexistent")
	if err == nil {
		t.Errorf("Get() error = %v, wantErr %v", err, "key not found")
	}
}

func TestKeys(t *testing.T) {
	db := New()
	key1 := "testKey1"
	tx1 := storage.Transaction{Hash: "tx1", From: "addr1", To: "addr2", BlockNumber: "0x4bad55"}
	key2 := "testKey2"
	tx2 := storage.Transaction{Hash: "tx2", From: "addr3", To: "addr4", BlockNumber: "0x4bad56"}

	_ = db.Store(key1, tx1)
	_ = db.Store(key2, tx2)

	keys, err := db.Keys()
	if err != nil {
		t.Errorf("Keys() error = %v, wantErr %v", err, nil)
	}
	if len(keys) != 2 {
		t.Errorf("Keys() keys len = %d, want %d", len(keys), 2)
	}
	if !contains(keys, key1) || !contains(keys, key2) {
		t.Errorf("Keys() keys = %v, want %v and %v", keys, key1, key2)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
