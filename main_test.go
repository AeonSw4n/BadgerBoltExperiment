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
	NumIteratedKeys25MBBatch = 20
	NumBatches1GB            = 40
)

type Key [32]byte

func NewKey(key []byte) Key {
	var k Key
	copy(k[:], key[:])
	return k
}

func (k Key) Bytes() []byte {
	copyKey := make([]byte, len(k[:]))
	copy(copyKey, k[:])
	return copyKey
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
	profiler := NewProfiler()
	for i := 0; i < NumBatches1GB; i++ {
		currentTime := timer.GetTotalElapsedTime("Experiment")
		timer.Start("Experiment")
		Generic25MBBatchTest(db, t)
		timer.End("Experiment")
		deltaTime := timer.GetTotalElapsedTime("Experiment") - currentTime

		memUsageLog := PrintMemUsage()
		profiler.Measure()
		t.Logf("Batch %v/%v\n%v\nTotal time: %vs\t Delta time: %vs\n",
			i, NumBatches1GB, memUsageLog, timer.GetTotalElapsedTime("Experiment"), deltaTime)
	}
	t.Logf("Timer results:\n%v", timer.Print("Experiment"))
	t.Logf("Profiler results:\n%v", profiler.Print())
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

	// Iterate over a couple values from the DB and confirm they match the original values.
	IterateOverDb(db, t)
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
			if err := tx.Set(key.Bytes(), val); err != nil {
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
			if err := tx.Delete(key.Bytes()); err != nil {
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
			val, err := tx.Get(key.Bytes())
			if err != nil {
				val = nil
			}
			values = append(values, val)
		}
		return nil
	}))
	return values
}

func IterateOverDb(db Database, t *testing.T) {
	require := require.New(t)
	require.NoError(db.Update(func(tx Transaction) error {
		it := tx.GetIterator()
		defer it.Close()
		for it.Next() {
			k := it.Key()
			v, err := it.Value()
			require.NoError(err)
			require.NotNil(k)
			require.NotNil(v)
		}
		return nil
	}))
}
