package syncmap

import (
	"sync"
	"testing"
)

// Намеренно красный тест: Load, потом Store — каждый вызов сам по себе
// безопасен, но вместе это не атомарный инкремент, гонка теряет
// обновления.
func TestSyncMapLosesUpdates(t *testing.T) {
	var m sync.Map
	m.Store("count", 0)

	var wg sync.WaitGroup
	for range 1000 {
		wg.Go(func() {
			v, _ := m.Load("count")
			m.Store("count", v.(int)+1)
		})
	}
	wg.Wait()

	got, _ := m.Load("count")
	if got.(int) != 1000 {
		t.Errorf("получили %d, ожидали 1000 — sync.Map потеряла обновления", got)
	}
}
