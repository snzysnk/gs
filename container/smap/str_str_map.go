package smap

import (
	"github.com/snzysnk/gs/v2/internal/srwmutex"
)

var _ IStrStrMap = (*strStrMap)(nil)

type IStrStrMap interface {
	Foreach(f func(k, v string))
	ToNewMap() map[string]string
	Clone() IStrStrMap
	Set(k string, v string)
	Get(k string) string
	Search(k string) (string, bool)
	GetOrSet(k string, v string) string
	GetOrSetWithLock(k string, v string) string
	LockFunc(f func(map[string]string))
	RLockFunc(f func(map[string]string))
	Remove(k string) string
}

type strStrMap struct {
	mu   srwmutex.IRWMutex
	data map[string]string
}

func (s *strStrMap) Remove(k string) (value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var ok bool
	if value, ok = s.data[k]; ok {
		delete(s.data, k)
	}
	return value
}

func NewStrStrMap(safe bool) IStrStrMap {
	return &strStrMap{
		mu:   srwmutex.New(safe),
		data: make(map[string]string),
	}
}

func NewStrStrFormMap(data map[string]string, safe bool) IStrStrMap {
	return &strStrMap{
		mu:   srwmutex.New(safe),
		data: data,
	}
}

func (s *strStrMap) Foreach(f func(k, v string)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.data {
		f(k, v)
	}
}

func (s *strStrMap) ToNewMap() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	newMap := make(map[string]string, len(s.data))
	for k, v := range s.data {
		newMap[k] = v
	}
	return newMap
}

func (s *strStrMap) Clone() IStrStrMap {
	_, isSafe := s.mu.(*srwmutex.SafeRWMutex)
	return NewStrStrFormMap(s.ToNewMap(), isSafe)
}

func (s *strStrMap) Set(k string, v string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[k] = v
}

func (s *strStrMap) LockFunc(f func(map[string]string)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f(s.data)
}

func (s *strStrMap) RLockFunc(f func(map[string]string)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f(s.data)
}

func (s *strStrMap) GetOrSetWithLock(k string, v string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if value, ok := s.data[k]; ok {
		return value
	}
	s.data[k] = v
	return s.data[k]
}

func (s *strStrMap) GetOrSet(k string, v string) string {
	if value, ok := s.Search(k); ok {
		return value
	}
	return s.GetOrSetWithLock(k, v)
}

func (s *strStrMap) Search(k string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[k]
	return value, ok
}

func (s *strStrMap) Get(k string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[k]
}
