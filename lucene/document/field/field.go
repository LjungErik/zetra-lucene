package field

import "github.com/LjungErik/zetra-lucene/lucene/document"

type Field struct {
	fieldType document.DocumentFieldType
	name      string
	value     string
	stored    bool
	indexable bool
}

var _ document.DocumentField = (*Field)(nil)

type Option func(*Field)

func WithValue(value string) Option {
	return func(f *Field) {
		f.value = value
	}
}

func WithStored(stored bool) Option {
	return func(f *Field) {
		f.stored = stored
	}
}

func WithType(fieldType document.DocumentFieldType) Option {
	return func(f *Field) {
		f.fieldType = fieldType
	}
}

func WithIndexable(indexable bool) Option {
	return func(f *Field) {
		f.indexable = indexable
	}
}

func New(name string, opts ...Option) *Field {
	f := &Field{name: name}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

func (f *Field) Name() string {
	return f.name
}

func (f *Field) Stored() bool {
	return f.stored
}

func (f *Field) Type() document.DocumentFieldType {
	return f.fieldType
}

func (f *Field) Indexable() bool {
	return f.indexable
}

func (f *Field) ValueAsString() string {
	return f.value
}

func (f *Field) ValueLength() int {
	return len(f.value)
}
