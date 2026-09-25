package syncmap

import (
	"sync"
	"testing"
)

// sync.Map выигрывает: все читают один и тот же существующий ключ.
func BenchmarkSyncMapRead(b *testing.B) {
	var m sync.Map
	m.Store("k", 1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Load("k")
		}
	})
}

func BenchmarkRWMutexMapRead(b *testing.B) {
	var mu sync.RWMutex
	m := map[string]int{"k": 1}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.RLock()
			_ = m["k"]
			mu.RUnlock()
		}
	})
}

// sync.Map проигрывает: без конкуренции платить за её бухгалтерию
// (атомарный указатель, entry в куче) уже нечем — окупать её нечем,
// когда никто ни за что не борется. Один вызывающий, никаких RunParallel.
func BenchmarkSyncMapInsertSerial(b *testing.B) {
	var m sync.Map
	i := 0
	for b.Loop() {
		m.Store(i, true)
		i++
	}
}

func BenchmarkMutexMapInsertSerial(b *testing.B) {
	var mu sync.Mutex
	m := make(map[int]bool)
	i := 0
	for b.Loop() {
		mu.Lock()
		m[i] = true
		mu.Unlock()
		i++
	}
}
