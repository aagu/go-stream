# go-stream

go-stream provides Java Stream API like stream operations in golang

## usage

```go
package main

import (
	"github.com/aagu/go-stream"
)

func main()  {
    ints := []int{1, 2, 3, 4, 5, 6, 7}
    stream.New(ints).Filter(func(i interface{}) bool {
		return i.(int) % 2 == 0
	}).FlatMap(func(i interface{}) []interface{} {
		return []interface{}{i, i.(int) + 2}
	}).Distinct().ForEach(func(i interface{}) {
		fmt.Println(i)
	})
}
```

current supports:

|function|describe|
| - | - |
| Filter | Filter uses a FilterFunc to filter out data |
| Map | Map transform data to another shape uses given MapFunc |
| FlatMap | transform datum to multiple data uses given FlatMapFunc |
| Distinct | pass only the different data to next stage use a golang built in map |
| Skip | not pass the first n elements it received to next stage |
| Limit | guarantee that no more than n elements pass to next stage |
| Sort | use a given ComparatorFunc to sort data |
| Group | use a given GroupFunc to split data into multiple groups |
| ForEach | call the given ForEachFunc to every element it received |
| Collect | transform stream to array |
| Count | return the count of elements in a stream |
| Max | return the maximum element in stream use the given ComparatorFunc |
| Min | return the minimal element in stream use the given ComparatorFunc |