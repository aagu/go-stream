package stream

type filterOp struct {
	baseStage
	filterFunc FilterFunc
}

func (f *filterOp) begin(_ int) {
	f.downStream.begin(0)
}

func (f *filterOp) accept(t interface{}) {
	if f.downStream.cancellationRequested() {
		return
	}
	if f.filterFunc(t) {
		f.downStream.accept(t)
	}
}

type mapperOp struct {
	baseStage
	mapperFunc MapFunc
}

func (m *mapperOp) accept(t interface{}) {
	if !m.downStream.cancellationRequested() {
		m.downStream.accept(m.mapperFunc(t))
	}
}

type flatMapperOp struct {
	baseStage
	flatMapFunc FlatMapFunc
}

func (f *flatMapperOp) begin(_ int) {
	f.downStream.begin(0)
}

func (f *flatMapperOp) accept(t interface{}) {
	flatted := f.flatMapFunc(t)
	for idx := range flatted {
		if f.downStream.cancellationRequested() {
			break
		}
		f.downStream.accept(flatted[idx])
	}
}
