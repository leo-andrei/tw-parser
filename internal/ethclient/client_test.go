package ethclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBlockNumber(t *testing.T) {
	// setup a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x4bad55"}`))
	}))
	defer ts.Close()

	client := New(ts.URL)

	blockNumber, err := client.BlockNumber()
	if err != nil {
		t.Fatalf("BlockNumber() error = %v, wantErr %v", err, nil)
	}
	expectedBlockNumber := 4959573 // 0x4bad55 in decimal
	if blockNumber != expectedBlockNumber {
		t.Errorf("BlockNumber() = %v, want %v", blockNumber, expectedBlockNumber)
	}
}

func TestBlockByNumber(t *testing.T) {
	// setup a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":{"number":"0x4bad55","transactions":[{"hash":"0xhash1","blockNumber":"0x5bad55","from":"0xfrom","to":"0xto"}]}}`))
	}))
	defer ts.Close()

	client := New(ts.URL)

	block, err := client.BlockByNumber(4959573)
	if err != nil {
		t.Fatalf("BlockByNumber() error = %v, wantErr %v", err, nil)
	}
	if block.Number != "0x4bad55" {
		t.Errorf("BlockByNumber() Number = %v, want %v", block.Number, "0x4bad55")
	}
	if len(block.Transactions) != 1 {
		t.Errorf("BlockByNumber() Transactions len = %d, want %d", len(block.Transactions), 1)
	}
	if block.Transactions[0].Hash != "0xhash1" {
		t.Errorf("BlockByNumber() Transaction Hash = %v, want %v", block.Transactions[0].Hash, "0xhash1")
	}
}

func TestInfo(t *testing.T) {
	endpoint := "http://localhost:8545"
	client := New(endpoint)

	info, err := client.Info()
	if err != nil {
		t.Fatalf("Info() error = %v, wantErr %v", err, nil)
	}
	if info.ApiVersion != DefaultApiVersion {
		t.Errorf("Info() ApiVersion = %v, want %v", info.ApiVersion, DefaultApiVersion)
	}
	if info.Endpoint != endpoint {
		t.Errorf("Info() Endpoint = %v, want %v", info.Endpoint, endpoint)
	}
}
