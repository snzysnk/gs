package sarray

import "github.com/snzysnk/gs/v2/internal/srwmutex"

var _ IStrArray = (*strArray)(nil)

type IStrArray interface {
	Clone() IStrArray
}

type strArray struct {
	mu   srwmutex.IRWMutex
	data []string
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
