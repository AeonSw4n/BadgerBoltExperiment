package main

type Database interface {
	Setup() error
	Update(func(Transaction) error) error
	Close() error
	Cleanup() error
	Id() string
}

type Transaction interface {
	Set(key []byte, value []byte) error
	Delete(key []byte) error
	Get(key []byte) ([]byte, error)
	GetIterator() Iterator
}

type Iterator interface {
	Value() ([]byte, error)
	Key() []byte
	Next() bool
	Close()
}
