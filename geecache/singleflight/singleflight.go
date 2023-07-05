package singleflight

import "sync"

// call is an in-flight or completed Do call
type packet struct { //正在进行或已结束的请求
	wg  sync.WaitGroup //避免重入
	val interface{}
	err error
}

// Flight represents a class of work and forms a namespace in which
// units of work can be executed with duplicate suppression.
type Flight struct {
	mu     sync.Mutex         // protects m
	flight map[string]*packet // lazily initialized
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
func (f *Flight) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	f.mu.Lock()
	if f.flight == nil {
		f.flight = make(map[string]*packet)
	}
	if c, ok := f.flight[key]; ok {
		f.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(packet)
	c.wg.Add(1)
	f.flight[key] = c
	f.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	f.mu.Lock()
	delete(f.flight, key)
	f.mu.Unlock()

	return c.val, c.err
}
