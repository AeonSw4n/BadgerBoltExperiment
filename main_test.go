package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type TestConfig struct {
	// ExperimentNumberOfBatches is the number of batches to write to the database.
	ExperimentNumberOfBatches int
	// BatchSizeBytes is the size of each batch in bytes.
	BatchSizeBytes int
	// BatchSizeItems is the number of items in each batch.
	BatchSizeItems int
	// BatchItemsRemoved is the number of items removed from each batch as part of the experiment.
	BatchItemsRemoved int
	// BatchItemsRetrieved is the number of items retrieved from each batch as part of the experiment.
	BatchItemsRetrieved int
	// BatchItemsIterated is the number of items iterated over in each batch as part of the experiment.
	BatchItemsIterated int
}

// TestBolt_5GB_Experiment_10MB_Batch is a BoltDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 10MB, totalling 500 batches. For each batch we:
// 		1. Write 100 equal size KV items to the database.
// 		2. Remove 20 items from the database.
// 		3. Retrieve 20 items from the database.
// 		4. Iterate over  items in the database.
func TestBolt_5GB_Experiment_10MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 500,
		BatchSizeBytes:            10000000,
		BatchSizeItems:            100,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}

	dir, err := os.MkdirTemp("", "boltdb-5gb-10mb")
	t.Logf("BoltDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	db := NewBoltDatabase(dir)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	ctx := db.GetContext([]byte("TestBucket"))
	GenericTest(db, ctx, testConfig, t)

	// Test nested context
	boltCtx, err := AssertContext[*BoltContext](ctx, BOLTDB)
	require.NoError(err)
	boltNestedCtx := NewBoltNestedContext([]byte("NestedBucket"), boltCtx)
	p := NewProfiler()
	timer := NewTimer()
	GenericBatchTest(db, boltNestedCtx, p, timer, testConfig, t)
}

// TestBolt_5GB_Experiment_25MB_Batch is a BoltDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 25MB, totalling 200 batches. For each batch we perform identical
// operations as in TestBolt_5GB_Experiment_10MB_Batch.
func TestBolt_5GB_Experiment_25MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 200,
		BatchSizeBytes:            25000000,
		BatchSizeItems:            250,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}

	dir, err := os.MkdirTemp("", "boltdb-5gb-25mb")
	t.Logf("BoltDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	db := NewBoltDatabase(dir)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	ctx := db.GetContext([]byte("TestBucket"))
	GenericTest(db, ctx, testConfig, t)
}

// TestBolt_5GB_Experiment_100MB_Batch is a BoltDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 100MB, totalling 50 batches. For each batch we perform identical
// operations as in TestBolt_5GB_Experiment_10MB_Batch.
func TestBolt_5GB_Experiment_100MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 50,
		BatchSizeBytes:            100000000,
		BatchSizeItems:            1000,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}

	dir, err := os.MkdirTemp("", "boltdb-5gb-100mb")
	t.Logf("BoltDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	db := NewBoltDatabase(dir)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	ctx := db.GetContext([]byte("TestBucket"))
	GenericTest(db, ctx, testConfig, t)
}

// TestBadger_Default_5GB_Experiment_10MB_Batch is a BadgerDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 10MB, totalling 500 batches. For each batch we perform identical
// operations as in TestBolt_5GB_Experiment_10MB_Batch. This experiment uses the "Default" Badger config.
func TestBadger_Default_5GB_Experiment_10MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 500,
		BatchSizeBytes:            10000000,
		BatchSizeItems:            100,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}
	dir, err := os.MkdirTemp("", "badgerdb-default-5gb-10mb")
	t.Logf("BadgerDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	opts := DefaultBadgerOptions(dir)
	db := NewBadgerDatabase(opts, false)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	badgerCtx := db.GetContext([]byte{})
	GenericTest(db, badgerCtx, testConfig, t)
}

