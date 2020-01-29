package xsync_test

import (
	"errors"
	"testing"
	. "github.com/aniketawati/xsync"
)

// should support all the tests in sync/once_test.go

type one int

func (o *one) Increment() {
	*o++
}

func run(t *testing.T, once *Once, o *one, c chan bool) {
	once.Do(func() error { o.Increment(); return nil })
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestOnce(t *testing.T) {
	o := new(one)
	once := new(Once)
	c := make(chan bool)
	const N = 10
	for i := 0; i < N; i++ {
		go run(t, once, o, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *o != 1 {
		t.Errorf("once failed outside run: %d is not 1", *o)
	}
}


/*
This is a major behavior change from standard library, as we want the function to be called
again till we have an error free execution. So even in cases where
called function panics, next call to once.Do with that function will
execute the function again.
 */

func TestOncePanic(t *testing.T) {
	var once Once
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("Once.Do did not panic")
			}
		}()
		once.Do(func() error {
			panic("failed")
		})
	}()
	called := false
	once.Do(func() error {
		called = true
		return nil
	})

	if !called {
		t.Fatalf("once.Do wasn't called again after panic")
	}
}

func TestOnceError(t *testing.T)  {
	var once Once
	once.Do(func() error {
		return errors.New("some random error")
	})

	called := false
	once.Do(func() error {
		called = true
		return nil
	})
	if !called {
		t.Fatalf("Once.Do wasn't called again in case of error")
	}
	once.Do(func() error {
		t.Fatalf("Once.Do called again!")
		return nil
	})
}

func BenchmarkOnce(b *testing.B) {
	var once Once
	f := func() error { return nil }
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			once.Do(f)
		}
	})
}
