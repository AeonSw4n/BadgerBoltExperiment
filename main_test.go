package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	MB25                     = 25000000
	NumKeys25MBBatch         = 250
	NumRemoveKeys25MBBatch   = 20
	NumRetrieveKeys25MBBatch = 20
	NumBatches1GB            = 40
)

type Key [32]byte

func NewKey(key []byte) Key {
	var k Key
	copy(k[:], key[:])
	return k
}

func TestBolt1GB(t *testing.T) {
	require := require.New(t)

	db := NewBoltDatabase()
	require.NoError(db.Setup())
	defer db.Cleanup()
	defer db.Close()
	Generic1GBTest(db, t)
}

func TestBadger1GB(t *testing.T) {
	require := require.New(t)

	db := NewBadgerDatabase()
	require.NoError(db.Setup())
	defer db.Cleanup()
	defer db.Close()
	Generic1GBTest(db, t)
}

func Generic1GBTest(db Database, t *testing.T) {
	timer := NewTimer()
	for i := 0; i < NumBatches1GB; i++ {
		currentTime := timer.GetTotalElapsedTime("Batch")
		timer.Start("Batch")
		Generic25MBBatchTest(db, t)
		timer.End("Batch")
		deltaTime := timer.GetTotalElapsedTime("Batch") - currentTime

		memUsageLog := PrintMemUsage()
		t.Logf("Batch %v/%v\n%v\nTotal time: %vs\t Delta time: %vs\n",
			i, NumBatches1GB, memUsageLog, timer.GetTotalElapsedTime("Batch"), deltaTime)
	}
}

func Generic25MBBatchTest(db Database, t *testing.T) {
	require := require.New(t)

	// Generate 25 MB of data and store it in the DB
	kvMap := Write25MBBatchToDb(db, t)

	// Choose keys to remove and keys to retrieve
	removedKeys := []Key{}
	retrievedKeys := []Key{}
	for key := range kvMap {
		if len(removedKeys) < NumRemoveKeys25MBBatch {
			removedKeys = append(removedKeys, key)
		} else if len(retrievedKeys) < NumRetrieveKeys25MBBatch {
			retrievedKeys = append(retrievedKeys, key)
		} else {
			break
		}
	}

	// Delete keys from DB and confirm they are deleted.
	DeleteFromDB(db, removedKeys, t)
	deletedValues := GetFromDb(db, removedKeys, t)
	for _, val := range deletedValues {
		require.Nil(val)
	}

	// Retrieve keys from DB and confirm they match the original values.
	retrievedValues := GetFromDb(db, retrievedKeys, t)
	for i, val := range retrievedValues {
		require.Equal(kvMap[retrievedKeys[i]], val)
	}
}

func Write25MBBatchToDb(db Database, t *testing.T) (_kv map[Key][]byte) {
	require := require.New(t)
	kvMap := make(map[Key][]byte)
	valueLenght := int32(MB25 / NumKeys25MBBatch)

	for i := 0; i < NumKeys25MBBatch; i++ {
		randomKey, err := RandomBytes(32)
		require.NoError(err)
		keyHash := NewKey(randomKey)
		kvMap[keyHash], err = RandomBytes(valueLenght)
		require.NoError(err)
	}

	require.NoError(db.Update(func(tx Transaction) error {
		for key, val := range kvMap {
			if err := tx.Set(key[:], val); err != nil {
				return err
			}
		}
		return nil
	}))
	return kvMap
}

func DeleteFromDB(db Database, keys []Key, t *testing.T) {
	require := require.New(t)
	require.NoError(db.Update(func(tx Transaction) error {
		for _, key := range keys {
			if err := tx.Delete(key[:]); err != nil {
				return err
			}
		}
		return nil
	}))
}

func GetFromDb(db Database, keys []Key, t *testing.T) [][]byte {
	require := require.New(t)
	var values [][]byte
	require.NoError(db.Update(func(tx Transaction) error {
		for _, key := range keys {
			val, err := tx.Get(key[:])
			if err != nil {
				val = nil
			}
			values = append(values, val)
		}
		return nil
	}))
	return values
}
