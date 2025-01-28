package storage

type Info struct {
	Name string
}

type Transaction struct {
	BlockNumber string
	Hash        string
	From        string
	To          string
	Type        string
}

type Storage interface {
	Info() (*Info, error)
	Store(key string, value Transaction) error
	Get(key string) ([]Transaction, error)
	Keys() ([]string, error)
}
