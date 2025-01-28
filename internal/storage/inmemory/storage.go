package inmemory

import (
	"fmt"
	"sync"

	"github.com/leo-andrei/tw-parser/internal/storage"
)

type db struct {
	data sync.Map
}

func New() storage.Storage {
	return &db{
		data: sync.Map{},
	}
}

func (s *db) Info() (*storage.Info, error) {
	return &storage.Info{
		Name: "inmemory",
	}, nil
}

func (s *db) Store(key string, v storage.Transaction) error {
	val, ok := s.data.Load(key)
	if ok {
		s.data.Store(key, append(val.([]storage.Transaction), v))
	} else {
		s.data.Store(key, []storage.Transaction{v})
	}
	return nil
}

func (s *db) Get(key string) ([]storage.Transaction, error) {
	value, ok := s.data.Load(key)
	if !ok {
		return nil, fmt.Errorf("key not found")
	}
	return value.([]storage.Transaction), nil
}

func (s *db) Keys() ([]string, error) {
	keys := []string{}
	s.data.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys, nil
}