// TestBadger_Performance_5GB_Experiment_10MB_Batch is a BadgerDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 10MB, totalling 500 batches. For each batch we perform identical
// operations as in TestBolt_5GB_Experiment_10MB_Batch. This experiment uses the "Performance" Badger config.
func TestBadger_Performance_5GB_Experiment_10MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 500,
		BatchSizeBytes:            10000000,
		BatchSizeItems:            100,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}
	dir, err := os.MkdirTemp("", "badgerdb-performance-5gb-10mb")
	t.Logf("BadgerDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	opts := PerformanceBadgerOptions(dir)
	db := NewBadgerDatabase(opts, false)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	badgerCtx := db.GetContext([]byte{})
	GenericTest(db, badgerCtx, testConfig, t)
}

// TestBadger_Performance_5GB_Experiment_25MB_Batch is a BadgerDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 25MB, totalling 200 batches. For each batch we perform identical
// operations as in TestBolt_5GB_Experiment_10MB_Batch. This experiment uses the "Performance" Badger config.
func TestBadger_Performance_5GB_Experiment_25MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 200,
		BatchSizeBytes:            25000000,
		BatchSizeItems:            250,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}
	dir, err := os.MkdirTemp("", "badgerdb-performance-5gb-25mb")
	t.Logf("BadgerDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	opts := PerformanceBadgerOptions(dir)
	db := NewBadgerDatabase(opts, false)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	badgerCtx := db.GetContext([]byte{})
	GenericTest(db, badgerCtx, testConfig, t)
}

// TestBadger_Performance_5GB_Experiment_100MB_Batch is a BadgerDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 100MB, totalling 50 batches. For each batch we perform identical
// operations as in TestBolt_5GB_Experiment_10MB_Batch. This experiment uses the "Performance" Badger config.
func TestBadger_Performance_5GB_Experiment_100MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 50,
		BatchSizeBytes:            100000000,
		BatchSizeItems:            1000,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}
	dir, err := os.MkdirTemp("", "badgerdb-performance-5gb-100mb")
	t.Logf("BadgerDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	opts := PerformanceBadgerOptions(dir)
	db := NewBadgerDatabase(opts, false)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	badgerCtx := db.GetContext([]byte{})
	GenericTest(db, badgerCtx, testConfig, t)
}

// TestBadger_Default_WriteBatch_5GB_Experiment_10MB_Batch is a BadgerDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 10MB, totalling 500 batches. For each batch we perform identical operations
// as in TestBolt_5GB_Experiment_10MB_Batch. This experiment uses the "Default" Badger config, and Badger's WriteBatch.
func TestBadger_Default_WriteBatch_5GB_Experiment_10MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 500,
		BatchSizeBytes:            10000000,
		BatchSizeItems:            100,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}
	dir, err := os.MkdirTemp("", "badgerdb-default-writebatch-5gb-10mb")
	t.Logf("BadgerDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	opts := DefaultBadgerOptions(dir)
	db := NewBadgerDatabase(opts, true)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	badgerCtx := db.GetContext([]byte{})
	GenericTest(db, badgerCtx, testConfig, t)
}

// TestBadger_Default_WriteBatch_5GB_Experiment_25MB_Batch is a BadgerDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 25MB, totalling 200 batches. For each batch we perform identical operations
// as in TestBolt_5GB_Experiment_10MB_Batch. This experiment uses the "Default" Badger config, and Badger's WriteBatch.
func TestBadger_Default_WriteBatch_5GB_Experiment_25MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 200,
		BatchSizeBytes:            25000000,
		BatchSizeItems:            250,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}
	dir, err := os.MkdirTemp("", "badgerdb-default-writebatch-5gb-25mb")
	t.Logf("BadgerDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	opts := DefaultBadgerOptions(dir)
	db := NewBadgerDatabase(opts, true)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	badgerCtx := db.GetContext([]byte{})
	GenericTest(db, badgerCtx, testConfig, t)
}

