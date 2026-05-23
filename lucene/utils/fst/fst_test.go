package fst_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/utils/fst"
	"github.com/stretchr/testify/assert"
)

type testValue struct {
	key            string
	expectedOffset uint64
	expectedBool   bool
}

type testEntry struct {
	key    string
	offset uint64
}

func Test_FST_Builder(t *testing.T) {
	tests := []struct {
		name       string
		input      []testEntry
		testValues []testValue
	}{
		{
			name: "Basic test",
			input: []testEntry{
				{
					key:    "act",
					offset: 5,
				},
				{
					key:    "air",
					offset: 10,
				},
				{
					key:    "bannana",
					offset: 20,
				},
				{
					key:    "batman",
					offset: 35,
				},
				{
					key:    "dog",
					offset: 58,
				},
				{
					key:    "cats",
					offset: 42,
				},
				{
					key:    "the",
					offset: 102,
				},
				{
					key:    "theater",
					offset: 133,
				},
			},
			testValues: []testValue{
				{
					key:            "act",
					expectedOffset: 5,
					expectedBool:   true,
				},
				{
					key:            "air",
					expectedOffset: 10,
					expectedBool:   true,
				},
				{
					key:            "bannana",
					expectedOffset: 20,
					expectedBool:   true,
				},
				{
					key:            "batman",
					expectedOffset: 35,
					expectedBool:   true,
				},
				{
					key:            "dog",
					expectedOffset: 58,
					expectedBool:   true,
				},
				{
					key:            "magic",
					expectedOffset: 0,
					expectedBool:   false,
				},
				{
					key:            "cats",
					expectedOffset: 42,
					expectedBool:   true,
				},
				{
					key:            "cat",
					expectedOffset: 0,
					expectedBool:   false,
				},
				{
					key:            "the",
					expectedOffset: 102,
					expectedBool:   true,
				},
				{
					key:            "theater",
					expectedOffset: 133,
					expectedBool:   true,
				},
				{
					key:            "cats",
					expectedOffset: 42,
					expectedBool:   true,
				},
				{
					key:            "there",
					expectedOffset: 0,
					expectedBool:   false,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			builder := fst.NewBuilder()
			for _, e := range test.input {
				builder.Insert(e.key, e.offset)
			}

			f := builder.Build()

			for _, tv := range test.testValues {
				offset, ok := f.Get(tv.key)
				assert.Equal(t, tv.expectedOffset, offset)
				assert.Equal(t, tv.expectedBool, ok)
			}
		})
	}
}
