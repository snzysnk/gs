package smap

import (
	"github.com/snzysnk/gs/v2/internal/srwmutex"
)

var _ IStrAnyMap = (*strAnyMap)(nil)

type IStrAnyMap interface {
	Foreach(f func(k, v interface{}))
	ToNewMap() map[string]interface{}
	Clone() IStrAnyMap
	Set(k string, v interface{})
	Get(k string) interface{}
	Search(k string) (interface{}, bool)
	GetOrSet(k string, v interface{}) interface{}
	GetOrSetWithLock(k string, v interface{}) interface{}
	LockFunc(f func(map[string]interface{}))
	RLockFunc(f func(map[string]interface{}))
	Remove(k string) interface{}
}

type strAnyMap struct {
	mu   srwmutex.IRWMutex
	data map[string]interface{}
}

func (s strAnyMap) Remove(k string) (value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var ok bool
	if value, ok = s.data[k]; ok {
		delete(s.data, k)
	}
	return value
}

func New(safe bool) IStrAnyMap {
	return &strAnyMap{
		mu:   srwmutex.New(safe),
		data: make(map[string]interface{}),
	}
}

func NewFormMap(data map[string]interface{}, safe bool) IStrAnyMap {
	return &strAnyMap{
		mu:   srwmutex.New(safe),
		data: data,
	}
}

func (s strAnyMap) Foreach(f func(k interface{}, v interface{})) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.data {
		f(k, v)
	}
}

func (s strAnyMap) ToNewMap() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	newMap := make(map[string]interface{}, len(s.data))
	for k, v := range s.data {
		newMap[k] = v
	}
	return newMap
}

func (s strAnyMap) Clone() IStrAnyMap {
	_, isSafe := s.mu.(*srwmutex.SafeRWMutex)
	return NewFormMap(s.ToNewMap(), isSafe)
}

func (s strAnyMap) Set(k string, v interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if f, ok := v.(func() interface{}); ok {
		s.data[k] = f()
	} else {
		s.data[k] = v
	}
}

func (s strAnyMap) LockFunc(f func(map[string]interface{})) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f(s.data)
}

func (s strAnyMap) RLockFunc(f func(map[string]interface{})) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f(s.data)
}

func (s strAnyMap) GetOrSetWithLock(k string, v interface{}) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	if value, ok := s.data[k]; ok {
		return value
	}
	if f, ok := v.(func() interface{}); ok {
		s.data[k] = f()
	} else {
		s.data[k] = v
	}
	return s.data[k]
}

func (s strAnyMap) GetOrSet(k string, v interface{}) interface{} {
	if value, ok := s.Search(k); ok {
		return value
	}
	return s.GetOrSetWithLock(k, v)
}

func (s strAnyMap) Search(k string) (interface{}, bool) {
	s.mu.RLock()
	s.mu.RUnlock()
	value, ok := s.data[k]
	return value, ok
}

func (s strAnyMap) Get(k string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[k]
}
