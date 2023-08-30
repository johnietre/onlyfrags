package main

import "sync"

type Mutex[T any] struct {
	valPtr *T
	mtx    sync.Mutex
}

func NewMutex[T any](val T) *Mutex[T] {
	return &Mutex[T]{valPtr: &val}
}

func (m *Mutex[T]) Lock() *T {
	m.mtx.Lock()
	return m.valPtr
}

func (m *Mutex[T]) Unlock() {
	m.mtx.Unlock()
}

func (m *Mutex[T]) Apply(f func(*T)) {
	m.mtx.Lock()
	f(m.valPtr)
	m.mtx.Unlock()
}

type RWMutex[T any] struct {
	valPtr *T
	mtx    sync.RWMutex
}

func NewRWMutex[T any](val T) *RWMutex[T] {
	return &RWMutex[T]{valPtr: &val}
}

func (m *RWMutex[T]) Lock() *T {
	m.mtx.Lock()
	return m.valPtr
}

func (m *RWMutex[T]) Unlock() {
	m.mtx.Unlock()
}

func (m *RWMutex[T]) RLock() *T {
	m.mtx.RLock()
	return m.valPtr
}

func (m *RWMutex[T]) RUnlock() {
	m.mtx.RUnlock()
}

func (m *RWMutex[T]) Apply(f func(*T)) {
	m.mtx.Lock()
	f(m.valPtr)
	m.mtx.Unlock()
}

func (m *RWMutex[T]) ApplyRead(f func(*T)) {
	m.mtx.RLock()
	f(m.valPtr)
	m.mtx.RUnlock()
}
