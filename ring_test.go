package ring

import (
	"github.com/valyala/fastrand"
	"sync"
	"sync/atomic"
	"testing"
)

const (
	elements = 100000000
	fpr      = 0.0001
)

func uint32bin(v uint32) []byte {
	return []byte{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)}
}

func BenchmarkAddConcurrent(b *testing.B) {
	f, _ := Init(elements, fpr)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f.Add(uint32bin(fastrand.Uint32()))
		}
	})
}

func BenchmarkTestConcurrent(b *testing.B) {
	f, _ := Init(elements, fpr)

	for i := 0; i < 1000000; i++ {
		f.Add(uint32bin(uint32(i)))
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f.Test(uint32bin(fastrand.Uint32() % 1000000))
		}
	})
}

func TestConcurrentErrors(t *testing.T) {
	f, _ := Init(elements, fpr)
	mutex := sync.Mutex{}
	data := make(map[uint32]bool)
	addCh := make(chan uint32, 1000)
	fnrTestCh := make(chan uint32, 1000)
	fprTestCh := make(chan uint32, 1000)

	wg1 := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}

	for i := 0; i < 4; i++ {
		wg1.Add(1)
		go func() {
			for v := range addCh {
				mutex.Lock()
				f.Add(uint32bin(v))
				data[v] = true
				mutex.Unlock()
			}
			wg1.Done()
		}()

		wg2.Add(1)
		go func() {
			for i := 0; i < elements/4; i++ {
				addCh <- fastrand.Uint32()
			}
			wg2.Done()
		}()
	}

	wg2.Wait()
	close(addCh)
	wg1.Wait()

	var fnrErr int64 = 0
	var fprErr int64 = 0

	for i := 0; i < 4; i++ {
		wg1.Add(1)
		go func() {
			for v := range fnrTestCh {
				if !f.Test(uint32bin(v)) {
					atomic.AddInt64(&fnrErr, 1)
				}
			}
			wg1.Done()
		}()

		wg1.Add(1)
		go func() {
			for v := range fprTestCh {
				_, ok := data[v]
				if !ok && f.Test(uint32bin(v)) {
					atomic.AddInt64(&fprErr, 1)
				}
			}
			wg1.Done()
		}()

		wg2.Add(1)
		go func() {
			for i := 0; i < elements/4; i++ {
				fprTestCh <- fastrand.Uint32()
			}
			wg2.Done()
		}()
	}

	for v := range data {
		fnrTestCh <- v
	}

	close(fnrTestCh)
	wg2.Wait()
	close(fprTestCh)
	wg1.Wait()

	t.Logf("FNR: %f, errors: %d", float64(fnrErr) / float64(elements) * 100, fnrErr)
	t.Logf("FPR: %f, errors: %d", float64(fprErr) / float64(elements) * 100, fprErr)
}
