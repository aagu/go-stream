package stream

import (
	"fmt"
	"testing"
)

func TestBaseStream(t *testing.T) {
	stream := Of(1, 2, 3, 4, 5, 6, 7).FlatMap(func(i interface{}) []interface{} {
		return []interface{}{i, i.(int) + 2}
	}).Distinct()

	stream.ForEach(func(i interface{}) {
		fmt.Println(i)
	})

	//fmt.Println(stream.Min(func(a interface{}, b interface{}) int {
	//	return a.(int) - b.(int)
	//}))
}
