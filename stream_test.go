package stream

import (
	"fmt"
	"testing"
)

func TestBaseStream(t *testing.T) {
	data := make([]interface{}, 7)
	for idx := range data {
		data[idx] = idx + 1
	}

	stream := New(data).FlatMap(func(i interface{}) []interface{} {
		return []interface{}{i, i.(int) + 2}
	}).Distinct()

	stream.ForEach(func(i interface{}) {
		fmt.Println(i)
	})

	//fmt.Println(stream.Min(func(a interface{}, b interface{}) int {
	//	return a.(int) - b.(int)
	//}))
}
