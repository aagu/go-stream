package stream

import "fmt"

type streamer int

const (
	opFilter streamer = iota
	opMapper
	opFlatMapper
	opSkipper
	opLimiter
	opSorter
	OpGrouper
	OpParalleled
	opCollector
	opDistincter
	opLooper
	opMaximizer
	opMinimizer
	opCounter
	opFirst
	opLast
)

// wrapSink is a helper function takes care of creating different kind of stages
func wrapSink(b *baseStage, s streamer, callback ...interface{}) stage {
	var nextStage stage
	switch s {
	case opFilter:
		downStream := new(filterOp)
		checkCallback("filter", callback)
		downStream.filterFunc = callback[0].(FilterFunc)
		nextStage = downStream
	case opMapper:
		downStream := new(mapperOp)
		checkCallback("mapper", callback)
		downStream.mapperFunc = callback[0].(MapFunc)
		nextStage = downStream
	case opFlatMapper:
		downStream := new(flatMapperOp)
		checkCallback("flatMap", callback)
		downStream.flatMapFunc = callback[0].(FlatMapFunc)
		nextStage = downStream
	case opDistincter:
		downStream := new(distinctOp)
		nextStage = downStream
	case opSkipper:
		downStream := new(skipperOp)
		checkCallback("skip", callback)
		downStream.skipSize = callback[0].(int)
		nextStage = downStream
	case opLimiter:
		downStream := new(limitOp)
		checkCallback("limit", callback)
		downStream.limitSize = callback[0].(int)
		nextStage = downStream
	case opSorter:
		downStream := new(sorterOp)
		checkCallback("sort", callback)
		downStream.comparator = callback[0].(ComparatorFunc)
		nextStage = downStream
	case OpGrouper:
		downStream := new(GroupOp)
		checkCallback("group", callback)
		downStream.groupFunc = callback[0].(GroupFunc)
		nextStage = downStream
	case OpParalleled:
		downStream := new(parallelStage)
		nextStage = downStream
		b.paralleled = true
	case opCollector:
		downStream := new(collectOp)
		nextStage = downStream
	case opLooper:
		downStream := new(forEachOp)
		checkCallback("forEach", callback)
		downStream.forEach = callback[0].(ForEachFunc)
		nextStage = downStream
	case opMaximizer:
		downStream := new(maxOp)
		checkCallback("max", callback)
		downStream.comparator = callback[0].(ComparatorFunc)
		nextStage = downStream
	case opMinimizer:
		downStream := new(minOp)
		checkCallback("min", callback)
		downStream.comparator = callback[0].(ComparatorFunc)
		nextStage = downStream
	case opCounter:
		downStream := new(countOp)
		nextStage = downStream
	case opFirst:
		downStream := new(firstOp)
		nextStage = downStream
	case opLast:
		downStream := new(lastOp)
		nextStage = downStream
	default:
		panic(fmt.Sprintf("unknown op %v", s))
	}
	b.setNextSink(nextStage)
	nextStage.setStartStage(b.getStartStage())
	return nextStage
}

func checkCallback(name string, callback ...interface{}) {
	if len(callback) == 0 {
		panic(fmt.Sprintf("not callback function found for %s", name))
	}
	if callback[0] == nil {
		panic("callback function could not be nil")
	}
}
