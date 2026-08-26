package stats

import "sync"

type Key struct {
	Int1  int
	Int2  int
	Limit int
	Str1  string
	Str2  string
}

type Store struct {
	mu     sync.RWMutex
	counts map[Key]int
}

func NewStore() *Store {
	return &Store{counts: make(map[Key]int)}
}

func (s *Store) Record(k Key) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[k]++
}

// Top returns the request with the highest hit count and its count.
// On a tie, the first key encountered during map iteration wins; since Go randomizes map iteration
// order, the winner among tied requests is not deterministic across calls.
func (s *Store) Top() (key Key, hits int, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for k, c := range s.counts {
		if !ok || c > hits {
			key, hits, ok = k, c, true
		}
	}
	return key, hits, ok
}
