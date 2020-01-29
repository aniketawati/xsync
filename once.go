package xsync

import (
	"sync"
	"sync/atomic"
)

// This implementation of Once keeps calling the supplied function
// until it has been successfully executed exactly once.
// So this guarantees exactly one successful execution of supplied function.
//

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 0 {
		return o.doSlow(f)
	}
	return nil
}

func (o *Once) doSlow(f func() error) error {
	o.m.Lock()
	defer o.m.Unlock()

	if o.done == 0 {
		if err := f(); err == nil {
			atomic.StoreUint32(&o.done, 1)
		} else {
			return err
		}
	}
	return nil
}