package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

// CAS на уровне процессора против мьютекса: обе горутины дерутся за одну
// и ту же переменную-счётчик, разница только в способе синхронизации.
func BenchmarkMutexCounter(b *testing.B) {
	var mu sync.Mutex
	var counter int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}

func BenchmarkAtomicCounter(b *testing.B) {
	var counter atomic.Int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Add(1)
		}
	})
}
