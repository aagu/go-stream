package stream

import "sort"

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

func (m *mapperOp) end() {
	m.downStream.end()
}

type skipperOp struct {
	baseStage
	skipSize  int
	skipCount int
}

func (s *skipperOp) begin(size int) {
	s.downStream.begin(size - s.skipSize)
}

func (s *skipperOp) accept(t interface{}) {
	if s.skipCount >= s.skipSize && !s.downStream.cancellationRequested() {
		s.downStream.accept(t)
	} else {
		s.skipCount++
	}
}

type sorterOp struct {
	baseStage
	comparator ComparatorFunc
	data       []interface{}
}

func (s *sorterOp) begin(size int) {
	if size > 0 {
		s.data = make([]interface{}, 0, size)
	} else {
		s.data = make([]interface{}, 0)
	}
}

func (s *sorterOp) accept(t interface{}) {
	s.data = append(s.data, t)
}

func (s *sorterOp) end() {
	sort.Slice(s.data, func(i, j int) bool {
		return s.comparator(s.data[i], s.data[j]) <= 0
	})
	s.downStream.begin(len(s.data))
	for idx := range s.data {
		s.downStream.accept(s.data[idx])
	}
	s.downStream.end()
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
		if !f.downStream.cancellationRequested() {
			f.downStream.accept(flatted[idx])
		}
	}
}

type limitOp struct {
	baseStage
	limitSize  int
	limitCount int
}

func (l *limitOp) begin(_ int) {
	l.downStream.begin(l.limitSize)
}

func (l *limitOp) accept(t interface{}) {
	l.limitCount++
	l.downStream.accept(t)
}

func (l *limitOp) cancellationRequested() bool {
	return l.limitCount >= l.limitSize
}

type distinctOp struct {
	terminalOp
	set map[interface{}]byte // temp storage
}

func (d *distinctOp) begin(_ int) {
	d.set = make(map[interface{}]byte)
}

func (d *distinctOp) accept(t interface{}) {
	d.set[t] = 0x01
}

func (d *distinctOp) end() {
	d.downStream.begin(len(d.set))
	for key := range d.set {
		if !d.downStream.cancellationRequested() {
			d.downStream.accept(key)
		}
	}
	d.downStream.end()
}
