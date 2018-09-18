package slice_test

import (
	"strconv"
	"sync"
	"testing"
)

// go test -bench="."
type Something struct {
	roomId   int
	roomName string
}

func BenchmarkDefaultSlice(b *testing.B) {
	b.ReportAllocs()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for i := 0; i < 120; i++ {
				output := make([]Something, 0)
				output = append(output, Something{
					roomId:   i,
					roomName: strconv.Itoa(i),
				})
			}
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}

func BenchmarkPreAllocSlice(b *testing.B) {
	b.ReportAllocs()
	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			output := make([]Something, 0, 120)
			for i := 0; i < 120; i++ {
				output = append(output, Something{
					roomId:   i,
					roomName: strconv.Itoa(i),
				})
			}
			wg.Done()
		}(&wg)
	}
	wg.Wait()

}

func BenchmarkSyncPoolSlice(b *testing.B) {
	b.ReportAllocs()
	var wg sync.WaitGroup
	var SomethingPool = sync.Pool{
		New: func() interface{} {
			b := make([]Something, 120)
			return &b
		},
	}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			obj := SomethingPool.Get().(*[]Something)
			for i := 0; i < 120; i++ {
				some := *obj
				some[i].roomId = i
				some[i].roomName = strconv.Itoa(i)
			}
			SomethingPool.Put(obj)
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}
