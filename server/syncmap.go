package main

import "sync"

type SyncMap[K, V any] struct {
	m sync.Map
}

func NewSyncMap[K, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{}
}

func (s *SyncMap[K, V]) Load(k K) (v V, ok bool) {
	var val any
	if val, ok = s.m.Load(k); ok {
		v = val.(V)
	}
	return
}

// Returns true if the value was stored
func (s *SyncMap[K, V]) LoadOrStore(k K, v V) (V, bool) {
	val, ok := s.m.LoadOrStore(k, v)
	return val.(V), ok
}

func (s *SyncMap[K, V]) Store(k K, v V) {
	s.m.Store(k, v)
}

func (s *SyncMap[K, V]) LoadAndDelete(k K) (v V, ok bool) {
	var val any
	if val, ok = s.LoadAndDelete(k); ok {
		v = val.(V)
	}
	return
}

func (s *SyncMap[K, V]) Delete(k K) {
	s.m.Delete(k)
}

func (s *SyncMap[K, V]) Range(f func(K, V) bool) {
	s.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}
