package stats

import (
	"sync"
	"testing"
)

func TestTop_Empty(t *testing.T) {
	s := NewStore()
	if _, _, ok := s.Top(); ok {
		t.Fatal("expected ok=false on empty store")
	}
}

func TestRecordAndTop(t *testing.T) {
	s := NewStore()
	a := Key{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"}
	b := Key{Int1: 2, Int2: 7, Limit: 50, Str1: "foo", Str2: "bar"}

	s.Record(a)
	s.Record(a)
	s.Record(b)

	key, hits, ok := s.Top()
	if !ok {
		t.Fatal("expected ok=true")
	}
	if key != a || hits != 2 {
		t.Errorf("got key=%+v hits=%d, want key=%+v hits=2", key, hits, a)
	}
}

func TestRecord_ConcurrentAccess(t *testing.T) {
	s := NewStore()
	k := Key{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"}

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			s.Record(k)
		}()
	}
	wg.Wait()

	_, hits, ok := s.Top()
	if !ok || hits != goroutines {
		t.Errorf("got hits=%d ok=%v, want hits=%d ok=true", hits, ok, goroutines)
	}
}
