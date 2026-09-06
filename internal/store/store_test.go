package store

import (
	"sync"
	"testing"
)

func TestEncodeBase62(t *testing.T) {
	tests := []struct {
		name string
		id   int
		want string
	}{
		{"zero", 0, "0"},
		{"one", 1, "1"},
		{"two", 2, "2"},
		{"ten", 10, "a"},
		{"sixty two", 62, "10"},
		{"125", 125, "21"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeBase62(tt.id)
			if got != tt.want {
				t.Errorf("encodeBase62(%d) got = %v, want %v", tt.id, got, tt.want)
			}
		})
	}
}
func TestStore_SaveAndGet(t *testing.T) {
	s := New()
	url := "http://example.com"
	code := s.Save(url)

	got, ok := s.Get(code)
	if !ok {
		t.Fatal("code should exist after save")
	}
	if got != url {
		t.Errorf("got %v, want %v", got, url)
	}
	_, ok = s.Get("nonexistent")
	if ok {
		t.Fatal("code should not exist after save")
	}
}
func TestStore_Concurrent(t *testing.T) {
	s := New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			code := s.Save("http://example.com")
			s.Get(code)
		}()
	}
	wg.Wait()
}