// TestBadger_Default_WriteBatch_5GB_Experiment_100MB_Batch is a BadgerDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 100MB, totalling 50 batches. For each batch we perform identical operations
// as in TestBolt_5GB_Experiment_10MB_Batch. This experiment uses the "Default" Badger config, and Badger's WriteBatch.
func TestBadger_Default_WriteBatch_5GB_Experiment_100MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 50,
		BatchSizeBytes:            100000000,
		BatchSizeItems:            1000,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}
	dir, err := os.MkdirTemp("", "badgerdb-default-writebatch-5gb-100mb")
	t.Logf("BadgerDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	opts := DefaultBadgerOptions(dir)
	db := NewBadgerDatabase(opts, true)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	badgerCtx := db.GetContext([]byte{})
	GenericTest(db, badgerCtx, testConfig, t)
}

// TestBadger_Performance_WriteBatch_5GB_Experiment_10MB_Batch is a BadgerDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 10MB, totalling 500 batches. For each batch we perform identical operations
// as in TestBolt_5GB_Experiment_10MB_Batch. This experiment uses the "Performance" Badger config, and Badger's WriteBatch.
func TestBadger_Performance_WriteBatch_5GB_Experiment_10MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 500,
		BatchSizeBytes:            10000000,
		BatchSizeItems:            100,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}
	dir, err := os.MkdirTemp("", "badgerdb-performance-writebatch-5gb-10mb")
	t.Logf("BadgerDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	opts := PerformanceBadgerOptions(dir)
	db := NewBadgerDatabase(opts, true)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	badgerCtx := db.GetContext([]byte{})
	GenericTest(db, badgerCtx, testConfig, t)
}

// TestBadger_Performance_WriteBatch_5GB_Experiment_25MB_Batch is a BadgerDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 25MB, totalling 200 batches. For each batch we perform identical operations
// as in TestBolt_5GB_Experiment_10MB_Batch. This experiment uses the "Performance" Badger config, and Badger's WriteBatch.
func TestBadger_Performance_WriteBatch_5GB_Experiment_25MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 200,
		BatchSizeBytes:            25000000,
		BatchSizeItems:            250,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}
	dir, err := os.MkdirTemp("", "badgerdb-performance-writebatch-5gb-25mb")
	t.Logf("BadgerDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	opts := PerformanceBadgerOptions(dir)
	db := NewBadgerDatabase(opts, true)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	badgerCtx := db.GetContext([]byte{})
	GenericTest(db, badgerCtx, testConfig, t)
}

// TestBadger_Performance_WriteBatch_5GB_Experiment_100MB_Batch is a BadgerDB test in which we write 5GB of data to the database.
// In the test, we write data in batches of 100MB, totalling 50 batches. For each batch we perform identical operations
// as in TestBolt_5GB_Experiment_10MB_Batch. This experiment uses the "Performance" Badger config, and Badger's WriteBatch.
func TestBadger_Performance_WriteBatch_5GB_Experiment_100MB_Batch(t *testing.T) {
	require := require.New(t)

	testConfig := &TestConfig{
		ExperimentNumberOfBatches: 50,
		BatchSizeBytes:            100000000,
		BatchSizeItems:            1000,
		BatchItemsRemoved:         20,
		BatchItemsRetrieved:       20,
		BatchItemsIterated:        20,
	}
	dir, err := os.MkdirTemp("", "badgerdb-performance-writebatch-5gb-100mb")
	t.Logf("BadgerDB directory: %s\nIt should be automatically removed at the end of the test", dir)
	require.NoError(err)

	opts := PerformanceBadgerOptions(dir)
	db := NewBadgerDatabase(opts, true)
	require.NoError(db.Setup())
	defer db.Erase()
	defer db.Close()

	badgerCtx := db.GetContext([]byte{})
	GenericTest(db, badgerCtx, testConfig, t)
}

func GenericTest(db Database, ctx Context, config *TestConfig, t *testing.T) {
	timer := NewTimer()
	p := NewProfiler()
	for ii := 0; ii < config.ExperimentNumberOfBatches; ii++ {
		currentTime := timer.GetTotalElapsedTime("Experiment")
		GenericBatchTest(db, ctx, p, timer, config, t)
		deltaTime := timer.GetTotalElapsedTime("Experiment") - currentTime

		memUsageLog := p.PrintLastMeasurement()
		if ii%(config.ExperimentNumberOfBatches/50) == 0 {
			t.Logf("Batch %v/%v\n%v\nTotal time: %vs\t Delta time: %vs\n",
				ii, config.ExperimentNumberOfBatches, memUsageLog,
				timer.GetTotalElapsedTime("Experiment"), deltaTime)
		}
	}
	t.Logf("Timer results:\n%v", timer.Print("Experiment"))
	t.Logf("Profiler results:\n%v", p.PrintStats())
}

