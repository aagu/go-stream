package stream

import (
	"sort"
	"sync"
)

type statefulOp struct {
	baseStage
	l sync.Mutex // for parallel stream synchronization
}

type skipperOp struct {
	statefulOp
	skipSize  int
	skipCount int
}

func (s *skipperOp) begin(size int) {
	s.downStream.begin(size - s.skipSize)
}

func (s *skipperOp) accept(t interface{}) {
	s.l.Lock()
	if s.skipCount >= s.skipSize && !s.downStream.cancellationRequested() {
		s.l.Unlock()
		s.downStream.accept(t)
	} else {
		s.skipCount++
		s.l.Unlock()
	}
}

type sorterOp struct {
	statefulOp
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
	s.l.Lock()
	s.data = append(s.data, t)
	s.l.Unlock()
}

func (s *sorterOp) end() {
	sort.Slice(s.data, func(i, j int) bool {
		return s.comparator(s.data[i], s.data[j]) <= 0
	})
	s.downStream.begin(len(s.data))
	for idx := range s.data {
		if s.downStream.cancellationRequested() { // check first, since accept may be called many times by upstream
			break
		}
		s.downStream.accept(s.data[idx])
	}
	s.downStream.end()
}

type limitOp struct {
	statefulOp
	limitSize  int
	limitCount int
}

func (l *limitOp) begin(_ int) {
	l.downStream.begin(l.limitSize)
}

func (l *limitOp) accept(t interface{}) {
	l.l.Lock()
	l.limitCount++
	l.l.Unlock()
	if !l.downStream.cancellationRequested() {
		l.downStream.accept(t)
	}
}

func (l *limitOp) cancellationRequested() bool {
	l.l.Lock()
	defer l.l.Unlock()
	return l.limitCount >= l.limitSize
}

type distinctOp struct {
	statefulOp
	set map[interface{}]struct{} // temp storage
}

func (d *distinctOp) begin(_ int) {
	d.set = make(map[interface{}]struct{})
}

func (d *distinctOp) accept(t interface{}) {
	d.set[t] = struct{}{}
}

func (d *distinctOp) end() {
	d.downStream.begin(len(d.set))
	for key := range d.set {
		if d.downStream.cancellationRequested() {
			break
		}
		d.downStream.accept(key)
	}
	d.downStream.end()
}

type funcDistinctOp struct {
	statefulOp
	set map[interface{}]interface{}
	fn  DistinctFunc
}

func (f *funcDistinctOp) begin(_ int) {
	f.set = make(map[interface{}]interface{})
}

func (f *funcDistinctOp) accept(t interface{}) {
	f.set[f.fn(t)] = t
}

func (f *funcDistinctOp) end() {
	f.downStream.begin(len(f.set))
	for _, v := range f.set {
		if f.downStream.cancellationRequested() {
			break
		}
		f.downStream.accept(v)
	}
	f.downStream.end()
}

type GroupOp struct {
	statefulOp
	groupFunc GroupFunc
	groups    map[interface{}][]interface{}
}

func (g *GroupOp) begin(_ int) {
	g.groups = make(map[interface{}][]interface{})
}

func (g *GroupOp) accept(t interface{}) {
	key := g.groupFunc(t)
	if g.groups[key] == nil {
		g.groups[key] = make([]interface{}, 0)
	}
	g.l.Lock()
	g.groups[key] = append(g.groups[key], t)
	g.l.Unlock()
}

func (g *GroupOp) end() {
	g.downStream.begin(len(g.groups))
	for _, value := range g.groups {
		if g.downStream.cancellationRequested() {
			break
		}
		g.downStream.accept(value)
	}
	g.downStream.end()
}
