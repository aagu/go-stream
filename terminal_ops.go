package stream

type terminalOp struct {
	baseStage
	started bool
}

func (t *terminalOp) begin(_ int) {}

func (t *terminalOp) end() {
	t.startStage.closed = true
}

type collectOp struct {
	terminalOp
	data []interface{}
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

func (c *countOp) accept(t interface{}) {
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
