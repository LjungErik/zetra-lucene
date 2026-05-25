package fst

import "errors"

const (
	BIT_FINAL_ARC = 0x01
	BIT_LAST_ARC  = 0x02
	BIT_STOP_NODE = 0x08
)

type Cursor struct {
	state   int
	output  uint64
	isFinal bool
	done    bool
}

type Arc struct {
	flags  byte
	label  byte
	target int
	output uint64
}

type State struct {
	arcs []*Arc
}

type FST struct {
	states []*State
	root   int
}

func newFST() *FST {
	return &FST{
		states: make([]*State, 0),
		root:   0,
	}
}

func (f *FST) addKey(key string, offset uint64) {
	stateIdx := f.root
	remaining := offset

	for i := range key {
		isLastByte := i == (len(key) - 1)
		arc, existing := f.findOrCreateArc(stateIdx, key[i], remaining)

		if existing {
			remaining -= arc.output

			if !isLastByte {
				arc.flags &^= BIT_STOP_NODE
			}
		} else {
			remaining = 0
		}

		if isLastByte {
			arc.flags |= BIT_FINAL_ARC
			if !f.hasOutgoingArcs(arc.target) {
				arc.flags |= BIT_STOP_NODE
			}
		}

		stateIdx = arc.target
	}
}

func (f *FST) getOrCreateState(stateIdx int) *State {
	for stateIdx >= len(f.states) {
		f.states = append(f.states, &State{})
	}

	return f.states[stateIdx]
}

func (f *FST) findOrCreateArc(stateIdx int, label byte, ouput uint64) (*Arc, bool) {
	state := f.getOrCreateState(stateIdx)

	for _, arc := range state.arcs {
		if arc.label == label {
			return arc, true
		}

		if arc.label > label {
			break
		}
	}

	if n := len(state.arcs); n > 0 {
		state.arcs[n-1].flags &^= (BIT_LAST_ARC)
	}

	targetIdx := len(f.states)
	f.states = append(f.states, &State{})

	arc := &Arc{
		flags:  BIT_LAST_ARC,
		label:  label,
		target: targetIdx,
		output: ouput,
	}

	state.arcs = append(state.arcs, arc)

	return arc, false
}

func (f *FST) hasOutgoingArcs(stateIdx int) bool {
	return len(f.states[stateIdx].arcs) > 0
}

func (f *FST) start() *Cursor {
	return &Cursor{state: f.root}
}

func (f *FST) step(c *Cursor, label byte) *Cursor {
	if c.done {
		return &Cursor{done: true}
	}

	state := f.states[c.state]
	for _, arc := range state.arcs {
		if arc.label == label {
			next := &Cursor{
				output:  c.output + arc.output,
				isFinal: arc.flags&BIT_FINAL_ARC != 0,
			}
			if arc.flags&BIT_STOP_NODE != 0 {
				next.done = true
			} else {
				next.state = arc.target
			}

			return next
		}
		if arc.label > label {
			break
		}

		if arc.flags&BIT_LAST_ARC != 0 {
			break
		}
	}

	return &Cursor{done: true}
}

func (f *FST) Get(key string) (uint64, bool) {
	c := f.start()
	for i := 0; i < len(key); i++ {
		c = f.step(c, key[i])
		if c.done && i < len(key)-1 {
			return 0, false
		}
	}

	if !c.isFinal {
		return 0, false
	}

	return c.output, true
}

func (f *FST) LookupBlock(key string) (uint64, bool) {
	c := f.start()
	var lastFinal uint64
	haveFinal := false
	for i := 0; i < len(key); i++ {
		c = f.step(c, key[i])
		if c.isFinal {
			lastFinal = c.output
			haveFinal = true
		}
		if c.done {
			break
		}
	}

	return lastFinal, haveFinal
}

var (
	ErrInvalidKey    = errors.New("invalid key, earlier keys are greater than given key")
	ErrInvalidOffset = errors.New("invalid offset, current highest offset is larger than given offset")
)

type entry struct {
	key    string
	offset uint64
}

type Builder struct {
	lastKey    string // determine if we get a incorrect value back (inserts must be preformed in order)
	lastOffset uint64
	registry   []entry
}

func NewBuilder() *Builder {
	return &Builder{
		lastKey:    "",
		lastOffset: 0,
		registry:   make([]entry, 0),
	}
}

func (b *Builder) Insert(key string, offset uint64) error {
	if key <= b.lastKey && b.lastKey != "" {
		return ErrInvalidKey
	}

	if offset <= b.lastOffset {
		return ErrInvalidOffset
	}

	b.registry = append(b.registry, entry{
		key:    key,
		offset: offset,
	})

	b.lastKey = key
	b.lastOffset = offset

	return nil
}

func (b *Builder) Build() *FST {
	fst := newFST()

	for _, e := range b.registry {
		fst.addKey(e.key, e.offset)
	}

	return fst
}
