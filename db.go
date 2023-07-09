package main

type Database interface {
	Setup() error
	Update(func(Transaction) error) error
	Close() error
	Cleanup() error
}

type Transaction interface {
	Set(key []byte, value []byte) error
	Delete(key []byte) error
	Get(key []byte) ([]byte, error)
}
