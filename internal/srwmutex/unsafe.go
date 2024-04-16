package srwmutex

var _ IRWMutex = (*unsafeRWMutex)(nil)

type unsafeRWMutex struct {
}

func (u unsafeRWMutex) Lock() {

}

func (u unsafeRWMutex) Unlock() {

}

func (u unsafeRWMutex) RLock() {

}

func (u unsafeRWMutex) RUnlock() {

}
