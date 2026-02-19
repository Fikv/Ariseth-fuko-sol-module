package scanner

import "sync"

type Store struct {
	mu     sync.RWMutex
	keys   map[string]struct{}
	events []Event
	subs   map[chan Event]struct{}
}

func NewStore() *Store {
	return &Store{
		keys: make(map[string]struct{}),
		subs: make(map[chan Event]struct{}),
	}
}

func (s *Store) Add(e Event) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := e.Key()
	if _, exists := s.keys[key]; exists {
		return false
	}

	s.keys[key] = struct{}{}
	s.events = append([]Event{e}, s.events...)

	for ch := range s.subs {
		select {
		case ch <- e:
		default:
		}
	}

	return true
}

func (s *Store) List(limit int) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 || limit > len(s.events) {
		limit = len(s.events)
	}

	out := make([]Event, limit)
	copy(out, s.events[:limit])
	return out
}

func (s *Store) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 32)

	s.mu.Lock()
	s.subs[ch] = struct{}{}
	s.mu.Unlock()

	unsubscribe := func() {
		s.mu.Lock()
		if _, ok := s.subs[ch]; ok {
			delete(s.subs, ch)
			close(ch)
		}
		s.mu.Unlock()
	}

	return ch, unsubscribe
}
