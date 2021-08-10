package stream

import (
	"fmt"
	"testing"
	"time"
)

func TestBaseStream(t *testing.T) {
	given := dataGenerator()
	start := time.Now()
	first := New(given).Filter(func(i interface{}) bool {
		return i.(int)%2 == 0
	}).FlatMap(func(i interface{}) []interface{} {
		return []interface{}{i, i.(int) + 2}
	}).Skip(3).Collect()
	fmt.Println(first)
	fmt.Println("\nstream", time.Now().Sub(start))

	//fmt.Println(stream.Min(func(a interface{}, b interface{}) int {
	//	return a.(int) - b.(int)
	//}))
}

func TestNonStream(t *testing.T) {
	given := dataGenerator()
	start := time.Now()
	//for i := 0; i < b.N; i++ {
	set := make(map[interface{}]byte)
	skip := 0
	for idx := range given {
		if given[idx]%2 == 0 {
			set[given[idx]] = 0x01
			set[given[idx]+2] = 0x01
		}
	}

	for _, _ = range set {
		if skip < 3 {
			skip++
			continue
		}
		fmt.Print(" ")
	}

	//for key := range set {
	//	fmt.Print(key)
	//}
	fmt.Println("\nnon stream", time.Now().Sub(start))
	//}
}

func dataGenerator() []int {
	res := make([]int, 0)
	for idx := 0; idx < 200; idx++ {
		res = append(res, idx+1)
	}
	return res
}

func BenchmarkBaseStream(b *testing.B) {
	given := dataGenerator()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := New(given)
		s.Parallel().Filter(func(i interface{}) bool {
			return i.(int)%2 == 0
		}).FlatMap(func(i interface{}) []interface{} {
			return []interface{}{i, i.(int) + 2}
		}).ForEach(func(i interface{}) {

		})
	}
}

func BenchmarkNonStream(b *testing.B) {
	given := dataGenerator()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set := make(map[interface{}]byte)
		for idx := range given {
			if given[idx]%2 == 0 {
				set[given[idx]] = 0x01
				set[given[idx]+2] = 0x01
			}
		}

		for key := range set {
			_ = key
		}
	}
}
