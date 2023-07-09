package main

import (
	"github.com/dgraph-io/badger/v3"
	"log"
	"os"
)

const (
	// PerformanceMemTableSize is 3072 MB. Increases the maximum
	// amount of data we can commit in a single transaction.
	PerformanceMemTableSize = 3072 << 20

	// PerformanceLogValueSize is 256 MB.
	PerformanceLogValueSize = 256 << 20
)

// ==========================
// BadgerDatabase
// ==========================

type BadgerDatabase struct {
	db   *badger.DB
	opts badger.Options
}

func NewBadgerDatabase() *BadgerDatabase {
	opts := GetTestBadgerOpts()
	return &BadgerDatabase{
		db:   nil,
		opts: opts,
	}
}

func (bdb *BadgerDatabase) Setup() error {
	db, err := badger.Open(bdb.opts)
	if err != nil {
		log.Fatal(err)
	}
	bdb.db = db
	return nil
}

func (bdb *BadgerDatabase) Update(fn func(Transaction) error) error {
	return bdb.db.Update(func(txn *badger.Txn) error {
		T := NewBadgerTransaction(txn)
		return fn(T)
	})
}

func (bdb *BadgerDatabase) Close() error {
	return bdb.db.Close()
}

func (bdb *BadgerDatabase) Cleanup() error {
	return os.RemoveAll(bdb.opts.Dir)
}

// ==========================
// BadgerTransaction
// ==========================

type BadgerTransaction struct {
	txn *badger.Txn
}

func NewBadgerTransaction(txn *badger.Txn) *BadgerTransaction {
	return &BadgerTransaction{txn: txn}
}

func (btx *BadgerTransaction) Set(key []byte, value []byte) error {
	return btx.txn.Set(key, value)
}

func (btx *BadgerTransaction) Delete(key []byte) error {
	return btx.txn.Delete(key)
}

func (btx *BadgerTransaction) Get(key []byte) ([]byte, error) {
	var value []byte
	item, err := btx.txn.Get(key)
	if err != nil {
		return value, err
	}
	return item.ValueCopy(nil)
}

func (btx *BadgerTransaction) GetIterator() Iterator {
	opts := badger.DefaultIteratorOptions
	it := btx.txn.NewIterator(opts)
	it.Seek([]byte{})
	return NewBadgerIterator(it)
}

// ==========================
// BadgerIterator
// ==========================

type BadgerIterator struct {
	it *badger.Iterator
}

func NewBadgerIterator(it *badger.Iterator) *BadgerIterator {
	return &BadgerIterator{it: it}
}

func (bit *BadgerIterator) Value() ([]byte, error) {
	item := bit.it.Item()
	return item.ValueCopy(nil)
}

func (bit *BadgerIterator) Key() []byte {
	return bit.it.Item().KeyCopy(nil)
}

func (bit *BadgerIterator) Next() bool {
	bit.it.Next()
	return bit.it.Valid()
}

func (bit *BadgerIterator) Close() {
	bit.it.Close()
}

// PerformanceBadgerOptions are performance geared
// BadgerDB options that use much more RAM than the
// default settings.
func PerformanceBadgerOptions(dir string) badger.Options {
	opts := badger.DefaultOptions(dir)

	// Use an extended table size for larger commits.
	opts.MemTableSize = PerformanceMemTableSize
	opts.ValueLogFileSize = PerformanceLogValueSize

	return opts
}

func GetTestBadgerOpts() badger.Options {
	dir, err := os.MkdirTemp("", "badgerdb")
	if err != nil {
		log.Fatal(err)
	}

	// Open a badgerdb in a temporary directory.
	opts := PerformanceBadgerOptions(dir)
	opts.Dir = dir
	opts.ValueDir = dir
	// Turn off logging for tests.
	opts.Logger = nil

	return opts
}
