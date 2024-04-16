package srwmutex

import "sync"

var _ IRWMutex = (*SafeRWMutex)(nil)

type SafeRWMutex struct {
	sync.RWMutex
}