func GenericBatchTest(db Database, ctx Context, p *Profiler, timer *Timer, config *TestConfig, t *testing.T) {
	require := require.New(t)

	timer.Start("Experiment")
	defer timer.End("Experiment")

	// Generate random data and store it in the DB
	kvMap := WriteBatchToDb(db, ctx, config, t)
	p.Measure()

	// Choose keys to remove and keys to retrieve
	removedKeys := []Key{}
	retrievedKeys := []Key{}
	for key := range kvMap {
		if len(removedKeys) < config.BatchItemsRemoved {
			removedKeys = append(removedKeys, key)
		} else if len(retrievedKeys) < config.BatchItemsRetrieved {
			retrievedKeys = append(retrievedKeys, key)
		} else {
			break
		}
	}

	// Delete keys from DB and confirm they are deleted.
	DeleteFromDB(db, removedKeys, ctx, t)
	deletedValues := GetFromDb(db, removedKeys, ctx, t)
	for _, val := range deletedValues {
		require.Nil(val)
	}
	p.Measure()

	// Retrieve keys from DB and confirm they match the original values.
	retrievedValues := GetFromDb(db, retrievedKeys, ctx, t)
	for i, val := range retrievedValues {
		require.Equal(kvMap[retrievedKeys[i]], val)
	}
	p.Measure()

	// Iterate over a couple values from the DB and confirm they match the original values.
	IterateOverDb(db, ctx, config, t)
	p.Measure()
}

func WriteBatchToDb(db Database, ctx Context, config *TestConfig, t *testing.T) (_kv map[Key][]byte) {
	require := require.New(t)
	kvMap := make(map[Key][]byte)
	valueLenght := int32(config.BatchSizeBytes / config.BatchSizeItems)

	for ii := 0; ii < config.BatchSizeItems; ii++ {
		randomKey, err := RandomBytes(32)
		require.NoError(err)
		keyHash := NewKey(randomKey)
		kvMap[keyHash], err = RandomBytes(valueLenght)
		require.NoError(err)
	}

	require.NoError(db.Update(ctx, func(tx Transaction, ctx Context) error {
		for key, val := range kvMap {
			if err := tx.Set(key.Bytes(), val, ctx); err != nil {
				return err
			}
		}
		return nil
	}))
	return kvMap
}

func DeleteFromDB(db Database, keys []Key, ctx Context, t *testing.T) {
	require := require.New(t)
	require.NoError(db.Update(ctx, func(tx Transaction, ctx Context) error {
		for _, key := range keys {
			if err := tx.Delete(key.Bytes(), ctx); err != nil {
				return err
			}
		}
		return nil
	}))
}

func GetFromDb(db Database, keys []Key, ctx Context, t *testing.T) [][]byte {
	require := require.New(t)
	var values [][]byte
	require.NoError(db.View(ctx, func(tx Transaction, ctx Context) error {
		for _, key := range keys {
			val, err := tx.Get(key.Bytes(), ctx)
			if err != nil {
				val = nil
			}
			values = append(values, val)
		}
		return nil
	}))
	return values
}

func IterateOverDb(db Database, ctx Context, config *TestConfig, t *testing.T) {
	require := require.New(t)
	iterationCount := 0
	require.NoError(db.Update(ctx, func(tx Transaction, ctx Context) error {
		it, err := tx.GetIterator(ctx)
		require.NoError(err)
		defer it.Close()
		for it.Next() {
			k := it.Key()
			v, err := it.Value()
			require.NoError(err)
			require.NotNil(k)
			require.NotNil(v)
			iterationCount++
			if iterationCount >= config.BatchItemsIterated {
				break
			}
		}
		return nil
	}))
}
