package main

import (
	"github.com/boltdb/bolt"
	"os"
	"path/filepath"
)

const (
	RandTxnsBucketName = "rand_txns"
)

// ==========================
// BoltDatabase
// ==========================

type BoltDatabase struct {
	db  *bolt.DB
	dir string
}

func NewBoltDatabase() *BoltDatabase {
	return &BoltDatabase{
		db: nil,
	}
}

func (bdb *BoltDatabase) Setup() error {
	dir, err := os.MkdirTemp("", "boltdb")
	if err != nil {
		return err
	}
	bdb.dir = dir
	dbFile := filepath.Join(dir, "bolt.db")

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return err
	}
	bdb.db = db

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(RandTxnsBucketName))
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func (bdb *BoltDatabase) Update(fn func(Transaction) error) error {
	return bdb.db.Update(func(tx *bolt.Tx) error {
		T := NewBoltTransaction(tx)
		return fn(T)
	})
}

func (bdb *BoltDatabase) Close() error {
	return bdb.db.Close()
}

func (bdb *BoltDatabase) Cleanup() error {
	return os.RemoveAll(bdb.dir)
}

// ==========================
// BoltTransaction
// ==========================

type BoltTransaction struct {
	tx *bolt.Tx
}

func NewBoltTransaction(tx *bolt.Tx) *BoltTransaction {
	return &BoltTransaction{
		tx: tx,
	}
}

func (bt *BoltTransaction) Set(key []byte, value []byte) error {
	b := bt.tx.Bucket([]byte(RandTxnsBucketName))
	return b.Put(key, value)
}

func (bt *BoltTransaction) Delete(key []byte) error {
	b := bt.tx.Bucket([]byte(RandTxnsBucketName))
	return b.Delete(key)
}

func (bt *BoltTransaction) Get(key []byte) ([]byte, error) {
	b := bt.tx.Bucket([]byte(RandTxnsBucketName))
	return b.Get(key), nil
}