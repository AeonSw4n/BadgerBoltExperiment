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

type Profiler struct {
	AllocMeasurements      []uint64
	TotalAllocMeasurements []uint64
	SysMeasurements        []uint64
	NumGCMeasurements      []uint64
}

func NewProfiler() *Profiler {
	return &Profiler{
		AllocMeasurements:      make([]uint64, 0),
		TotalAllocMeasurements: make([]uint64, 0),
		SysMeasurements:        make([]uint64, 0),
		NumGCMeasurements:      make([]uint64, 0),
	}
}

func (p *Profiler) Measure() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	p.AllocMeasurements = append(p.AllocMeasurements, bToMb(m.Alloc))
	p.TotalAllocMeasurements = append(p.TotalAllocMeasurements, bToMb(m.TotalAlloc))
	p.SysMeasurements = append(p.SysMeasurements, bToMb(m.Sys))
	p.NumGCMeasurements = append(p.NumGCMeasurements, uint64(m.NumGC))
}

func computeMean(values []uint64) uint64 {
	var sum uint64
	for _, value := range values {
		sum += value
	}
	return sum / uint64(len(values))
}

func computeMax(values []uint64) uint64 {
	var max uint64
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
}

func (p *Profiler) Print() string {
	log := fmt.Sprintf("MEAN STATS \t|\t Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v\n",
		computeMean(p.AllocMeasurements), computeMean(p.TotalAllocMeasurements),
		computeMean(p.SysMeasurements), computeMean(p.NumGCMeasurements))
	log += fmt.Sprintf("MAX STATS \t|\t Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v\n",
		computeMax(p.AllocMeasurements), computeMax(p.TotalAllocMeasurements),
		computeMax(p.SysMeasurements), computeMax(p.NumGCMeasurements))
	return log
}
