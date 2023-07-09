package main

import (
	"crypto/rand"
	"fmt"
	"github.com/pkg/errors"
	"runtime"
)

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
// Taken from: https://golangcode.com/print-the-current-memory-usage/
func PrintMemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	log := fmt.Sprintf("Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v",
		bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
	return log
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// RandomBytes returns a []byte with random values.
func RandomBytes(numBytes int32) ([]byte, error) {
	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "Problem reading random bytes")
	}
	return randomBytes, nil
}
