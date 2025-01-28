# tw-parser
Tx parser - Ethereum blockchain parser that will allow to query transactions for subscribed addresses

A main file main.go is provided to run the example. The example will run the parser and print the result to the console. All the parser interface methods are called in the example.

## Interface 
```go
type Parser interface {
        // last parsed block
        GetCurrentBlock() int
        // add address to observer
        Subscribe(address string) bool
        // list of inbound or outbound transactions for an address
        GetTransactions(address string) []Transaction
   }
```
## Running the example
```bash
go run main.go
```

## Notes

When initializing the parser, a scanning process is started in the background. The scanning process will scan the blockchain for new transactions and update the storage. The scanning process will run every 30 miliseconds (a chosen hardcoded variable at the moment. it can be changed and passed as a parameter in the future).

The subscription process only stores the address in the storage. The scanning process will update the transactions for all the addresses found when scanning the blockchain, not only the subscribed addresses.

The parser can be easily switched to processing only the subscribed addresses by adding a filter to the scanning process (`scanBlocksTxs`). This will reduce the amount of data stored in memory and the amount of data processed but depending on the use case, it might not be the best approach.

The inmemory storage is thread safe. It uses `sync.Map` under the hood.
