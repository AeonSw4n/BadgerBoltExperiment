package main

import (
	"fmt"
	"sync"
	"time"
)

type Timer struct {
	totalElapsedTimes map[string]float64
	lastTimes         map[string]time.Time
	mut               sync.RWMutex
}

func NewTimer() *Timer {
	return &Timer{
		totalElapsedTimes: make(map[string]float64),
		lastTimes:         make(map[string]time.Time),
	}
}

func (t *Timer) Start(eventName string) {
	t.mut.Lock()
	defer t.mut.Unlock()
	if _, exists := t.lastTimes[eventName]; !exists {
		t.totalElapsedTimes[eventName] = 0.0
	}
	t.lastTimes[eventName] = time.Now()
}

func (t *Timer) End(eventName string) {
	t.mut.Lock()
	defer t.mut.Unlock()
	if _, exists := t.totalElapsedTimes[eventName]; !exists {
		return
	}
	t.totalElapsedTimes[eventName] += time.Since(t.lastTimes[eventName]).Seconds()
}

func (t *Timer) GetTotalElapsedTime(eventName string) float64 {
	t.mut.RLock()
	defer t.mut.RUnlock()
	if _, exists := t.totalElapsedTimes[eventName]; exists {
		return t.totalElapsedTimes[eventName]
	}
	return 0.0
}

func (t *Timer) Print(eventName string) string {
	t.mut.RLock()
	defer t.mut.RUnlock()

	totalElapsedTime := float64(0.0)
	if _, exists := t.lastTimes[eventName]; exists {
		totalElapsedTime = t.totalElapsedTimes[eventName]
	}
	return fmt.Sprintf("Timer.End: event (%s) total elapsed time (%v)", eventName, totalElapsedTime)
}
