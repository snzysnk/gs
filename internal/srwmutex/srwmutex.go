package srwmutex

type IRWMutex interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

func New(safe bool) IRWMutex {
	if safe {
		return new(SafeRWMutex)
	}
	return new(unsafeRWMutex)
}
