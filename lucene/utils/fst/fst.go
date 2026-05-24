package fst

import "errors"

const (
	BIT_FINAL_ARC = 0x01
	BIT_LAST_ARC  = 0x02
)

type Arc struct {
	flags  byte
	label  byte
	target int
	output uint64
}

type State struct {
	arcs []Arc
}

type FST struct {
	states []State
	root   int
}

func (f *FST) Get(key string) (uint64, bool) {
	state := f.states[f.root]
	var output uint64

	for i := range key {
		label := key[i] // byte key
		for j, arc := range state.arcs {
			if arc.label == label {
				output += arc.output
				// We have found the final matching arc
				if arc.flags&BIT_FINAL_ARC != 0 {
					return output, true
				}

				state = f.states[arc.target]
				break
			} else if arc.flags&BIT_LAST_ARC != 0 || j == (len(state.arcs)-1) {
				// This means we have no matches
				return 0, false
			}
		}
	}

	return 0, false

}

var (
	ErrInvalidInsertOrder = errors.New("invalid insert order, earlier keys are greater than given key")
)

type entry struct {
	key    string
	offset uint64
}

type Builder struct {
	lastKey  string // determine if we get a incorrect value back (inserts must be preformed in order)
	registry []entry
}

func commonPrefixLen(s1, s2 string) int {
	n := len(s1)
	if len(s2) < n {
		n = len(s2)
	}

	for i := range n {
		if s1[i] != s2[i] {
			return i
		}
	}

	return n
}

func (b *Builder) Insert(key string, offset uint64) error {
	if key <= b.lastKey && b.lastKey != "" {
		return ErrInvalidInsertOrder
	}

	prefixLen := commonPrefixLen(key, b.lastKey)
}
