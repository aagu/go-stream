package stream

// the basic struct of a terminal operation
// each Stream ends with a terminal operation
// and all the operations in Stream will not be performed before
// calling terminal operation
type terminalOp struct {
	baseStage
	started bool
}

func (t *terminalOp) end() {
	t.startStage.closed = true
}

type collectOp struct {
	terminalOp
	data []interface{}
}

func (c *collectOp) accept(t interface{}) {
	c.data = append(c.data, t)
}

type forEachOp struct {
	terminalOp
	forEach ForEachFunc
}

func (f *forEachOp) accept(t interface{}) {
	f.forEach(t)
}

type countOp struct {
	terminalOp
	count int
}

func (c *countOp) accept(_ interface{}) {
	c.count++
}

type maxOp struct {
	terminalOp
	comparator ComparatorFunc
	max        interface{}
}

func (m *maxOp) accept(t interface{}) {
	if m.max == nil {
		m.max = t
	} else if m.comparator(m.max, t) < 0 {
		m.max = t
	}
}

type minOp struct {
	terminalOp
	comparator ComparatorFunc
	min        interface{}
}

func (m *minOp) accept(t interface{}) {
	if m.min == nil {
		m.min = t
	} else if m.comparator(m.min, t) > 0 {
		m.min = t
	}
}

type firstOp struct {
	terminalOp
	val    interface{}
	cancel bool
}

func (f *firstOp) accept(v interface{}) {
	f.val = v
	f.cancel = true
}

func (f *firstOp) cancellationRequested() bool {
	return f.cancel
}

type lastOp struct {
	terminalOp
	val interface{}
}

func (l *lastOp) accept(v interface{}) {
	l.val = v
}
