package sarray

import (
	"errors"
	"github.com/snzysnk/gs/v2/internal/srwmutex"
	"sort"
)

var _ IStrArray = (*strArray)(nil)

type IStrArray interface {
	Get(i int) (value string, found bool)
	Set(i int, value string) error
	Len() int
	SortFunc(func(i, k string) bool)
	ToStrSlice() []string
	Insert(i int, v string) error
	Remove(i int) (value string, found bool)
	RemoveWithOutLock(i int) (value string, found bool)
	FindIndexWithOutLock(value string) (index int)
	Clear() IStrArray
	Unique() IStrArray
}

type strArray struct {
	mu   srwmutex.IRWMutex
	data []string
}

func (s strArray) Unique() IStrArray {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.data) == 0 {
		return s
	}
	var (
		uniqueArr = make([]string, 0, len(s.data))
		uniqueSet = make(map[string]struct{})
	)
	for _, v := range s.data {
		if _, ok := uniqueSet[v]; ok {
			continue
		}
		uniqueSet[v] = struct{}{}
		uniqueArr = append(uniqueArr, v)
	}
	s.data = uniqueArr
	return s
}

func (s strArray) Clear() IStrArray {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.data) > 0 {
		s.data = make([]string, 0)
	}

	return s
}

func (s strArray) FindIndexWithOutLock(value string) (index int) {
	for i, v := range s.data {
		if v == value {
			return i
		}
	}
	return -1
}

func (s strArray) Remove(i int) (value string, found bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.RemoveWithOutLock(i)
}

func (s strArray) RemoveWithOutLock(i int) (value string, found bool) {
	if i < 0 || i >= len(s.data) {
		return "", false
	}
	value = s.data[i]
	s.data = append(s.data[:i], s.data[i+1:]...)
	return value, true
}

func (s strArray) Insert(i int, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i < 0 || i > len(s.data) {
		return errors.New("index out of range")
	}
	s.data = append(s.data[:i], v, s.data[i])
	return nil
}

func (s strArray) ToStrSlice() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data := make([]string, len(s.data))
	copy(data, s.data)
	return data
}

func (s strArray) SortFunc(f func(i, j string) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sort.Slice(s.data, func(i, j int) bool {
		return f(s.data[i], s.data[j])
	})
}

func (s strArray) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.data)
}

func (s strArray) Get(i int) (value string, found bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if i < 0 || i >= len(s.data) {
		return "", false
	}
	return s.data[i], true
}

func (s strArray) Set(i int, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i < 0 || i >= len(s.data) {
		return errors.New("index out of range or index letter zero")
	}
	s.data[i] = value
	return nil
}

func (s strArray) Clone() IStrArray {
	s.mu.RLock()
	defer s.mu.RUnlock()
	newStrArr := make([]string, len(s.data))
	_, isSafe := s.mu.(*srwmutex.SafeRWMutex)
	return NewStrArrayFrom(newStrArr, isSafe)
}

func NewStrArraySize(size int, cap int, safe bool) IStrArray {
	return &strArray{
		mu:   srwmutex.New(safe),
		data: make([]string, size, cap),
	}
}

func NewStrArray(safe bool) IStrArray {
	return NewStrArraySize(0, 0, safe)
}

func NewStrArrayFrom(array []string, safe bool) IStrArray {
	return &strArray{
		mu:   srwmutex.New(safe),
		data: array,
	}
}

func NewStrArrayFromCopy(array []string, safe bool) IStrArray {
	newArray := make([]string, len(array))
	copy(newArray, array)
	return &strArray{
		mu:   srwmutex.New(safe),
		data: newArray,
	}
}
