package fst

import "sort"

type State struct {
	arcs    map[byte]*Arc
	isFinal bool
}

type Arc struct {
	input  byte
	output uint64
	target *State
}

type FST struct {
	// Need to store the transitions
	root *State
}

func (f *FST) Get(key string) (uint64, bool) {
	state := f.root
	var output uint64

	for i := range key {
		b := key[i]
		t, ok := state.arcs[b]
		if !ok {
			return 0, false
		}
		output += t.output
		state = t.target
	}

	if !state.isFinal {
		return 0, false
	}

	return output, true
}

type Builder struct {
	register map[string]uint64
}

type entry struct {
	key    string
	offset uint64
}

func NewBuilder() *Builder {
	return &Builder{
		register: make(map[string]uint64),
	}
}

func (b *Builder) Insert(key string, offset uint64) {
	if key == "" {
		// Empty string is not allowed
		return
	}

	b.register[key] = offset
}

func (b *Builder) Build() *FST {
	fst := &FST{
		root: &State{
			arcs:    make(map[byte]*Arc),
			isFinal: false,
		},
	}

	registry := make([]entry, 0, len(b.register))

	for k, v := range b.register {
		registry = append(registry, entry{
			key:    k,
			offset: v,
		})
	}

	sort.SliceStable(registry, func(i, j int) bool {
		return registry[i].key < registry[j].key
	})

	for _, r := range registry {
		state := fst.root
		offsetLeft := r.offset
		for i := range r.key {
			k := r.key[i]
			arc, ok := state.arcs[k]
			if !ok {
				arc = &Arc{
					input:  k,
					output: offsetLeft,
					target: &State{
						arcs:    make(map[byte]*Arc),
						isFinal: false,
					},
				}
				offsetLeft = 0
				state.arcs[k] = arc
			} else {
				offsetLeft -= arc.output
			}

			state = arc.target
		}

		state.isFinal = true
	}

	return fst
}
