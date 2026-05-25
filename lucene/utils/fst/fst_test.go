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

func Test_FST_Get(t *testing.T) {
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
					key:    "ba",
					offset: 20,
				},
				{
					key:    "bat",
					offset: 35,
				},
				{
					key:    "cat",
					offset: 42,
				},
				{
					key:    "dog",
					offset: 58,
				},
				{
					key:    "the",
					offset: 102,
				},
				{
					key:    "thea",
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
					key:            "ba",
					expectedOffset: 20,
					expectedBool:   true,
				},
				{
					key:            "bannana",
					expectedOffset: 0,
					expectedBool:   false,
				},
				{
					key:            "bat",
					expectedOffset: 35,
					expectedBool:   true,
				},
				{
					key:            "batman",
					expectedOffset: 0,
					expectedBool:   false,
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
					expectedOffset: 0,
					expectedBool:   false,
				},
				{
					key:            "cat",
					expectedOffset: 42,
					expectedBool:   true,
				},
				{
					key:            "the",
					expectedOffset: 102,
					expectedBool:   true,
				},
				{
					key:            "theater",
					expectedOffset: 0,
					expectedBool:   false,
				},
				{
					key:            "thea",
					expectedOffset: 133,
					expectedBool:   true,
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

func Test_FST_LookupBlock(t *testing.T) {
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
					key:    "ba",
					offset: 20,
				},
				{
					key:    "bat",
					offset: 35,
				},
				{
					key:    "c",
					offset: 42,
				},
				{
					key:    "d",
					offset: 58,
				},
				{
					key:    "the",
					offset: 102,
				},
				{
					key:    "thea",
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
					key:            "actor",
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
					expectedOffset: 42,
					expectedBool:   true,
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
					key:            "there",
					expectedOffset: 102,
					expectedBool:   true,
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
				offset, ok := f.LookupBlock(tv.key)
				assert.Equal(t, tv.expectedOffset, offset)
				assert.Equal(t, tv.expectedBool, ok)
			}
		})
	}
}
